package mcptransportplane

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJSONRPCCodecRejectsUnknownFieldsAndAmbiguousResponse(t *testing.T) {
	if _, err := DecodeRequest([]byte(`{
		"jsonrpc":"2.0",
		"id":"1",
		"method":"tools/list",
		"unexpected":true
	}`)); err == nil {
		t.Fatal("DecodeRequest() accepted unknown field")
	}
	if _, err := DecodeResponse([]byte(`{
		"jsonrpc":"2.0",
		"id":"1",
		"result":{},
		"error":{"code":-1,"message":"failure"}
	}`)); err == nil {
		t.Fatal("DecodeResponse() accepted result and error together")
	}
}

func TestJSONRPCCodecRoundTrip(t *testing.T) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"request-1"`),
		Method:  "tools/list",
		Params:  json.RawMessage(`{}`),
	}
	data, err := EncodeRequest(request)
	if err != nil {
		t.Fatalf("EncodeRequest() failed: %v", err)
	}
	decoded, err := DecodeRequest(data)
	if err != nil {
		t.Fatalf("DecodeRequest() failed: %v", err)
	}
	if decoded.Method != request.Method {
		t.Fatalf("DecodeRequest().Method = %q, want %q", decoded.Method, request.Method)
	}
	notificationData, err := EncodeNotification(JSONRPCNotification{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	})
	if err != nil {
		t.Fatalf("EncodeNotification() failed: %v", err)
	}
	if _, err := DecodeRequest(notificationData); err == nil {
		t.Fatal("DecodeRequest() accepted notification without id")
	}
}

func TestValidateHandshakeRejectsExtraToolAndSchemaDrift(t *testing.T) {
	pinned := PinnedManifest{
		ServerID:                "gsc",
		AllowedProtocolVersions: []string{LatestProtocolVersion},
		RequiredCapabilities:    []string{"tools"},
		Tools: []ToolManifest{{
			Name:       "query",
			SchemaHash: "sha256:abc",
		}},
	}
	live := ServerManifest{
		ServerID:        "gsc",
		ProtocolVersion: LatestProtocolVersion,
		Capabilities:    []string{"tools"},
		Tools:           append([]ToolManifest(nil), pinned.Tools...),
	}
	if err := ValidateHandshake(live, pinned); err != nil {
		t.Fatalf("ValidateHandshake() failed: %v", err)
	}
	live.Tools = append(live.Tools, ToolManifest{Name: "exfiltrate", SchemaHash: "sha256:def"})
	if err := ValidateHandshake(live, pinned); err == nil {
		t.Fatal("ValidateHandshake() accepted extra tool")
	}
	live.Tools = []ToolManifest{{Name: "query", SchemaHash: "sha256:changed"}}
	if err := ValidateHandshake(live, pinned); err == nil {
		t.Fatal("ValidateHandshake() accepted schema drift")
	}
}

func TestSchemaHashCanonicalizesJSON(t *testing.T) {
	left, err := SchemaHash([]byte(`{"type":"object","properties":{"x":{"type":"string"}}}`))
	if err != nil {
		t.Fatalf("SchemaHash(left) failed: %v", err)
	}
	right, err := SchemaHash([]byte(`{
		"type": "object",
		"properties": {"x": {"type": "string"}}
	}`))
	if err != nil {
		t.Fatalf("SchemaHash(right) failed: %v", err)
	}
	if left != right {
		t.Fatalf("SchemaHash() values differ: %s != %s", left, right)
	}
}

func TestCircuitBreakerOpensAndAllowsSingleProbe(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	breaker, err := NewCircuitBreaker(2, time.Minute, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewCircuitBreaker() failed: %v", err)
	}
	breaker.Failure()
	if err := breaker.Allow(); err != nil {
		t.Fatalf("CircuitBreaker.Allow() failed before threshold: %v", err)
	}
	breaker.Failure()
	if err := breaker.Allow(); !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("CircuitBreaker.Allow() error = %v, want open", err)
	}
	now = now.Add(time.Minute)
	if err := breaker.Allow(); err != nil {
		t.Fatalf("CircuitBreaker.Allow() probe failed: %v", err)
	}
	if err := breaker.Allow(); !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("second half-open probe error = %v, want open", err)
	}
	breaker.Success()
	if breaker.State() != CircuitClosed {
		t.Fatalf("CircuitBreaker.State() = %s, want closed", breaker.State())
	}
}

func TestHTTPClientInitializeCallAndClose(t *testing.T) {
	var sawInitializedNotification bool
	var sawCall bool
	var sawClose bool
	client, server := newMCPHTTPClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			method := requestMethod(t, request)
			switch method {
			case "initialize":
				if got := request.Header.Get(headerProtocolVersion); got != "" {
					t.Fatalf("initialize %s header = %q, want empty", headerProtocolVersion, got)
				}
				response.Header().Set("Content-Type", "application/json; charset=utf-8")
				response.Header().Set(headerSessionID, "session-1")
				_, _ = response.Write([]byte(`{
					"jsonrpc":"2.0",
					"id":"initialize-1",
					"result":{
						"protocolVersion":"2025-11-25",
						"capabilities":{"tools":{"listChanged":true}},
						"serverInfo":{"name":"test-mcp","version":"1.0.0"}
					}
				}`))
			case "notifications/initialized":
				requireMCPHeaders(t, request)
				sawInitializedNotification = true
				response.WriteHeader(http.StatusAccepted)
			case "tools/list":
				requireMCPHeaders(t, request)
				sawCall = true
				response.Header().Set("Content-Type", "application/json")
				_, _ = response.Write([]byte(`{"jsonrpc":"2.0","id":"tools-1","result":{"tools":[]}}`))
			default:
				t.Fatalf("unexpected method %q", method)
			}
		case http.MethodDelete:
			requireMCPHeaders(t, request)
			sawClose = true
			response.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected HTTP method %s", request.Method)
		}
	}))
	defer server.Close()

	result, err := client.Initialize(context.Background(), PeerInfo{Name: "aeolyzer", Version: "test"})
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	if result.ProtocolVersion != LatestProtocolVersion {
		t.Fatalf("Initialize().ProtocolVersion = %q, want %q", result.ProtocolVersion, LatestProtocolVersion)
	}
	if got, want := client.SessionID(), "session-1"; got != want {
		t.Fatalf("SessionID() = %q, want %q", got, want)
	}
	_, err = client.Call(context.Background(), JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"tools-1"`),
		Method:  "tools/list",
		Params:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("Call() failed: %v", err)
	}
	if err := client.Close(context.Background()); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}
	if !sawInitializedNotification || !sawCall || !sawClose {
		t.Fatalf("server saw initialized=%t call=%t close=%t, want all true", sawInitializedNotification, sawCall, sawClose)
	}
}

func TestHTTPClientRejectsCallBeforeInitialize(t *testing.T) {
	client, server := newMCPHTTPClient(t, http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("server should not be called before initialization")
	}))
	defer server.Close()
	_, err := client.Call(context.Background(), JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"tools-1"`),
		Method:  "tools/list",
	})
	if err == nil {
		t.Fatal("Call() succeeded before Initialize()")
	}
}

func TestHTTPClientParsesSSEResponse(t *testing.T) {
	client, server := newMCPHTTPClient(t, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodPost:
			method := requestMethod(t, request)
			switch method {
			case "initialize":
				response.Header().Set("Content-Type", "application/json")
				response.Header().Set(headerSessionID, "session-1")
				_, _ = response.Write([]byte(`{
					"jsonrpc":"2.0",
					"id":"initialize-1",
					"result":{
						"protocolVersion":"2025-11-25",
						"capabilities":{},
						"serverInfo":{"name":"test-mcp","version":"1.0.0"}
					}
				}`))
			case "notifications/initialized":
				requireMCPHeaders(t, request)
				response.WriteHeader(http.StatusAccepted)
			case "tools/list":
				requireMCPHeaders(t, request)
				response.Header().Set("Content-Type", "text/event-stream")
				_, _ = response.Write([]byte("id: event-1\ndata: {\"jsonrpc\":\"2.0\",\"id\":\"tools-1\",\"result\":{\"tools\":[]}}\n\n"))
			default:
				t.Fatalf("unexpected method %q", method)
			}
		case http.MethodGet:
			requireMCPHeaders(t, request)
			if got, want := request.Header.Get(headerLastEventID), "event-1"; got != want {
				t.Fatalf("%s = %q, want %q", headerLastEventID, got, want)
			}
			response.Header().Set("Content-Type", "text/event-stream")
			_, _ = response.Write([]byte("id: event-2\ndata: {\"jsonrpc\":\"2.0\",\"id\":\"stream-1\",\"result\":{\"status\":\"ok\"}}\n\n"))
		default:
			t.Fatalf("unexpected HTTP method %s", request.Method)
		}
	}))
	defer server.Close()

	if _, err := client.Initialize(context.Background(), PeerInfo{Name: "aeolyzer", Version: "test"}); err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}
	response, err := client.Call(context.Background(), JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      json.RawMessage(`"tools-1"`),
		Method:  "tools/list",
		Params:  json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("Call() failed: %v", err)
	}
	if len(response.Result) == 0 {
		t.Fatal("Call() returned empty result")
	}
	responses, err := client.Listen(context.Background(), "event-1")
	if err != nil {
		t.Fatalf("Listen() failed: %v", err)
	}
	if len(responses) != 1 || string(responses[0].ID) != `"stream-1"` {
		t.Fatalf("Listen() responses = %#v, want stream-1", responses)
	}
}

func newMCPHTTPClient(t *testing.T, handler http.Handler) (*HTTPClient, *httptest.Server) {
	t.Helper()
	server := httptest.NewTLSServer(handler)
	roots := x509.NewCertPool()
	roots.AddCert(server.Certificate())
	client, err := NewHTTPClient(server.URL, &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    roots,
	}, 5*time.Second, 1<<20)
	if err != nil {
		server.Close()
		t.Fatalf("NewHTTPClient() failed: %v", err)
	}
	return client, server
}

func requestMethod(t *testing.T, request *http.Request) string {
	t.Helper()
	body, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("ReadAll() failed: %v", err)
	}
	var envelope struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}
	return envelope.Method
}

func requireMCPHeaders(t *testing.T, request *http.Request) {
	t.Helper()
	if got, want := request.Header.Get(headerProtocolVersion), LatestProtocolVersion; got != want {
		t.Fatalf("%s = %q, want %q", headerProtocolVersion, got, want)
	}
	if got, want := request.Header.Get(headerSessionID), "session-1"; got != want {
		t.Fatalf("%s = %q, want %q", headerSessionID, got, want)
	}
}
