package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRoutingSchema(t *testing.T) {
	path := filepath.Join("..", "config", "routing-schema.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read routing-schema.json: %v", err)
	}

	var schema struct {
		Version        int      `json:"version"`
		AllowedIntents []string `json:"allowed_intents"`
	}

	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}

	if schema.Version != 2 {
		t.Errorf("Expected version 2, got %d", schema.Version)
	}

	foundTopicDiscovery := false
	for _, intent := range schema.AllowedIntents {
		if intent == "topic_discovery" {
			foundTopicDiscovery = true
			break
		}
	}

	if !foundTopicDiscovery {
		t.Errorf("Expected topic_discovery in allowed intents")
	}
}
