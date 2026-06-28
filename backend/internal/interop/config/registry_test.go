package interopconfig

import "testing"

func TestEmbeddedLayer7ConfigValidates(t *testing.T) {
	config, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	if len(config.Registry.Connectors) != 1 {
		t.Fatalf("connectors = %d, want 1", len(config.Registry.Connectors))
	}
	if err := config.MCPManifest.ValidateJSON([]byte("{}")); err == nil {
		t.Fatal("MCP manifest schema accepted placeholder")
	}
}

func TestLayer7ConfigRejectsUnsafeConnector(t *testing.T) {
	config, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}
	config.Registry.Connectors[0].RedirectsAllowed = true
	if err := config.Validate(); err == nil {
		t.Fatal("Config.Validate() accepted redirects")
	}
}
