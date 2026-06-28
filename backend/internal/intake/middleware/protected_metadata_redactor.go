package middleware

import (
	"errors"
	"regexp"
	"strings"

	"aeolyzer/internal/intake/contracts"
)

var ErrProtectedMetadataOutbound = errors.New("protected metadata detected in outbound response")

type Redaction struct {
	Pattern  string
	Replaced bool
}

var protectedPatterns = []struct {
	name    string
	pattern *regexp.Regexp
}{
	{"workflow_id", regexp.MustCompile(`(?i)\bworkflow[_ -]?id\b(?:\s*[:=]\s*\S+)?`)},
	{"profile_id", regexp.MustCompile(`(?i)\bprofile[_ -]?id\b(?:\s*[:=]\s*\S+)?`)},
	{"route_id", regexp.MustCompile(`(?i)\broute[_ -]?id\b(?:\s*[:=]\s*\S+)?`)},
	{"trace_id", regexp.MustCompile(`(?i)\btrace[_ -]?id\b(?:\s*[:=]\s*\S+)?`)},
	{"skill_path", regexp.MustCompile(`(?i)(?:[A-Za-z]:)?[\\/]?[^\s]*SKILL\.md|internal[\\/]skills[\\/][^\s]+`)},
	{"mcp_endpoint", regexp.MustCompile(`(?i)\bmcp://[^\s]+|\bhttps://mcp[.-][^\s]+`)},
	{"bearer_token", regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9._~+/-]{16,}=*\b`)},
	{"google_api_key", regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{20,}\b`)},
	{"github_token", regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}\b`)},
	{"private_key", regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----`)},
	{"cookie", regexp.MustCompile(`(?i)\b(?:set-cookie|cookie)\s*:\s*[^\r\n]+`)},
	{"password", regexp.MustCompile(`(?i)\b(?:password|passwd|client_secret|api_key)\s*[:=]\s*\S+`)},
}

// RedactProtectedMetadata replaces complete protected values, not only labels.
func RedactProtectedMetadata(text string) (string, []Redaction, error) {
	if len(text) > 1<<20 {
		return "", nil, errors.New("outbound response exceeds redaction limit")
	}
	var redactions []Redaction
	for _, candidate := range protectedPatterns {
		if candidate.pattern.MatchString(text) {
			text = candidate.pattern.ReplaceAllString(text, "[REDACTED]")
			redactions = append(redactions, Redaction{
				Pattern:  candidate.name,
				Replaced: true,
			})
		}
	}
	return text, redactions, nil
}

// ContainsProtectedMetadata reports whether any protected pattern is present.
func ContainsProtectedMetadata(text string) bool {
	for _, candidate := range protectedPatterns {
		if candidate.pattern.MatchString(text) {
			return true
		}
	}
	return false
}

// SafeCapabilitySummary returns product-level capability language.
func SafeCapabilitySummary(_ contracts.Intent) string {
	return "I can help with topic discovery, content briefs, source-backed research, SEO planning, page analysis, article planning, guarded drafting, optimization, repurposing, and tone preference handling."
}

func normalizeDisclosureText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
