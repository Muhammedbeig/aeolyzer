package security_test

import (
	"testing"
	"aeolyzer/internal/extensions/security"
)

func TestURLSanitization(t *testing.T) {
	_, err := security.SanitizeURL("javascript:alert(1)")
	if err == nil {
		t.Fatal("expected javascript URL to be rejected")
	}
}
