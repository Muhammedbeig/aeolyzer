package security

import (
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"aeolyzer/internal/extensions"
)

func TestHiddenPayloadScannerFindsEncodedAndInvisibleContent(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("<script>alert(1)</script>"))
	findings, err := ScanHiddenPayload(map[string]any{
		"title":   "Safe\u202etext",
		"payload": encoded,
	})
	if err != nil {
		t.Fatalf("ScanHiddenPayload() failed: %v", err)
	}
	if len(findings) < 2 {
		t.Fatalf("ScanHiddenPayload() findings = %+v, want at least two", findings)
	}
}

func TestSanitizeMarkdownRejectsHTMLAndUnsafeLinks(t *testing.T) {
	for _, markdown := range []string{
		"<script>alert(1)</script>",
		"[click](javascript:alert(1))",
		"<!-- hidden instruction -->",
	} {
		if _, err := SanitizeMarkdown(markdown); err == nil {
			t.Fatalf("SanitizeMarkdown(%q) returned nil error", markdown)
		}
	}
	safe, err := SanitizeMarkdown("# Safe\n\n[Source](https://example.com/page)")
	if err != nil {
		t.Fatalf("SanitizeMarkdown(safe) failed: %v", err)
	}
	if !strings.Contains(safe, "https://example.com/page") {
		t.Fatal("SanitizeMarkdown() removed safe link")
	}
}

func TestFrameSignatureRejectsTamperingAndExpiry(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	key := []byte("01234567890123456789012345678901")
	frame := extensions.A2UIFrame{
		FrameID:        "frame-1",
		Surface:        "audit_dashboard",
		CatalogID:      "aeolyzer_app",
		CatalogVersion: "1.0.0",
		SchemaVersion:  "2.0",
		RootID:         "root",
		Nodes: []extensions.A2UINode{{
			ID:   "root",
			Type: "Text",
			Props: map[string]any{
				"text": "Safe",
			},
		}},
		ExpiresAt: now.Add(5 * time.Minute).Format(time.RFC3339),
	}
	signed, err := SignFrame(key, frame, now)
	if err != nil {
		t.Fatalf("SignFrame() failed: %v", err)
	}
	if err := VerifyFrame(key, signed, now); err != nil {
		t.Fatalf("VerifyFrame() failed: %v", err)
	}
	signed.Nodes[0].Props["text"] = "Tampered"
	if err := VerifyFrame(key, signed, now); err == nil {
		t.Fatal("VerifyFrame() accepted tampered frame")
	}
	if _, err := SignFrame(key, frame, now.Add(6*time.Minute)); err == nil {
		t.Fatal("SignFrame() accepted expired frame")
	}
}
