package runtime

import "testing"

func TestLayer6SchemasCompile(t *testing.T) {
	if err := CompileSchemas(); err != nil {
		t.Fatalf("CompileSchemas() failed: %v", err)
	}
}
