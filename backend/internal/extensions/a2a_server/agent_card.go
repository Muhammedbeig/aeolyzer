// Package a2aserver validates public A2A application contracts.
package a2aserver

import (
	"encoding/json"
	"errors"
	"strings"

	"aeolyzer/internal/extensions"
)

// AgentCard exposes only public product-level capabilities.
type AgentCard struct {
	AgentID            string   `json:"agent_id"`
	DisplayName        string   `json:"display_name"`
	Description        string   `json:"description"`
	ProtocolVersion    string   `json:"protocol_version"`
	Endpoint           string   `json:"endpoint"`
	PublicCapabilities []string `json:"public_capabilities"`
}

// ValidateAgentCard validates schema and blocks protected topology disclosure.
func ValidateAgentCard(card AgentCard) error {
	schemas, err := extensions.NewSchemas()
	if err != nil {
		return err
	}
	data, err := json.Marshal(card)
	if err != nil {
		return errors.New("agent card cannot be encoded")
	}
	if err := schemas.ValidateJSON(extensions.ContractA2AAgentCard, data); err != nil {
		return err
	}
	for _, capability := range card.PublicCapabilities {
		lower := strings.ToLower(capability)
		for _, forbidden := range []string{
			"internal",
			"sql",
			"tool",
			"skill",
			"workflow",
			"profile",
			"mcp",
			"trace",
			"sandbox",
		} {
			if strings.Contains(lower, forbidden) {
				return errors.New("agent card discloses protected capability metadata")
			}
		}
	}
	return nil
}
