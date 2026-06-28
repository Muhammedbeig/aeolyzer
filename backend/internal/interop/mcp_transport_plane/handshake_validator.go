package mcptransportplane

import (
	"errors"
	"fmt"
	"sort"
)

// ToolManifest describes one MCP tool exposed during handshake.
type ToolManifest struct {
	Name       string `json:"name"`
	SchemaHash string `json:"schema_hash"`
}

// ServerManifest is the validated handshake surface.
type ServerManifest struct {
	ServerID        string         `json:"server_id"`
	ProtocolVersion string         `json:"protocol_version"`
	Capabilities    []string       `json:"capabilities"`
	Tools           []ToolManifest `json:"tools"`
}

// PinnedManifest is the registry-approved server surface.
type PinnedManifest struct {
	ServerID                string
	AllowedProtocolVersions []string
	RequiredCapabilities    []string
	Tools                   []ToolManifest
}

// ValidateHandshake rejects identity, protocol, capability, tool-list, and
// schema-hash drift.
func ValidateHandshake(live ServerManifest, pinned PinnedManifest) error {
	if live.ServerID == "" ||
		pinned.ServerID == "" ||
		live.ServerID != pinned.ServerID ||
		!containsString(pinned.AllowedProtocolVersions, live.ProtocolVersion) {
		return errors.New("mcp handshake identity or protocol mismatch")
	}
	liveCapabilities := append([]string(nil), live.Capabilities...)
	requiredCapabilities := append([]string(nil), pinned.RequiredCapabilities...)
	sort.Strings(liveCapabilities)
	sort.Strings(requiredCapabilities)
	if !equalStrings(liveCapabilities, requiredCapabilities) {
		return errors.New("mcp capability set drift")
	}
	if len(live.Tools) != len(pinned.Tools) {
		return errors.New("mcp tool list drift")
	}
	expected := make(map[string]string, len(pinned.Tools))
	for _, tool := range pinned.Tools {
		if tool.Name == "" || tool.SchemaHash == "" {
			return errors.New("pinned mcp tool manifest is invalid")
		}
		if _, duplicate := expected[tool.Name]; duplicate {
			return errors.New("pinned mcp tool manifest contains duplicates")
		}
		expected[tool.Name] = tool.SchemaHash
	}
	for _, tool := range live.Tools {
		hash, found := expected[tool.Name]
		if !found || hash != tool.SchemaHash {
			return fmt.Errorf("mcp tool %q schema or registration drift", tool.Name)
		}
		delete(expected, tool.Name)
	}
	if len(expected) != 0 {
		return errors.New("mcp handshake omitted pinned tools")
	}
	return nil
}

func containsString(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
