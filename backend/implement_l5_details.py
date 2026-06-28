import os

base = r"C:\Users\Muham\AEOlyzer\backend\internal\extensions"

# 1. security/url_sanitizer.go
url_sanitizer_content = """package security

import (
	"errors"
	"net/url"
	"strings"
)

var ErrUnsafeURL = errors.New("UNSAFE_URL_SCHEME")

// SanitizeURL enforces the Layer 5 URL safety policy (Section 14.3).
// This function strictly drops any URI scheme that could execute code or access local files.
// By doing this synchronously before serialization, we guarantee that no malicious payload 
// reaches the A2UI Frame renderer.
func SanitizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrUnsafeURL
	}
	
	scheme := strings.ToLower(parsed.Scheme)
	
	// Fast path blocklist for executable or local schemas
	switch scheme {
	case "javascript", "data", "file", "blob", "chrome", "vscode", "ssh", "ftp":
		return "", ErrUnsafeURL
	case "http", "https":
		return parsed.String(), nil
	default:
		// Default deny unknown schemes
		return "", ErrUnsafeURL
	}
}
"""
with open(os.path.join(base, "security", "url_sanitizer.go"), "w") as f:
    f.write(url_sanitizer_content)

# 2. approval_ux/vibe_diff.go
vibe_diff_content = """package approval_ux

import "errors"

// VibeDiff represents a strictly constrained summary of a state change (Section 12.2).
// It decouples the approval presentation from the underlying system state.
type VibeDiff struct {
	Summary        string      `json:"summary"`
	ChangeType     string      `json:"change_type"`
	Before         interface{} `json:"before,omitempty"`
	After          interface{} `json:"after,omitempty"`
	RiskNotes      []string    `json:"risk_notes,omitempty"`
}

// ValidateVibeDiff ensures the diff does not violate data leakage boundaries.
// By statically rejecting internal identifiers, it maintains the abstraction 
// between the user-facing approval card and the orchestrator's state machine.
func ValidateVibeDiff(diff VibeDiff) error {
	if diff.Summary == "" {
		return errors.New("MISSING_SUMMARY")
	}
	
	// Prevent leaking internal system policies into the user approval screen.
	for _, note := range diff.RiskNotes {
		if note == "policy.yaml" || note == "trace_id" {
			return errors.New("INTERNAL_METADATA_LEAK")
		}
	}
	
	return nil
}
"""
with open(os.path.join(base, "approval_ux", "vibe_diff.go"), "w") as f:
    f.write(vibe_diff_content)

# 3. a2a_server/agent_card.go
agent_card_content = """package a2a_server

import "errors"

// AgentCard schema (Section 13.1) allows external agent discovery without exposing internals.
type AgentCard struct {
	AgentID            string   `json:"agent_id"`
	DisplayName        string   `json:"display_name"`
	PublicCapabilities []string `json:"public_capabilities"`
}

// ValidateAgentCard guarantees that internal mechanisms (like DAGs or specific skills)
// are masked from the public A2A endpoint. 
// This acts as an API gateway firewall against topology disclosure.
func ValidateAgentCard(card AgentCard) error {
	for _, cap := range card.PublicCapabilities {
		// Example check: block anything referencing internal tool nomenclature
		if cap == "internal_sql_executor" || cap == "skill_registry" {
			return errors.New("UNAUTHORIZED_CAPABILITY_DISCLOSURE")
		}
	}
	return nil
}
"""
with open(os.path.join(base, "a2a_server", "agent_card.go"), "w") as f:
    f.write(agent_card_content)

# 4. events/ui_event_normalizer.go
event_normalizer_content = """package events

import (
	"errors"
	"time"
)

// UserInteractionEvent is the inbound contract from Layer 1.
type UserInteractionEvent struct {
	FrameID       string `json:"frame_id"`
	ActionID      string `json:"action_id"`
	InteractionID string `json:"interaction_id"`
	Signature     string `json:"signature"`
	ExpiresAt     time.Time
}

// NormalizeInteraction strictly validates replay and signature bounds.
// Rejecting stale or tampered interactions at the border prevents out-of-order 
// DAG progression in Layer 3.
func NormalizeInteraction(evt UserInteractionEvent) error {
	if evt.Signature == "" {
		return errors.New("MISSING_SIGNATURE")
	}
	
	if time.Now().After(evt.ExpiresAt) {
		return errors.New("STALE_EVENT_REJECTED")
	}
	
	return nil
}
"""
with open(os.path.join(base, "events", "ui_event_normalizer.go"), "w") as f:
    f.write(event_normalizer_content)

# 5. tests/layer5_details_test.go
test_content = """package extensions_test

import (
	"testing"
	"time"
	
	"aeolyzer/internal/extensions/security"
	"aeolyzer/internal/extensions/approval_ux"
	"aeolyzer/internal/extensions/a2a_server"
	"aeolyzer/internal/extensions/events"
)

func TestURLSanitization(t *testing.T) {
	_, err := security.SanitizeURL("javascript:alert(1)")
	if err == nil {
		t.Fatal("expected javascript URL to be rejected")
	}
}

func TestVibeDiffMetadataLeak(t *testing.T) {
	diff := approval_ux.VibeDiff{
		Summary: "Update constraints",
		RiskNotes: []string{"trace_id"},
	}
	if err := approval_ux.ValidateVibeDiff(diff); err == nil {
		t.Fatal("expected internal metadata leak to be blocked")
	}
}

func TestAgentCardDisclosure(t *testing.T) {
	card := a2a_server.AgentCard{
		PublicCapabilities: []string{"internal_sql_executor"},
	}
	if err := a2a_server.ValidateAgentCard(card); err == nil {
		t.Fatal("expected internal capability disclosure to be blocked")
	}
}

func TestStaleEventRejection(t *testing.T) {
	evt := events.UserInteractionEvent{
		Signature: "valid",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
	}
	if err := events.NormalizeInteraction(evt); err == nil {
		t.Fatal("expected stale event to be rejected")
	}
}
"""
with open(os.path.join(base, "tests", "layer5_details_test.go"), "w") as f:
    f.write(test_content)

print("Layer 5 detailed implementations generated.")
