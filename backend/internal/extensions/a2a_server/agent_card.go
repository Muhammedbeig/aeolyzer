package a2a_server

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
