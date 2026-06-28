package extensions

import "testing"

func TestLayer5SchemasCompileAndCatalogValidates(t *testing.T) {
	if _, err := NewSchemas(); err != nil {
		t.Fatalf("NewSchemas() failed: %v", err)
	}
}

func TestLayer5SchemasRejectPlaceholders(t *testing.T) {
	schemas, err := NewSchemas()
	if err != nil {
		t.Fatalf("NewSchemas() failed: %v", err)
	}
	for _, contract := range []Contract{
		ContractPresentation,
		ContractA2UIFrame,
		ContractA2UICatalog,
		ContractUIEvent,
		ContractApproval,
		ContractSurfacePatch,
		ContractA2AAgentCard,
		ContractA2AEnvelope,
	} {
		t.Run(string(contract), func(t *testing.T) {
			if err := schemas.ValidateJSON(contract, []byte("{}")); err == nil {
				t.Fatal("placeholder unexpectedly validated")
			}
		})
	}
}
