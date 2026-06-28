package backendconfig

import "testing"

func TestEmbeddedConfigValidates(t *testing.T) {
	config, err := LoadEmbedded()
	if err != nil {
		t.Fatalf("LoadEmbedded() failed: %v", err)
	}
	if _, found := config.Policies.ActionClasses["canvas_write"]; !found {
		t.Fatal("action policies are missing canvas_write")
	}
}

func TestConfigRejectsOpenRouting(t *testing.T) {
	config, err := LoadEmbedded()
	if err != nil {
		t.Fatalf("LoadEmbedded() failed: %v", err)
	}
	config.Routing.AllowUserSuppliedRoutes = true
	if err := config.Validate(); err == nil {
		t.Fatal("Config.Validate() accepted user-supplied routes")
	}
}
