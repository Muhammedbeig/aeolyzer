package a2uitranslator

import (
	"testing"

	"aeolyzer/internal/extensions"
)

func TestSchemaAndCatalogRejectCyclesAndExecutableProps(t *testing.T) {
	schemas, err := extensions.NewSchemas()
	if err != nil {
		t.Fatalf("extensions.NewSchemas() failed: %v", err)
	}
	manager, err := NewSchemaManager(schemas)
	if err != nil {
		t.Fatalf("NewSchemaManager() failed: %v", err)
	}
	catalog, err := LoadCatalog(extensions.CatalogLock(), schemas)
	if err != nil {
		t.Fatalf("LoadCatalog() failed: %v", err)
	}
	frame := extensions.A2UIFrame{
		FrameID:        "frame-1",
		Surface:        "audit_dashboard",
		CatalogID:      catalog.CatalogID,
		CatalogVersion: catalog.CatalogVersion,
		SchemaVersion:  catalog.SchemaVersion,
		RootID:         "root",
		Nodes: []extensions.A2UINode{{
			ID:       "root",
			Type:     "Container",
			Props:    map[string]any{"layout": "column"},
			Children: []string{"text"},
		}, {
			ID:    "text",
			Type:  "Text",
			Props: map[string]any{"text": "Safe content"},
		}},
		Signature: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ",
	}
	if err := manager.ValidateFrame(frame); err != nil {
		t.Fatalf("SchemaManager.ValidateFrame() failed: %v", err)
	}
	if err := ValidateCatalogNodes(frame.Nodes, catalog); err != nil {
		t.Fatalf("ValidateCatalogNodes() failed: %v", err)
	}

	frame.Nodes[1].Children = []string{"root"}
	if err := manager.ValidateFrame(frame); err == nil {
		t.Fatal("SchemaManager.ValidateFrame() accepted cycle")
	}
	frame.Nodes[1].Children = nil
	frame.Nodes[1].Props["onClick"] = "javascript:alert(1)"
	if err := ValidateCatalogNodes(frame.Nodes, catalog); err == nil {
		t.Fatal("ValidateCatalogNodes() accepted executable prop")
	}
}
