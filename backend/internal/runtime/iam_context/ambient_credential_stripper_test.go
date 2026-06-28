package iamcontext

import (
	"errors"
	"testing"
)

func TestSanitizeEnvironmentRejectsCredentialVariants(t *testing.T) {
	for _, key := range []string{
		"AWS_ACCESS_KEY_ID",
		"GITHUB_TOKEN",
		"gemini_api_key",
		"MY_CLIENT_SECRET",
		"DATABASE_PASSWORD",
		"authorization",
	} {
		t.Run(key, func(t *testing.T) {
			_, err := SanitizeEnvironment(
				map[string]string{key: "secret"},
				[]string{key},
			)
			if !errors.Is(err, ErrEnvironmentVariableDenied) &&
				!errors.Is(err, ErrAmbientCredentialsBlocked) {
				t.Fatalf("SanitizeEnvironment() error = %v, want credential denial", err)
			}
		})
	}
}

func TestSanitizeEnvironmentRequiresExplicitAllowlist(t *testing.T) {
	_, err := SanitizeEnvironment(
		map[string]string{"LANG": "en_US.UTF-8", "PATH": "/usr/bin"},
		[]string{"LANG"},
	)
	if !errors.Is(err, ErrEnvironmentVariableDenied) {
		t.Fatalf("SanitizeEnvironment() error = %v, want %v", err, ErrEnvironmentVariableDenied)
	}
}

func TestSanitizeEnvironmentReturnsDefensiveCopy(t *testing.T) {
	env := map[string]string{"LANG": "en_US.UTF-8"}
	result, err := SanitizeEnvironment(env, []string{"LANG"})
	if err != nil {
		t.Fatalf("SanitizeEnvironment() failed: %v", err)
	}
	env["LANG"] = "changed"
	if result["LANG"] != "en_US.UTF-8" {
		t.Fatal("SanitizeEnvironment() returned caller-owned state")
	}
}

func TestStripAmbientCredentialsRejectsSecret(t *testing.T) {
	if err := StripAmbientCredentials(
		map[string]string{"AWS_ACCESS_KEY_ID": "secret"},
	); !errors.Is(err, ErrAmbientCredentialsBlocked) {
		t.Fatalf("StripAmbientCredentials() error = %v, want %v", err, ErrAmbientCredentialsBlocked)
	}
}
