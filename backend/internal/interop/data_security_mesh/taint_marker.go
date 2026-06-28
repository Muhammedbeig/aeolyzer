package datasecuritymesh

import (
	"errors"
	"regexp"
	"strings"
)

var taintPatterns = []struct {
	class   string
	pattern *regexp.Regexp
	weight  float64
}{
	{
		class:   "instruction_override",
		pattern: regexp.MustCompile(`(?i)\b(ignore|disregard|override)\b.{0,80}\b(previous|prior|system|developer)\b.{0,40}\b(instruction|message|prompt)s?\b`),
		weight:  0.45,
	},
	{
		class:   "role_impersonation",
		pattern: regexp.MustCompile(`(?i)\b(system|developer)\s*(message|prompt)\s*:`),
		weight:  0.35,
	},
	{
		class:   "tool_lure",
		pattern: regexp.MustCompile(`(?i)\b(call|invoke|execute|run)\b.{0,40}\b(tool|function|command|shell)\b`),
		weight:  0.30,
	},
	{
		class:   "data_exfiltration",
		pattern: regexp.MustCompile(`(?i)\b(send|upload|post|exfiltrate)\b.{0,80}\b(secret|token|credential|cookie|key|private data)\b`),
		weight:  0.55,
	},
	{
		class:   "memory_poisoning",
		pattern: regexp.MustCompile(`(?i)\b(remember|store in memory|persist)\b.{0,100}\b(instruction|rule|secret|credential)\b`),
		weight:  0.40,
	},
}

// TaintResult is metadata only; it never contains the inspected source text.
type TaintResult struct {
	Tainted bool     `json:"tainted"`
	Score   float64  `json:"score"`
	Classes []string `json:"classes,omitempty"`
}

// DetectTaint detects prompt-like instructions in retrieved evidence.
func DetectTaint(text string) (TaintResult, error) {
	if text == "" {
		return TaintResult{}, errors.New("taint input is empty")
	}
	if len(text) > 2<<20 {
		return TaintResult{}, errors.New("taint input exceeds size limit")
	}
	normalized := strings.ToValidUTF8(text, "")
	score := 0.0
	var classes []string
	for _, candidate := range taintPatterns {
		if candidate.pattern.MatchString(normalized) {
			score += candidate.weight
			classes = append(classes, candidate.class)
		}
	}
	if score > 1 {
		score = 1
	}
	return TaintResult{
		Tainted: score >= 0.30,
		Score:   score,
		Classes: classes,
	}, nil
}
