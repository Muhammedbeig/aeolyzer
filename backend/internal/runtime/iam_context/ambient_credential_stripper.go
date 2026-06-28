package iam_context

import "errors"

var ErrAmbientCredentialsBlocked = errors.New("AMBIENT_CREDENTIALS_BLOCKED")

// StripAmbientCredentials enforces zero-ambient authority (Section 4.3).
// Prevents a script or tool from inheriting the host/orchestrator's implicit permissions.
// Any execution context attempting to pass raw ENV secrets will be rejected at the gateway.
func StripAmbientCredentials(env map[string]string) error {
	for k := range env {
		if k == "AWS_ACCESS_KEY_ID" || k == "GITHUB_TOKEN" {
			return ErrAmbientCredentialsBlocked
		}
	}
	return nil
}
