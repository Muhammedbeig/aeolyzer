package a2aserver

import "testing"

func TestAgentCardDisclosure(t *testing.T) {
	card := AgentCard{
		AgentID:         "aeolyzer",
		DisplayName:     "AEOlyzer",
		Description:     "Audits websites and prepares safe content recommendations.",
		ProtocolVersion: "1.0",
		Endpoint:        "https://api.aeolyzer.example/a2a",
		PublicCapabilities: []string{
			"internal_sql_executor",
		},
	}
	if err := ValidateAgentCard(card); err == nil {
		t.Fatal("ValidateAgentCard() accepted internal capability disclosure")
	}
	card.PublicCapabilities = []string{"website_audit"}
	if err := ValidateAgentCard(card); err != nil {
		t.Fatalf("ValidateAgentCard(safe) failed: %v", err)
	}
}
