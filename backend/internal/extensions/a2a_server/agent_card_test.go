package a2a_server_test

import (
	"aeolyzer/internal/extensions/a2a_server"
	"testing"
)

func TestAgentCardDisclosure(t *testing.T) {
	card := a2a_server.AgentCard{
		PublicCapabilities: []string{"internal_sql_executor"},
	}
	if err := a2a_server.ValidateAgentCard(card); err == nil {
		t.Fatal("expected internal capability disclosure to be blocked")
	}
}
