package iam_context_test

import (
	"aeolyzer/internal/runtime/iam_context"
	"testing"
)

func TestAmbientCredentialStripping(t *testing.T) {
	env := map[string]string{
		"GITHUB_TOKEN": "secret",
	}
	if err := iam_context.StripAmbientCredentials(env); err == nil {
		t.Fatal("expected ambient credentials to be blocked")
	}
}
