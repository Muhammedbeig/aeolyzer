package security

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

const maxMarkdownBytes = 64 << 10

var (
	rawHTMLPattern = regexp.MustCompile(`(?i)<\s*/?\s*[a-z][^>]*>`)
	linkPattern    = regexp.MustCompile(`!?\[[^\]]*\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
)

// SanitizeMarkdown validates conservative user-facing Markdown. Unsafe input
// is rejected rather than rewritten ambiguously.
func SanitizeMarkdown(markdown string) (string, error) {
	if markdown == "" {
		return "", errors.New("markdown is empty")
	}
	if len(markdown) > maxMarkdownBytes || !utf8.ValidString(markdown) {
		return "", errors.New("markdown size or encoding is invalid")
	}
	if rawHTMLPattern.MatchString(markdown) ||
		strings.Contains(markdown, "<!--") ||
		strings.Contains(markdown, "-->") {
		return "", errors.New("raw html is not allowed in markdown")
	}
	for _, character := range markdown {
		if unicode.IsControl(character) &&
			character != '\n' &&
			character != '\r' &&
			character != '\t' {
			return "", errors.New("markdown contains control characters")
		}
	}
	matches := linkPattern.FindAllStringSubmatch(markdown, -1)
	for _, match := range matches {
		if len(match) != 2 {
			return "", errors.New("markdown link is malformed")
		}
		if _, err := SanitizeURL(match[1]); err != nil {
			return "", errors.New("markdown contains an unsafe link")
		}
	}
	for _, pattern := range hiddenPatterns {
		if pattern.pattern.MatchString(markdown) {
			return "", errors.New("markdown contains hidden or executable content")
		}
	}
	return strings.ReplaceAll(markdown, "\r\n", "\n"), nil
}
