package httpapi

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"aeolyzer/layer_02_intake"
	"aeolyzer/layer_03_orchestration"
	"aeolyzer/layer_06_runtime"
	"aeolyzer/layer_08_observability"
)

type testResolver struct{}

func (testResolver) LookupIPAddr(_ context.Context, _ string) ([]net.IPAddr, error) {
	return []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}, nil
}

type testAdapter struct{}

func (testAdapter) Inspect(_ context.Context, _ string, _ int64) (runtime.ExecutionResult, error) {
	return runtime.ExecutionResult{
		Title:       "AEOlyzer",
		Description: "Answer engine visibility",
		IconURL:     "https://example.com/favicon.ico",
	}, nil
}

func TestCompleteOnboarding(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)
	body := `{
		"session_id":"guest-1",
		"account_type":"brand",
		"domain":"https://example.com",
		"brand_name":"AEOlyzer",
		"reach":"global",
		"country_code":"PK",
		"country_name":"Pakistan",
		"language":"English (UK)",
		"competitors":["peer.example"]
	}`
	request := httptest.NewRequest(http.MethodPost, "/v1/onboarding/complete", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("complete onboarding status = %d, body = %s", response.Code, response.Body.String())
	}
	var frame struct {
		Surface string   `json:"surface"`
		Prompts []string `json:"prompts"`
	}
	if err := json.NewDecoder(response.Body).Decode(&frame); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if frame.Surface != "audit_dashboard" {
		t.Fatalf("surface = %q", frame.Surface)
	}
	if len(frame.Prompts) != 12 {
		t.Fatalf("prompts = %d, want 12", len(frame.Prompts))
	}
}

func TestRejectsUnknownJSONField(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)
	request := httptest.NewRequest(
		http.MethodPost,
		"/v1/onboarding/inspect",
		bytes.NewBufferString(`{"session_id":"guest-1","url":"example.com","secret":"no"}`),
	)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("inspect status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestRejectsUnapprovedOrigin(t *testing.T) {
	t.Parallel()

	handler := newTestHandler(t)
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set("Origin", "https://untrusted.example")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("health status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()

	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		t.Fatalf("generate key: %v", err)
	}
	now := func() time.Time { return time.Unix(100, 0) }
	intakeService := intake.NewService(func() string { return "trace-1" }, key, now)
	executor := runtime.NewExecutor(testResolver{}, testAdapter{}, key, now)
	handler := NewHandler(
		intakeService,
		orchestration.NewService(),
		executor,
		observability.NewSink(10),
		slog.Default(),
		"http://localhost:3000",
	)
	handler.now = now

	return handler.Routes()
}
