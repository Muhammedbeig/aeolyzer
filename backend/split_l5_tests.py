import os

base = r"C:\Users\Muham\AEOlyzer\backend\internal\extensions"

# Split tests into subpackages
security_test = """package security_test

import (
	"testing"
	"aeolyzer/internal/extensions/security"
)

func TestURLSanitization(t *testing.T) {
	_, err := security.SanitizeURL("javascript:alert(1)")
	if err == nil {
		t.Fatal("expected javascript URL to be rejected")
	}
}
"""
with open(os.path.join(base, "security", "url_sanitizer_test.go"), "w") as f:
    f.write(security_test)

approval_test = """package approval_ux_test

import (
	"testing"
	"aeolyzer/internal/extensions/approval_ux"
)

func TestVibeDiffMetadataLeak(t *testing.T) {
	diff := approval_ux.VibeDiff{
		Summary: "Update constraints",
		RiskNotes: []string{"trace_id"},
	}
	if err := approval_ux.ValidateVibeDiff(diff); err == nil {
		t.Fatal("expected internal metadata leak to be blocked")
	}
}
"""
with open(os.path.join(base, "approval_ux", "vibe_diff_test.go"), "w") as f:
    f.write(approval_test)

a2a_test = """package a2a_server_test

import (
	"testing"
	"aeolyzer/internal/extensions/a2a_server"
)

func TestAgentCardDisclosure(t *testing.T) {
	card := a2a_server.AgentCard{
		PublicCapabilities: []string{"internal_sql_executor"},
	}
	if err := a2a_server.ValidateAgentCard(card); err == nil {
		t.Fatal("expected internal capability disclosure to be blocked")
	}
}
"""
with open(os.path.join(base, "a2a_server", "agent_card_test.go"), "w") as f:
    f.write(a2a_test)

events_test = """package events_test

import (
	"testing"
	"time"
	"aeolyzer/internal/extensions/events"
)

func TestStaleEventRejection(t *testing.T) {
	evt := events.UserInteractionEvent{
		Signature: "valid",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := events.NormalizeInteraction(evt); err == nil {
		t.Fatal("expected stale event to be rejected")
	}
}
"""
with open(os.path.join(base, "events", "ui_event_normalizer_test.go"), "w") as f:
    f.write(events_test)

# Remove the consolidated test file
os.remove(os.path.join(base, "tests", "layer5_details_test.go"))

print("Tests moved into subpackages.")
