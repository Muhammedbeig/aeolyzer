// Package eventingestion validates and redacts telemetry before storage.
package eventingestion

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	observabilityconfig "aeolyzer/internal/observability/config"
)

const (
	maxRedactionDepth = 32
	maxRedactionKeys  = 10_000
)

var secretValuePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\bBearer\s+[A-Za-z0-9._~+/-]{16,}=*\b`),
	regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{20,}\b`),
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}\b`),
	regexp.MustCompile(`-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----`),
}

// Redactor applies the validated Layer 8 redaction policy. The HMAC key must be
// tenant-scoped or deployment-scoped and supplied by a secret manager.
type Redactor struct {
	neverStore     map[string]struct{}
	storeAsHash    map[string]struct{}
	storeAsSummary map[string]struct{}
	protected      map[string]struct{}
	hmacKey        []byte
}

// NewRedactor builds a redactor from validated policy and a non-exportable HMAC
// key. The key is copied and never included in errors or outputs.
func NewRedactor(policy observabilityconfig.RedactionPolicy, hmacKey []byte) (*Redactor, error) {
	if len(hmacKey) < 32 {
		return nil, errors.New("redaction hmac key must contain at least 32 bytes")
	}
	if len(policy.NeverStore) == 0 {
		return nil, errors.New("redaction never-store policy is empty")
	}
	return &Redactor{
		neverStore:     normalizedSet(policy.NeverStore),
		storeAsHash:    normalizedSet(policy.StoreAsHash),
		storeAsSummary: normalizedSet(policy.StoreAsRedactedSummary),
		protected:      normalizedSet(policy.ProtectedInternalMetadata),
		hmacKey:        append([]byte(nil), hmacKey...),
	}, nil
}

// RedactMap returns a deep redacted copy. It never mutates caller-owned data.
func (r *Redactor) RedactMap(payload map[string]any) (map[string]any, error) {
	if r == nil {
		return nil, errors.New("redactor is not configured")
	}
	if payload == nil {
		return nil, errors.New("redaction payload is required")
	}
	keyCount := 0
	value, err := r.redactValue(payload, 0, &keyCount)
	if err != nil {
		return nil, err
	}
	result, ok := value.(map[string]any)
	if !ok {
		return nil, errors.New("redaction result is not an object")
	}
	return result, nil
}

// RedactJSON strictly decodes, redacts, and re-encodes one JSON object.
func (r *Redactor) RedactJSON(payload []byte) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errors.New("redaction payload is empty")
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	var object map[string]any
	if err := decoder.Decode(&object); err != nil {
		return nil, fmt.Errorf("decode redaction payload: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("redaction payload contains trailing data")
	}
	redacted, err := r.RedactMap(object)
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(redacted)
	if err != nil {
		return nil, fmt.Errorf("encode redacted payload: %w", err)
	}
	return result, nil
}

// RedactEvent provides conservative redaction for legacy unstructured events.
// New event ingestion should use RedactMap or RedactJSON.
func RedactEvent(payload string) string {
	for _, field := range []string{
		"hidden_chain_of_thought",
		"raw_system_prompt",
		"raw_developer_prompt",
		"raw_user_prompt",
		"raw_api_key",
		"raw_oauth_token",
		"password",
		"private_key",
		"cookie",
	} {
		payload = strings.ReplaceAll(payload, field, "[REDACTED_FIELD]")
	}
	for _, pattern := range secretValuePatterns {
		payload = pattern.ReplaceAllString(payload, "[REDACTED_SECRET]")
	}
	return payload
}

func (r *Redactor) redactValue(value any, depth int, keyCount *int) (any, error) {
	if depth > maxRedactionDepth {
		return nil, errors.New("redaction payload exceeds nesting limit")
	}
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, child := range typed {
			*keyCount++
			if *keyCount > maxRedactionKeys {
				return nil, errors.New("redaction payload exceeds key limit")
			}
			normalized := normalizeField(key)
			if r.mustDrop(normalized) {
				continue
			}
			if _, hash := r.storeAsHash[normalized]; hash {
				result[key] = r.hashValue(child)
				continue
			}
			if _, summary := r.storeAsSummary[normalized]; summary {
				text, ok := child.(string)
				if !ok || containsSecret(text) {
					result[key] = "[REDACTED_SUMMARY]"
					continue
				}
				result[key] = truncate(text, 2000)
				continue
			}
			redacted, err := r.redactValue(child, depth+1, keyCount)
			if err != nil {
				return nil, err
			}
			result[key] = redacted
		}
		return result, nil
	case []any:
		result := make([]any, len(typed))
		for i, child := range typed {
			redacted, err := r.redactValue(child, depth+1, keyCount)
			if err != nil {
				return nil, err
			}
			result[i] = redacted
		}
		return result, nil
	case string:
		if containsSecret(typed) {
			return "[REDACTED_SECRET]", nil
		}
		return typed, nil
	case nil, bool, json.Number, float64:
		return typed, nil
	default:
		return nil, fmt.Errorf("redaction payload contains unsupported type %T", value)
	}
}

func (r *Redactor) mustDrop(field string) bool {
	if _, found := r.neverStore[field]; found {
		return true
	}
	if _, found := r.protected[field]; found {
		return true
	}
	return secretLikeField(field)
}

func (r *Redactor) hashValue(value any) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "[HASH_UNAVAILABLE]"
	}
	mac := hmac.New(sha256.New, r.hmacKey)
	_, _ = mac.Write(encoded)
	return "hmac-sha256:" + hex.EncodeToString(mac.Sum(nil))
}

func normalizedSet(values []string) map[string]struct{} {
	result := make(map[string]struct{}, len(values))
	for _, value := range values {
		result[normalizeField(value)] = struct{}{}
	}
	return result
}

func normalizeField(field string) string {
	return strings.ToLower(strings.TrimSpace(field))
}

func secretLikeField(field string) bool {
	for _, fragment := range []string{
		"api_key",
		"apikey",
		"access_token",
		"refresh_token",
		"oauth_token",
		"authorization",
		"password",
		"passwd",
		"private_key",
		"client_secret",
		"cookie",
	} {
		if strings.Contains(field, fragment) {
			return true
		}
	}
	return false
}

func containsSecret(value string) bool {
	for _, pattern := range secretValuePatterns {
		if pattern.MatchString(value) {
			return true
		}
	}
	return false
}

func truncate(value string, limit int) string {
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}
