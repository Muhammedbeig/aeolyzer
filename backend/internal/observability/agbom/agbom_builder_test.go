package agbom_test

import (
	"testing"
	"aeolyzer/internal/observability/agbom"
)

func TestBuildAgBOM(t *testing.T) {
	_, err := agbom.BuildAgBOM("")
	if err == nil {
		t.Fatal("expected trace ID requirement for agbom")
	}
}
