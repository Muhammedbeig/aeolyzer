package middleware

import (
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"unicode"

	"aeolyzer/internal/intake/contracts"
	"golang.org/x/text/unicode/norm"
)

var ErrPromptInjection = errors.New("prompt injection detected")

var injectionPatterns = []string{
	"ignore all previous instructions",
	"ignore previous instructions",
	"disregard prior instructions",
	"override system instructions",
	"reveal system prompt",
	"print all instructions",
	"show developer message",
	"you are now",
	"act as the system",
	"bypass safety policy",
	"disable security",
	"exfiltrate credentials",
	"call hidden tool",
}

var encodedBlockPattern = regexp.MustCompile(`[A-Za-z0-9+/]{24,}={0,2}`)

// CheckForPromptInjection normalizes Unicode and token splitting, scans plain
// and compact forms, and recursively inspects bounded base64 fragments.
func CheckForPromptInjection(input contracts.SanitizedInput) error {
	if len(input.RawText) > 64<<10 {
		return ErrPromptInjection
	}
	if detectInjection(input.RawText) {
		return ErrPromptInjection
	}
	for _, encoded := range encodedBlockPattern.FindAllString(input.RawText, 8) {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err == nil && len(decoded) <= 16<<10 && detectInjection(string(decoded)) {
			return ErrPromptInjection
		}
	}
	return nil
}

func detectInjection(value string) bool {
	normalized := normalizeSecurityText(value)
	compact := strings.ReplaceAll(normalized, " ", "")
	for _, pattern := range injectionPatterns {
		normalizedPattern := normalizeSecurityText(pattern)
		if strings.Contains(normalized, normalizedPattern) ||
			strings.Contains(compact, strings.ReplaceAll(normalizedPattern, " ", "")) {
			return true
		}
	}
	for _, role := range []string{
		"system:",
		"developer:",
		"assistant to=functions.",
		"<system>",
		"</system>",
		"[system message]",
	} {
		if strings.Contains(strings.ToLower(value), role) {
			return true
		}
	}
	return false
}

func normalizeSecurityText(value string) string {
	value = norm.NFKC.String(strings.ToLower(value))
	var builder strings.Builder
	builder.Grow(len(value))
	space := false
	for _, character := range value {
		if unicode.IsLetter(character) || unicode.IsDigit(character) {
			builder.WriteRune(character)
			space = false
			continue
		}
		if !space {
			builder.WriteByte(' ')
			space = true
		}
	}
	return strings.TrimSpace(builder.String())
}
