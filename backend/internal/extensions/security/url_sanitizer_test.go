package security_test

import (
	"aeolyzer/internal/extensions/security"
	"testing"
)

func TestURLSanitization(t *testing.T) {
	_, err := security.SanitizeURL("javascript:alert(1)")
	if err == nil {
		t.Fatal("expected javascript URL to be rejected")
	}
}
