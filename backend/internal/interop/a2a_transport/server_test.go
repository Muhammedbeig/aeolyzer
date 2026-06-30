package a2atransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"iter"
	"net/http"
	"net/http/httptest"
	"testing"

	"aeolyzer/internal/extensions/a2a_server"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func TestServerServesCanonicalAgentCard(t *testing.T) {
	server := newTestServer(t, fakeVerifier{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, WellKnownAgentCardPath(), nil)
	server.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("Agent Card status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var card map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &card); err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}
	if _, ok := card["supportedInterfaces"]; !ok {
		t.Fatalf("Agent Card missing supportedInterfaces: %s", recorder.Body.String())
	}
}

func TestServerRejectsUnauthenticatedJSONRPC(t *testing.T) {
	server := newTestServer(t, fakeVerifier{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/a2a", bytes.NewBufferString(`{
		"jsonrpc":"2.0",
		"id":"request-1",
		"method":"SendMessage",
		"params":{"message":{"messageId":"message-1","role":"ROLE_USER","parts":[{"text":"hello"}]}}
	}`))
	server.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("A2A JSON-RPC status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"error"`)) {
		t.Fatalf("A2A JSON-RPC response did not contain error: %s", recorder.Body.String())
	}
}

func TestServerAllowsAuthenticatedJSONRPC(t *testing.T) {
	server := newTestServer(t, fakeVerifier{subject: "caller-1"})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/a2a", bytes.NewBufferString(`{
		"jsonrpc":"2.0",
		"id":"request-1",
		"method":"SendMessage",
		"params":{"message":{"messageId":"message-1","role":"ROLE_USER","parts":[{"text":"hello"}]}}
	}`))
	request.Header.Set("Authorization", "Bearer good")
	server.Routes().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("A2A JSON-RPC status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if bytes.Contains(recorder.Body.Bytes(), []byte(`"error"`)) {
		t.Fatalf("A2A JSON-RPC response contained error: %s", recorder.Body.String())
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte("AEOlyzer received a safe public request")) {
		t.Fatalf("A2A JSON-RPC response missing agent text: %s", recorder.Body.String())
	}
}

func TestServerRequiresVerifier(t *testing.T) {
	adkAgent := newADKAgent(t)
	card, err := a2aserver.NewAgentCard(a2aserver.AgentCardConfig{
		Name:          "AEOlyzer",
		Description:   "Provides guarded website visibility and content-planning capabilities.",
		PublicBaseURL: "https://api.aeolyzer.example",
		Skills:        a2aserver.DefaultPublicSkills(),
	})
	if err != nil {
		t.Fatalf("NewAgentCard() failed: %v", err)
	}
	if _, err := NewServer(Config{Agent: adkAgent, Card: card}); err == nil {
		t.Fatal("NewServer() accepted nil verifier")
	}
}

func newTestServer(t *testing.T, verifier TokenVerifier) *Server {
	t.Helper()
	adkAgent := newADKAgent(t)
	card, err := a2aserver.NewAgentCard(a2aserver.AgentCardConfig{
		Name:          "AEOlyzer",
		Description:   "Provides guarded website visibility and content-planning capabilities.",
		PublicBaseURL: "https://api.aeolyzer.example",
		Skills:        a2aserver.DefaultPublicSkills(),
	})
	if err != nil {
		t.Fatalf("NewAgentCard() failed: %v", err)
	}
	server, err := NewServer(Config{
		Agent:    adkAgent,
		Card:     card,
		Verifier: verifier,
	})
	if err != nil {
		t.Fatalf("NewServer() failed: %v", err)
	}
	return server
}

func newADKAgent(t *testing.T) agent.Agent {
	t.Helper()
	adkAgent, err := agent.New(agent.Config{
		Name:        "public_site_guidance",
		Description: "Explains safe public website visibility options.",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			return func(yield func(*session.Event, error) bool) {
				event := session.NewEvent(ctx.InvocationID())
				event.Author = ctx.Agent().Name()
				event.LLMResponse = model.LLMResponse{
					Content: genai.NewContentFromText("AEOlyzer received a safe public request.", genai.RoleModel),
				}
				yield(event, nil)
			}
		},
	})
	if err != nil {
		t.Fatalf("agent.New() failed: %v", err)
	}
	return adkAgent
}

type fakeVerifier struct {
	subject string
	err     error
}

func (v fakeVerifier) VerifyBearer(_ context.Context, token string) (Identity, error) {
	if v.err != nil {
		return Identity{}, v.err
	}
	if token != "good" {
		return Identity{}, errors.New("bad token")
	}
	subject := v.subject
	if subject == "" {
		subject = "caller"
	}
	return Identity{Subject: subject, Tenant: "tenant-1"}, nil
}
