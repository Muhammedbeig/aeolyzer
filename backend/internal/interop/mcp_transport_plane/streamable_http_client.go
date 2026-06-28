package mcptransportplane

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPClient performs bounded MCP Streamable HTTP calls over required TLS.
type HTTPClient struct {
	endpoint *url.URL
	client   *http.Client
	maxBytes int64
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

// Call sends one JSON-RPC request and validates one JSON-RPC response.
func (c *HTTPClient) Call(
	ctx context.Context,
	request JSONRPCRequest,
) (JSONRPCResponse, error) {
	if c == nil || c.endpoint == nil || c.client == nil {
		return JSONRPCResponse{}, errors.New("mcp http client is not configured")
	}
	data, err := EncodeRequest(request)
	if err != nil {
		return JSONRPCResponse{}, err
	}
	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint.String(),
		bytes.NewReader(data),
	)
	if err != nil {
		return JSONRPCResponse{}, fmt.Errorf("build mcp http request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept", "application/json, text/event-stream")
	response, err := c.client.Do(httpRequest)
	if err != nil {
		return JSONRPCResponse{}, fmt.Errorf("send mcp http request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return JSONRPCResponse{}, fmt.Errorf("mcp http status %d", response.StatusCode)
	}
	if contentType := response.Header.Get("Content-Type"); contentType != "application/json" {
		return JSONRPCResponse{}, errors.New("mcp http response content type is invalid")
	}
	body, err := io.ReadAll(io.LimitReader(response.Body, c.maxBytes+1))
	if err != nil {
		return JSONRPCResponse{}, fmt.Errorf("read mcp http response: %w", err)
	}
	if int64(len(body)) > c.maxBytes {
		return JSONRPCResponse{}, errors.New("mcp http response exceeds size limit")
	}
	return DecodeResponse(body)
}
