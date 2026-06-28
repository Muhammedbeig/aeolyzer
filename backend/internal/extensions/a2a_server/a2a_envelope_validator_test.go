package a2aserver

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"testing"
	"time"
)

type keyResolver struct {
	sender string
	key    ed25519.PublicKey
}

func (r keyResolver) PublicKey(
	_ context.Context,
	sender string,
) (ed25519.PublicKey, error) {
	if sender != r.sender {
		return nil, errors.New("sender is unknown")
	}
	return r.key, nil
}

func TestEnvelopeValidatorRejectsTamperingReplayAndHiddenPayload(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("ed25519.GenerateKey() failed: %v", err)
	}
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	validator, err := NewEnvelopeValidator(
		keyResolver{sender: "audit-agent", key: publicKey},
		func() time.Time { return now },
		100,
	)
	if err != nil {
		t.Fatalf("NewEnvelopeValidator() failed: %v", err)
	}
	envelope := Envelope{
		MessageID:   "message-1",
		SenderID:    "audit-agent",
		RecipientID: "content-agent",
		TaskType:    "audit_summary",
		Payload:     map[string]any{"summary": "Safe audit result."},
		IssuedAt:    now,
		ExpiresAt:   now.Add(time.Minute),
		Nonce:       "0123456789abcdef",
	}
	signed, err := SignEnvelope(envelope, privateKey)
	if err != nil {
		t.Fatalf("SignEnvelope() failed: %v", err)
	}
	if err := validator.Validate(context.Background(), signed); err != nil {
		t.Fatalf("EnvelopeValidator.Validate() failed: %v", err)
	}
	if err := validator.Validate(context.Background(), signed); err == nil {
		t.Fatal("EnvelopeValidator.Validate() accepted replay")
	}

	tampered := signed
	tampered.Nonce = "fedcba9876543210"
	tampered.Payload = map[string]any{"summary": "Tampered"}
	if err := validator.Validate(context.Background(), tampered); err == nil {
		t.Fatal("EnvelopeValidator.Validate() accepted tampering")
	}

	unsafe := envelope
	unsafe.Nonce = "abcdef0123456789"
	unsafe.Payload = map[string]any{
		"summary": "<script>alert(1)</script>",
	}
	unsafe, err = SignEnvelope(unsafe, privateKey)
	if err != nil {
		t.Fatalf("SignEnvelope(unsafe) failed: %v", err)
	}
	if err := validator.Validate(context.Background(), unsafe); err == nil {
		t.Fatal("EnvelopeValidator.Validate() accepted hidden payload")
	}
}
