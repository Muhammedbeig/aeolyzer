package mcptransportplane

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	LatestProtocolVersion = "2025-11-25"

	headerProtocolVersion = "MCP-Protocol-Version"
	headerSessionID       = "MCP-Session-Id"
	headerLastEventID     = "Last-Event-ID"
)

// PeerInfo identifies the MCP client or server during initialization.
type PeerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// InitializeResult is the MCP initialize response payload.
type InitializeResult struct {
	ProtocolVersion string         `json:"protocolVersion"`
	Capabilities    map[string]any `json:"capabilities"`
	ServerInfo      PeerInfo       `json:"serverInfo"`
}

// HTTPClient performs bounded MCP Streamable HTTP calls over required TLS.
type HTTPClient struct {
	endpoint *url.URL
	client   *http.Client
	maxBytes int64

	mu              sync.Mutex
	initialized     bool
	protocolVersion string
	sessionID       string
}

// NewHTTPClient creates an exact-endpoint, no-redirect, mTLS-capable client.
func NewHTTPClient(
	endpoint string,
	tlsConfig *tls.Config,
	timeout time.Duration,
	maxBytes int64,
) (*HTTPClient, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil ||
		parsed.Scheme != "https" ||
		parsed.Hostname() == "" ||
		parsed.User != nil ||
		parsed.Fragment != "" ||
		tlsConfig == nil ||
		tlsConfig.RootCAs == nil ||
		tlsConfig.MinVersion < tls.VersionTLS12 ||
		timeout < time.Second ||
		timeout > 2*time.Minute ||
		maxBytes < 1 ||
		maxBytes > 8<<20 {
		return nil, errors.New("mcp http client policy is invalid")
	}
	transport := &http.Transport{
		Proxy:             nil,
		DisableKeepAlives: true,
		TLSClientConfig:   tlsConfig.Clone(),
	}
	return &HTTPClient{
		endpoint: parsed,
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				return errors.New("mcp redirects are blocked")
			},
		},
		maxBytes: maxBytes,
	}, nil
}

// Initialize performs the MCP 2025-11-25 initialize handshake.
func (c *HTTPClient) Initialize(ctx context.Context, clientInfo PeerInfo) (InitializeResult, error) {
	if strings.TrimSpace(clientInfo.Name) == "" || strings.TrimSpace(clientInfo.Version) == "" {
		return InitializeResult{}, errors.New("mcp client info is invalid")
	}
	params, err := json.Marshal(struct {
		ProtocolVersion string         `json:"protocolVersion"`
		Capabilities    map[string]any `json:"capabilities"`
		ClientInfo      PeerInfo       `json:"clientInfo"`
	}{
		ProtocolVersion: LatestProtocolVersion,
		Capabilities:    map[string]any{},
		ClientInfo:      clientInfo,
	})
	if err != nil {
		return InitializeResult{}, fmt.Errorf("encode mcp initialize params: %w", err)
	}
	requestID := json.RawMessage(`"initialize-1"`)
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      requestID,
		Method:  "initialize",
		Params:  params,
	}
	data, err := EncodeRequest(request)
	if err != nil {
		return InitializeResult{}, err
	}
	response, sessionID, err := c.doJSONRPC(ctx, http.MethodPost, data, "", false, requestID)
	if err != nil {
		return InitializeResult{}, err
	}
	var result InitializeResult
	if err := json.Unmarshal(response.Result, &result); err != nil {
		return InitializeResult{}, fmt.Errorf("decode mcp initialize result: %w", err)
	}
	if result.ProtocolVersion != LatestProtocolVersion {
		return InitializeResult{}, errors.New("mcp protocol version is unsupported")
	}
	c.mu.Lock()
	c.initialized = true
	c.protocolVersion = result.ProtocolVersion
	c.sessionID = sessionID
	c.mu.Unlock()

	if err := c.Notify(ctx, JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}); err != nil {
		c.mu.Lock()
		c.initialized = false
		c.protocolVersion = ""
		c.sessionID = ""
		c.mu.Unlock()
		return InitializeResult{}, err
	}
	return result, nil
}

// Call sends one initialized JSON-RPC request and validates one JSON-RPC response.
func (c *HTTPClient) Call(ctx context.Context, request JSONRPCRequest) (JSONRPCResponse, error) {
	if !c.isInitialized() {
		return JSONRPCResponse{}, errors.New("mcp http client is not initialized")
	}
	data, err := EncodeRequest(request)
	if err != nil {
		return JSONRPCResponse{}, err
	}
	response, _, err := c.doJSONRPC(ctx, http.MethodPost, data, "", true, request.ID)
	return response, err
}

// Notify sends one initialized JSON-RPC notification.
func (c *HTTPClient) Notify(ctx context.Context, notification JSONRPCNotification) error {
	if !c.isInitialized() {
		return errors.New("mcp http client is not initialized")
	}
	data, err := EncodeNotification(notification)
	if err != nil {
		return err
	}
	return c.doNotification(ctx, data)
}

// Listen opens a bounded server event stream for initialized MCP sessions.
func (c *HTTPClient) Listen(ctx context.Context, lastEventID string) ([]JSONRPCResponse, error) {
	if !c.isInitialized() {
		return nil, errors.New("mcp http client is not initialized")
	}
	response, _, err := c.doJSONRPC(ctx, http.MethodGet, nil, lastEventID, true, nil)
	if err != nil {
		return nil, err
	}
	return []JSONRPCResponse{response}, nil
}

// Close closes a stateful MCP session when the server issued one.
func (c *HTTPClient) Close(ctx context.Context) error {
	protocolVersion, sessionID, initialized := c.state()
	if !initialized || sessionID == "" {
		return nil
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, c.endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("build mcp close request: %w", err)
	}
	request.Header.Set(headerProtocolVersion, protocolVersion)
	request.Header.Set(headerSessionID, sessionID)
	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("send mcp close request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("mcp close status %d", response.StatusCode)
	}
	c.mu.Lock()
	c.initialized = false
	c.protocolVersion = ""
	c.sessionID = ""
	c.mu.Unlock()
	return nil
}

// SessionID returns the server-issued MCP session identifier.
func (c *HTTPClient) SessionID() string {
	_, sessionID, _ := c.state()
	return sessionID
}

func (c *HTTPClient) doJSONRPC(
	ctx context.Context,
	method string,
	body []byte,
	lastEventID string,
	requireInitialized bool,
	expectedID json.RawMessage,
) (JSONRPCResponse, string, error) {
	if c == nil || c.endpoint == nil || c.client == nil {
		return JSONRPCResponse{}, "", errors.New("mcp http client is not configured")
	}
	requestBody := io.Reader(nil)
	if body != nil {
		requestBody = bytes.NewReader(body)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.endpoint.String(), requestBody)
	if err != nil {
		return JSONRPCResponse{}, "", fmt.Errorf("build mcp http request: %w", err)
	}
	request.Header.Set("Accept", "application/json, text/event-stream")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if requireInitialized {
		c.attachStateHeaders(request)
	}
	if lastEventID != "" {
		request.Header.Set(headerLastEventID, lastEventID)
	}
	response, err := c.client.Do(request)
	if err != nil {
		return JSONRPCResponse{}, "", fmt.Errorf("send mcp http request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return JSONRPCResponse{}, "", fmt.Errorf("mcp http status %d", response.StatusCode)
	}
	decoded, err := c.decodeHTTPResponse(response, expectedID)
	if err != nil {
		return JSONRPCResponse{}, "", err
	}
	return decoded, response.Header.Get(headerSessionID), nil
}

func (c *HTTPClient) doNotification(ctx context.Context, body []byte) error {
	if c == nil || c.endpoint == nil || c.client == nil {
		return errors.New("mcp http client is not configured")
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint.String(), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build mcp notification request: %w", err)
	}
	request.Header.Set("Accept", "application/json, text/event-stream")
	request.Header.Set("Content-Type", "application/json")
	c.attachStateHeaders(request)
	response, err := c.client.Do(request)
	if err != nil {
		return fmt.Errorf("send mcp notification request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusAccepted &&
		response.StatusCode != http.StatusOK &&
		response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("mcp notification status %d", response.StatusCode)
	}
	if response.StatusCode == http.StatusOK && response.Body != http.NoBody {
		body, err := io.ReadAll(io.LimitReader(response.Body, c.maxBytes+1))
		if err != nil {
			return fmt.Errorf("read mcp notification response: %w", err)
		}
		if len(bytes.TrimSpace(body)) != 0 {
			return errors.New("mcp notification returned an unexpected response body")
		}
	}
	return nil
}

func (c *HTTPClient) decodeHTTPResponse(response *http.Response, expectedID json.RawMessage) (JSONRPCResponse, error) {
	contentType, _, err := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if err != nil {
		return JSONRPCResponse{}, errors.New("mcp http response content type is invalid")
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, c.maxBytes+1))
	if err != nil {
		return JSONRPCResponse{}, fmt.Errorf("read mcp http response: %w", err)
	}
	if int64(len(body)) > c.maxBytes {
		return JSONRPCResponse{}, errors.New("mcp http response exceeds size limit")
	}
	var decoded JSONRPCResponse
	switch contentType {
	case "application/json":
		decoded, err = DecodeResponse(body)
	case "text/event-stream":
		var responses []JSONRPCResponse
		responses, err = parseSSEResponses(body)
		if err == nil {
			decoded, err = selectResponse(responses, expectedID)
		}
	default:
		err = errors.New("mcp http response content type is invalid")
	}
	if err != nil {
		return JSONRPCResponse{}, err
	}
	if len(expectedID) != 0 && !bytes.Equal(decoded.ID, expectedID) {
		return JSONRPCResponse{}, errors.New("mcp json-rpc response id mismatch")
	}
	return decoded, nil
}

func parseSSEResponses(data []byte) ([]JSONRPCResponse, error) {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 0, 64<<10), maxJSONRPCBytes)
	var responses []JSONRPCResponse
	var dataLines []string
	flush := func() error {
		if len(dataLines) == 0 {
			return nil
		}
		response, err := DecodeResponse([]byte(strings.Join(dataLines, "\n")))
		if err != nil {
			return err
		}
		responses = append(responses, response)
		dataLines = nil
		return nil
	}
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" {
			if err := flush(); err != nil {
				return nil, err
			}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		field, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = strings.TrimPrefix(value, " ")
		if field == "data" {
			dataLines = append(dataLines, value)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read mcp sse response: %w", err)
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return nil, errors.New("mcp sse response is empty")
	}
	return responses, nil
}

func selectResponse(responses []JSONRPCResponse, expectedID json.RawMessage) (JSONRPCResponse, error) {
	if len(expectedID) == 0 {
		return responses[0], nil
	}
	for _, response := range responses {
		if bytes.Equal(response.ID, expectedID) {
			return response, nil
		}
	}
	return JSONRPCResponse{}, errors.New("mcp sse response id mismatch")
}

func (c *HTTPClient) attachStateHeaders(request *http.Request) {
	protocolVersion, sessionID, _ := c.state()
	request.Header.Set(headerProtocolVersion, protocolVersion)
	if sessionID != "" {
		request.Header.Set(headerSessionID, sessionID)
	}
}

func (c *HTTPClient) state() (string, string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.protocolVersion, c.sessionID, c.initialized
}

func (c *HTTPClient) isInitialized() bool {
	_, _, initialized := c.state()
	return initialized
}
