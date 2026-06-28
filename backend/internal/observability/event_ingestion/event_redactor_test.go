package eventingestion

import (
	"encoding/json"
	"strings"
	"testing"

	observabilityconfig "aeolyzer/internal/observability/config"
)

func TestRedactorDropsSecretsAndProtectedMetadata(t *testing.T) {
	redactor := newTestRedactor(t)
	input := map[string]any{
		"event_type":         "runtime.completed",
		"raw_api_key":        "do-not-store",
		"workflow_id":        "internal-workflow",
		"authorization":      "Bearer fake-token-that-is-long-enough",
		"selected_text_hash": "stable-selection",
		"reasoning_summary":  "Safe summary.",
		"nested": map[string]any{
			"password": "secret",
			"status":   "ok",
		},
	}

	output, err := redactor.RedactMap(input)
	if err != nil {
		t.Fatalf("RedactMap() failed: %v", err)
	}
	for _, field := range []string{"raw_api_key", "workflow_id", "authorization"} {
		if _, found := output[field]; found {
			t.Errorf("RedactMap() retained forbidden field %q", field)
		}
	}
	nested := output["nested"].(map[string]any)
	if _, found := nested["password"]; found {
		t.Error("RedactMap() retained nested password")
	}
	if nested["status"] != "ok" {
		t.Errorf("RedactMap() nested status = %v, want ok", nested["status"])
	}
	hash, ok := output["selected_text_hash"].(string)
	if !ok || !strings.HasPrefix(hash, "hmac-sha256:") || strings.Contains(hash, "stable-selection") {
		t.Errorf("RedactMap() hash = %v, want non-reversible HMAC", output["selected_text_hash"])
	}
}

func TestRedactorDoesNotMutateInput(t *testing.T) {
	redactor := newTestRedactor(t)
	input := map[string]any{"raw_api_key": "secret", "status": "ok"}
	if _, err := redactor.RedactMap(input); err != nil {
		t.Fatalf("RedactMap() failed: %v", err)
	}
	if input["raw_api_key"] != "secret" {
		t.Fatal("RedactMap() mutated caller input")
	}
}

func TestRedactJSONRejectsTrailingAndUnsupportedData(t *testing.T) {
	redactor := newTestRedactor(t)
	if _, err := redactor.RedactJSON([]byte(`{"status":"ok"} {}`)); err == nil {
		t.Fatal("RedactJSON() accepted trailing data")
	}
	if _, err := redactor.RedactJSON(nil); err == nil {
		t.Fatal("RedactJSON() accepted empty payload")
	}
}

func TestRedactJSONReturnsValidRedactedObject(t *testing.T) {
	redactor := newTestRedactor(t)
	output, err := redactor.RedactJSON([]byte(`{
		"status":"ok",
		"raw_user_prompt":"never store",
		"message":"Authorization: Bearer abcdefghijklmnopqrstuvwxyz"
	}`))
	if err != nil {
		t.Fatalf("RedactJSON() failed: %v", err)
	}
	var object map[string]any
	if err := json.Unmarshal(output, &object); err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}
	if _, found := object["raw_user_prompt"]; found {
		t.Fatal("RedactJSON() retained raw user prompt")
	}
	if object["message"] != "[REDACTED_SECRET]" {
		t.Errorf("RedactJSON() message = %v, want redacted secret", object["message"])
	}
}

func TestLegacyEventRedaction(t *testing.T) {
	output := RedactEvent(
		"raw_system_prompt Authorization: Bearer abcdefghijklmnopqrstuvwxyz",
	)
	if strings.Contains(output, "raw_system_prompt") ||
		strings.Contains(output, "abcdefghijklmnopqrstuvwxyz") {
		t.Fatalf("RedactEvent() leaked protected input: %s", output)
	}
}

func TestNewRedactorRequiresStrongHMACKey(t *testing.T) {
	policies, err := observabilityconfig.LoadEmbeddedPolicies()
	if err != nil {
		t.Fatalf("LoadEmbeddedPolicies() failed: %v", err)
	}
	if _, err := NewRedactor(policies.Redaction, []byte("short")); err == nil {
		t.Fatal("NewRedactor() accepted a weak HMAC key")
	}
}

func newTestRedactor(t *testing.T) *Redactor {
	t.Helper()
	policies, err := observabilityconfig.LoadEmbeddedPolicies()
	if err != nil {
		t.Fatalf("LoadEmbeddedPolicies() failed: %v", err)
	}
	redactor, err := NewRedactor(
		policies.Redaction,
		[]byte("01234567890123456789012345678901"),
	)
	if err != nil {
		t.Fatalf("NewRedactor() failed: %v", err)
	}
	return redactor
}
