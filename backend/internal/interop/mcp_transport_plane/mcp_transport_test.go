package mcptransportplane

import (
	"encoding/json"
	"errors"
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
}

func TestValidateHandshakeRejectsExtraToolAndSchemaDrift(t *testing.T) {
	pinned := PinnedManifest{
		ServerID:                "gsc",
		AllowedProtocolVersions: []string{"2025-03-26"},
		RequiredCapabilities:    []string{"tools"},
		Tools: []ToolManifest{{
			Name:       "query",
			SchemaHash: "sha256:abc",
		}},
	}
	live := ServerManifest{
		ServerID:        "gsc",
		ProtocolVersion: "2025-03-26",
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
