package security

import (
	"errors"
	"net/url"
	"strings"
)

var ErrUnsafeURL = errors.New("UNSAFE_URL_SCHEME")

// SanitizeURL enforces the Layer 5 URL safety policy (Section 14.3).
// This function strictly drops any URI scheme that could execute code or access local files.
// By doing this synchronously before serialization, we guarantee that no malicious payload 
// reaches the A2UI Frame renderer.
func SanitizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrUnsafeURL
	}
	
	scheme := strings.ToLower(parsed.Scheme)
	
	// Fast path blocklist for executable or local schemas
	switch scheme {
	case "javascript", "data", "file", "blob", "chrome", "vscode", "ssh", "ftp":
		return "", ErrUnsafeURL
	case "http", "https":
		return parsed.String(), nil
	default:
		// Default deny unknown schemes
		return "", ErrUnsafeURL
	}
}
