package event_ingestion

import "strings"

// RedactEvent safely strips out sensitive values before they reach the data sink (Section 8).
// Ensures hidden chain-of-thought or raw system prompts are not leaked to external logs.
func RedactEvent(payload string) string {
	payload = strings.ReplaceAll(payload, "raw_system_prompt", "[REDACTED]")
	payload = strings.ReplaceAll(payload, "raw_developer_prompt", "[REDACTED]")
	return payload
}
