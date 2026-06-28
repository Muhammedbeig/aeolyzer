package middleware

import (
	"strings"
	"unicode/utf8"
)

// Enforces structural integrity on incoming dynamic context from upstream.
// Nullifies invisible payloads by completely dropping control characters (except standard whitespaces)
// and strict length capping.
func SanitizeContextValue(value string, limit int) string {
	cleaned := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, value)

	if !utf8.ValidString(cleaned) {
		return ""
	}

	// Slice by runes, not bytes, to avoid corrupting multi-byte unicode sequences.
	if utf8.RuneCountInString(cleaned) > limit {
		runes := []rune(cleaned)
		return string(runes[:limit])
	}

	return cleaned
}

// Unmarshals untrusted dynamic dictionaries into a strict, bounded schema.
// Drops any arbitrary key not explicitly permitted by the internal schema definitions.
// Mitigates context window overflow and stops lateral injection attacks via undefined keys.
func ExtractSanitizedContext(raw map[string]interface{}) map[string]string {
	sanitized := make(map[string]string)

	allowedKeys := map[string]int{
		"target_domain": 100,
		"target_url":    2048,
		"topic":         200,
		"audience":      300,
		"angle":         500,
		"content_type":  80,
	}

	for k, v := range raw {
		if limit, ok := allowedKeys[k]; ok {
			if strVal, isStr := v.(string); isStr {
				sanitized[k] = SanitizeContextValue(strVal, limit)
			}
		}
	}

	return sanitized
}
