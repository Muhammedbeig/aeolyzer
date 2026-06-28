package security

import (
	"errors"
	"net/netip"
	"net/url"
	"strings"
)

var ErrUnsafeURL = errors.New("url is unsafe")

// SanitizeURL enforces the Layer 5 URL safety policy (Section 14.3).
// This function strictly drops any URI scheme that could execute code or access local files.
// By doing this synchronously before serialization, we guarantee that no malicious payload
// reaches the A2UI Frame renderer.
func SanitizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil ||
		parsed.Hostname() == "" ||
		parsed.User != nil {
		return "", ErrUnsafeURL
	}

	scheme := strings.ToLower(parsed.Scheme)

	// Fast path blocklist for executable or local schemas
	switch scheme {
	case "javascript", "data", "file", "blob", "chrome", "vscode", "ssh", "ftp":
		return "", ErrUnsafeURL
	case "http", "https":
		host := strings.ToLower(strings.TrimSuffix(parsed.Hostname(), "."))
		if host == "localhost" || strings.HasSuffix(host, ".localhost") {
			return "", ErrUnsafeURL
		}
		if address, err := netip.ParseAddr(host); err == nil {
			address = address.Unmap()
			if !address.IsGlobalUnicast() ||
				address.IsPrivate() ||
				address.IsLoopback() ||
				address.IsLinkLocalUnicast() ||
				address.IsMulticast() ||
				address.IsUnspecified() {
				return "", ErrUnsafeURL
			}
		}
		parsed.Scheme = scheme
		parsed.Host = strings.ToLower(parsed.Host)
		return parsed.String(), nil
	default:
		// Default deny unknown schemes
		return "", ErrUnsafeURL
	}
}
