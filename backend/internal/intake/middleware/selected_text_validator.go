package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"

	"aeolyzer/internal/intake/contracts"
)

var (
	ErrEditSelectionRequired = errors.New("EDIT_SELECTION_REQUIRED")
	// Triggered if the client-submitted selection hash diverges from the server-computed one.
	// Typically occurs when the UI editor state drifts from the orchestration backend.
	ErrSelectedTextMismatch = errors.New("SELECTED_TEXT_HASH_MISMATCH")
	ErrUnsafeTextPayload    = errors.New("UNSAFE_TEXT_PAYLOAD")
)

// Sanitizes user selection buffers before they reach the orchestration DAG.
// Guards against control-character poisoning and ensures the buffer complies
// with the maximum context window limits.
func ValidateSelectedText(text string) (contracts.SanitizedSelectedText, error) {
	if text == "" {
		return contracts.SanitizedSelectedText{}, ErrEditSelectionRequired
	}

	// Hard limit of 12k bytes to avoid memory exhaustion during AST parsing
	// or subsequent prompt construction.
	if len(text) > 12000 {
		return contracts.SanitizedSelectedText{}, errors.New("selected text exceeds 12000 limit")
	}

	if err := RejectSelectedTextWithHiddenPayload(text); err != nil {
		return contracts.SanitizedSelectedText{}, err
	}

	// Strip unprintable control characters, but preserve structural whitespace (CR, LF, TAB).
	// This prevents invisible homoglyphs/zero-width chars from corrupting the diff matching.
	sanitized := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			if r != '\n' && r != '\r' && r != '\t' {
				return -1
			}
		}
		return r
	}, text)

	// Bind a cryptographic hash to the selection. Later stages use this hash
	// to enforce exact-match patches (preventing hallucinated diffs).
	return contracts.SanitizedSelectedText{
		Text: sanitized,
		Hash: HashSelectedText(sanitized),
	}, nil
}

func HashSelectedText(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

func ValidateSelectedTextHash(text string, expectedHash string) error {
	if HashSelectedText(text) != expectedHash {
		return ErrSelectedTextMismatch
	}
	return nil
}

// Simple heuristic to block the most common roleplay injection primitives.
// Note: Advanced semantic evasion is expected to be handled by `llm_firewall.go`.
func RejectSelectedTextWithHiddenPayload(text string) error {
	if strings.Contains(text, "SYSTEM:") || strings.Contains(text, "TOOL:") {
		return ErrUnsafeTextPayload
	}
	return nil
}
