// Package iamcontext enforces zero-ambient authority and scoped credentials.
package iamcontext

import (
	"errors"
	"strings"
)

var (
	// ErrAmbientCredentialsBlocked indicates that the caller attempted to pass
	// a credential-bearing environment variable into a runtime.
	ErrAmbientCredentialsBlocked = errors.New("ambient credentials are blocked")
	// ErrEnvironmentVariableDenied indicates that a variable was not explicitly
	// allowlisted for the sandbox.
	ErrEnvironmentVariableDenied = errors.New("environment variable is not allowlisted")
)

// SanitizeEnvironment returns only explicitly allowlisted, non-secret
// environment variables. Keys are compared case-insensitively.
func SanitizeEnvironment(
	env map[string]string,
	allowedKeys []string,
) (map[string]string, error) {
	allowed := make(map[string]struct{}, len(allowedKeys))
	for _, key := range allowedKeys {
		normalized := strings.ToUpper(strings.TrimSpace(key))
		if normalized == "" || secretEnvironmentKey(normalized) {
			return nil, ErrEnvironmentVariableDenied
		}
		allowed[normalized] = struct{}{}
	}

	result := make(map[string]string, len(env))
	for key, value := range env {
		normalized := strings.ToUpper(strings.TrimSpace(key))
		if secretEnvironmentKey(normalized) {
			return nil, ErrAmbientCredentialsBlocked
		}
		if _, ok := allowed[normalized]; !ok {
			return nil, ErrEnvironmentVariableDenied
		}
		result[normalized] = value
	}
	return result, nil
}

// StripAmbientCredentials rejects any credential-like environment. It remains
// for compatibility with the earlier Layer 6 boundary API.
func StripAmbientCredentials(env map[string]string) error {
	_, err := SanitizeEnvironment(env, nonSecretKeys(env))
	return err
}

func nonSecretKeys(env map[string]string) []string {
	result := make([]string, 0, len(env))
	for key := range env {
		if !secretEnvironmentKey(strings.ToUpper(strings.TrimSpace(key))) {
			result = append(result, key)
		}
	}
	return result
}

func secretEnvironmentKey(key string) bool {
	for _, fragment := range []string{
		"ACCESS_KEY",
		"API_KEY",
		"APIKEY",
		"AUTHORIZATION",
		"BEARER",
		"CLIENT_SECRET",
		"COOKIE",
		"CREDENTIAL",
		"GH_TOKEN",
		"GITHUB_TOKEN",
		"OAUTH",
		"PASSWORD",
		"PASSWD",
		"PRIVATE_KEY",
		"REFRESH_TOKEN",
		"SECRET",
		"SESSION_TOKEN",
	} {
		if strings.Contains(key, fragment) {
			return true
		}
	}
	return false
}
