package a2aserver

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"aeolyzer/internal/extensions"
	"aeolyzer/internal/extensions/security"
)

// Envelope is a signed, expiring A2A application message.
type Envelope struct {
	MessageID   string         `json:"message_id"`
	SenderID    string         `json:"sender_id"`
	RecipientID string         `json:"recipient_id"`
	TaskType    string         `json:"task_type"`
	Payload     map[string]any `json:"payload"`
	IssuedAt    time.Time      `json:"issued_at"`
	ExpiresAt   time.Time      `json:"expires_at"`
	Nonce       string         `json:"nonce"`
	Signature   string         `json:"signature"`
}

// PublicKeyResolver resolves the expected sender signing identity.
type PublicKeyResolver interface {
	PublicKey(context.Context, string) (ed25519.PublicKey, error)
}

// EnvelopeValidator verifies schema, time, sender signature, hidden payload,
// and nonce replay.
type EnvelopeValidator struct {
	mu        sync.Mutex
	schemas   *extensions.Schemas
	keys      PublicKeyResolver
	now       func() time.Time
	nonces    map[string]time.Time
	maxNonces int
}

// NewEnvelopeValidator creates a bounded validator.
func NewEnvelopeValidator(
	keys PublicKeyResolver,
	now func() time.Time,
	maxNonces int,
) (*EnvelopeValidator, error) {
	if keys == nil || now == nil || maxNonces < 100 || maxNonces > 1_000_000 {
		return nil, errors.New("a2a envelope validator is not configured")
	}
	schemas, err := extensions.NewSchemas()
	if err != nil {
		return nil, err
	}
	return &EnvelopeValidator{
		schemas:   schemas,
		keys:      keys,
		now:       now,
		nonces:    make(map[string]time.Time),
		maxNonces: maxNonces,
	}, nil
}

// SignEnvelope signs an envelope with Ed25519.
func SignEnvelope(
	envelope Envelope,
	privateKey ed25519.PrivateKey,
) (Envelope, error) {
	if len(privateKey) != ed25519.PrivateKeySize {
		return Envelope{}, errors.New("a2a private key is invalid")
	}
	envelope.Signature = ""
	payload, err := json.Marshal(envelope)
	if err != nil {
		return Envelope{}, fmt.Errorf("encode a2a envelope: %w", err)
	}
	envelope.Signature = hex.EncodeToString(ed25519.Sign(privateKey, payload))
	return envelope, nil
}

// Validate validates and consumes the envelope nonce.
func (v *EnvelopeValidator) Validate(
	ctx context.Context,
	envelope Envelope,
) error {
	if v == nil || v.schemas == nil || v.keys == nil || v.now == nil {
		return errors.New("a2a envelope validator is not configured")
	}
	data, err := json.Marshal(envelope)
	if err != nil {
		return errors.New("a2a envelope cannot be encoded")
	}
	if err := v.schemas.ValidateJSON(extensions.ContractA2AEnvelope, data); err != nil {
		return err
	}
	now := v.now()
	if envelope.IssuedAt.After(now.Add(time.Minute)) ||
		!now.Before(envelope.ExpiresAt) ||
		!envelope.ExpiresAt.After(envelope.IssuedAt) ||
		envelope.ExpiresAt.Sub(envelope.IssuedAt) > 5*time.Minute {
		return errors.New("a2a envelope validity window is invalid")
	}
	findings, err := security.ScanHiddenPayload(envelope.Payload)
	if err != nil {
		return err
	}
	if len(findings) > 0 {
		return errors.New("a2a envelope contains unsafe hidden payload")
	}
	publicKey, err := v.keys.PublicKey(ctx, envelope.SenderID)
	if err != nil {
		return fmt.Errorf("resolve a2a sender key: %w", err)
	}
	if len(publicKey) != ed25519.PublicKeySize {
		return errors.New("a2a sender public key is invalid")
	}
	signature, err := hex.DecodeString(envelope.Signature)
	if err != nil {
		return errors.New("a2a envelope signature is malformed")
	}
	unsigned := envelope
	unsigned.Signature = ""
	payload, err := json.Marshal(unsigned)
	if err != nil {
		return errors.New("a2a envelope cannot be encoded")
	}
	if !ed25519.Verify(publicKey, payload, signature) {
		return errors.New("a2a envelope signature verification failed")
	}
	return v.consumeNonce(envelope.SenderID+"\x00"+envelope.Nonce, envelope.ExpiresAt, now)
}

func (v *EnvelopeValidator) consumeNonce(
	nonce string,
	expiresAt, now time.Time,
) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	for existing, expiry := range v.nonces {
		if !now.Before(expiry) {
			delete(v.nonces, existing)
		}
	}
	if _, replay := v.nonces[nonce]; replay {
		return errors.New("a2a envelope nonce replay detected")
	}
	if len(v.nonces) >= v.maxNonces {
		return errors.New("a2a envelope replay cache is full")
	}
	v.nonces[nonce] = expiresAt
	return nil
}
