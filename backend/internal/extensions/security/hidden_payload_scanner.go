package security

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	maxPayloadDepth = 32
	maxPayloadKeys  = 10_000
)

var hiddenPatterns = []struct {
	class   string
	pattern *regexp.Regexp
}{
	{"script", regexp.MustCompile(`(?i)<\s*(script|iframe|object|embed|svg|math)\b`)},
	{"event_handler", regexp.MustCompile(`(?i)\bon[a-z]{3,}\s*=`)},
	{"executable_url", regexp.MustCompile(`(?i)(javascript|vbscript|data\s*:\s*text/html)\s*:`)},
	{"html_comment", regexp.MustCompile(`(?s)<!--.*?-->`)},
	{"css_hiding", regexp.MustCompile(`(?i)(display\s*:\s*none|visibility\s*:\s*hidden|opacity\s*:\s*0(?:\D|$))`)},
	{"instruction_override", regexp.MustCompile(`(?i)\b(ignore|override|disregard)\b.{0,80}\b(previous|system|developer)\b.{0,40}\b(instruction|prompt|message)s?\b`)},
}

// HiddenPayloadFinding identifies unsafe payload structure without returning
// the raw value.
type HiddenPayloadFinding struct {
	Path  string `json:"path"`
	Class string `json:"class"`
}

// ScanHiddenPayload recursively scans a declarative payload.
func ScanHiddenPayload(payload any) ([]HiddenPayloadFinding, error) {
	keys := 0
	var findings []HiddenPayloadFinding
	if err := scanPayloadValue(payload, "$", 0, &keys, &findings); err != nil {
		return nil, err
	}
	return findings, nil
}

func scanPayloadValue(
	value any,
	path string,
	depth int,
	keys *int,
	findings *[]HiddenPayloadFinding,
) error {
	if depth > maxPayloadDepth {
		return errors.New("ui payload exceeds nesting limit")
	}
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			*keys++
			if *keys > maxPayloadKeys {
				return errors.New("ui payload exceeds key limit")
			}
			lower := strings.ToLower(key)
			if strings.HasPrefix(lower, "on") ||
				strings.Contains(lower, "script") ||
				strings.Contains(lower, "html") ||
				strings.Contains(lower, "style") {
				*findings = append(*findings, HiddenPayloadFinding{
					Path:  path + "." + key,
					Class: "forbidden_key",
				})
			}
			if err := scanPayloadValue(
				child,
				path+"."+key,
				depth+1,
				keys,
				findings,
			); err != nil {
				return err
			}
		}
	case []any:
		for i, child := range typed {
			if err := scanPayloadValue(
				child,
				fmt.Sprintf("%s[%d]", path, i),
				depth+1,
				keys,
				findings,
			); err != nil {
				return err
			}
		}
	case string:
		scanString(typed, path, findings)
	case nil, bool, float64, int, int64:
	default:
		return fmt.Errorf("ui payload contains unsupported type %T", value)
	}
	return nil
}

func scanString(value, path string, findings *[]HiddenPayloadFinding) {
	if !utf8.ValidString(value) {
		*findings = append(*findings, HiddenPayloadFinding{Path: path, Class: "invalid_utf8"})
		return
	}
	for _, character := range value {
		switch character {
		case '\u200b', '\u200c', '\u200d', '\u2060',
			'\u202a', '\u202b', '\u202c', '\u202d', '\u202e',
			'\u2066', '\u2067', '\u2068', '\u2069':
			*findings = append(*findings, HiddenPayloadFinding{
				Path:  path,
				Class: "invisible_or_bidi_text",
			})
			return
		}
	}
	for _, candidate := range hiddenPatterns {
		if candidate.pattern.MatchString(value) {
			*findings = append(*findings, HiddenPayloadFinding{
				Path:  path,
				Class: candidate.class,
			})
		}
	}
	if len(value) >= 24 && len(value) <= 16<<10 {
		decoded, err := base64.StdEncoding.DecodeString(value)
		if err == nil && utf8.Valid(decoded) {
			for _, candidate := range hiddenPatterns {
				if candidate.pattern.Match(decoded) {
					*findings = append(*findings, HiddenPayloadFinding{
						Path:  path,
						Class: "encoded_" + candidate.class,
					})
				}
			}
		}
	}
}
