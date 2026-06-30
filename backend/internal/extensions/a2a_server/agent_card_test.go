package a2aserver

import (
	"encoding/json"
	"testing"

	"github.com/a2aproject/a2a-go/v2/a2a"
	"google.golang.org/adk/agent"
)

func TestNewAgentCardUsesCanonicalA2AFields(t *testing.T) {
	card, err := NewAgentCard(AgentCardConfig{
		Name:          "AEOlyzer",
		Description:   "Provides guarded website visibility and content-planning capabilities.",
		PublicBaseURL: "https://api.aeolyzer.example",
		Skills:        DefaultPublicSkills(),
	})
	if err != nil {
		t.Fatalf("NewAgentCard() failed: %v", err)
	}
	data, err := agentCardJSON(card)
	if err != nil {
		t.Fatalf("agentCardJSON() failed: %v", err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}
	for _, key := range []string{"supportedInterfaces", "defaultInputModes", "defaultOutputModes", "securityRequirements"} {
		if _, ok := raw[key]; !ok {
			t.Fatalf("Agent Card missing canonical key %q in %s", key, data)
		}
	}
	for _, forbidden := range []string{"agent_id", "display_name", "public_capabilities"} {
		if _, ok := raw[forbidden]; ok {
			t.Fatalf("Agent Card exposed obsolete key %q in %s", forbidden, data)
		}
	}
}

func TestAgentCardDisclosure(t *testing.T) {
	card, err := NewAgentCard(AgentCardConfig{
		Name:          "AEOlyzer",
		Description:   "Provides guarded website visibility and content-planning capabilities.",
		PublicBaseURL: "https://api.aeolyzer.example",
		Skills: []a2a.AgentSkill{{
			ID:          "site_visibility_guidance",
			Name:        "Site visibility guidance",
			Description: "Calls internal_sql_executor through an mcp server for hidden workflow ids.",
			Tags:        []string{"aeo"},
		}},
	})
	if err == nil {
		t.Fatalf("NewAgentCard() accepted internal disclosure: %#v", card)
	}
}

func TestAgentCardRequiresHTTPSAndA2AProtocol(t *testing.T) {
	card := &a2a.AgentCard{
		Name:                "AEOlyzer",
		Description:         "Provides guarded website visibility and content-planning capabilities.",
		SupportedInterfaces: []*a2a.AgentInterface{a2a.NewAgentInterface("http://api.aeolyzer.example/a2a", a2a.TransportProtocolJSONRPC)},
		Capabilities:        a2a.AgentCapabilities{},
		DefaultInputModes:   []string{"text/plain"},
		DefaultOutputModes:  []string{"text/plain"},
		Skills:              DefaultPublicSkills(),
		Version:             "1.0.0",
	}
	if err := ValidateAgentCard(card); err == nil {
		t.Fatal("ValidateAgentCard() accepted non-HTTPS A2A endpoint")
	}
	card.SupportedInterfaces[0] = a2a.NewAgentInterface("https://api.aeolyzer.example/a2a", a2a.TransportProtocolGRPC)
	if err := ValidateAgentCard(card); err == nil {
		t.Fatal("ValidateAgentCard() accepted non-JSON-RPC protocol binding")
	}
}

func TestBuildAgentSkillsUsesADKMetadata(t *testing.T) {
	adkAgent, err := agent.New(agent.Config{
		Name:        "public_site_guidance",
		Description: "Explains safe public website visibility options.",
	})
	if err != nil {
		t.Fatalf("agent.New() failed: %v", err)
	}
	skills := BuildAgentSkills(adkAgent)
	if len(skills) != 1 {
		t.Fatalf("BuildAgentSkills() returned %d skills, want 1", len(skills))
	}
	if skills[0].ID == "" || skills[0].Name == "" {
		t.Fatalf("BuildAgentSkills()[0] missing public identifiers: %#v", skills[0])
	}
	if got, want := skills[0].Description, "Explains safe public website visibility options."; got != want {
		t.Fatalf("BuildAgentSkills()[0].Description = %q, want %q", got, want)
	}
}
