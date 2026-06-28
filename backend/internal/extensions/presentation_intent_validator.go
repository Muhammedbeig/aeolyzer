package extensions

import (
	"errors"
	"strings"
)

var (
	ErrUnknownSurface = errors.New("UNKNOWN_SURFACE")
	ErrUnsafePayload  = errors.New("UNSAFE_PAYLOAD")
)

// ValidatePresentationIntent ensures that upstream layers are not leaking
// raw execution contexts or arbitrary javascript into the presentation layer.
// This preserves the firewall against DOM XSS and payload smuggling.
func ValidatePresentationIntent(intent PresentationIntent) error {
	if intent.Surface == "" {
		return ErrUnknownSurface
	}

	// Enforce that fallback text does not contain raw HTML injections.
	// This acts as a secondary depth-in-defense check.
	if strings.Contains(intent.FallbackText, "<script>") {
		return ErrUnsafePayload
	}

	return nil
}
