package datasecuritymesh

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ProvenanceSigner signs evidence packet attestations.
type ProvenanceSigner interface {
	KeyID() string
	Sign(context.Context, []byte) ([]byte, error)
}

// ProvenanceInput identifies one tenant-bound source payload.
type ProvenanceInput struct {
	TenantID      string
	SourceID      string
	SourceVersion string
	RetrievedAt   time.Time
	ProjectedData map[string]any
}

// ProvenanceAttestation is safe evidence metadata.
type ProvenanceAttestation struct {
	TenantID      string    `json:"tenant_id"`
	SourceID      string    `json:"source_id"`
	SourceVersion string    `json:"source_version"`
	RetrievedAt   time.Time `json:"retrieved_at"`
	PayloadHash   string    `json:"payload_hash"`
	SigningKeyID  string    `json:"signing_key_id"`
	Signature     string    `json:"signature"`
}

// AttestProvenance hashes canonical JSON and signs its tenant/source binding.
func AttestProvenance(
	ctx context.Context,
	input ProvenanceInput,
	signer ProvenanceSigner,
) (ProvenanceAttestation, error) {
	if signer == nil ||
		signer.KeyID() == "" ||
		input.TenantID == "" ||
		input.SourceID == "" ||
		input.SourceVersion == "" ||
		input.RetrievedAt.IsZero() ||
		input.ProjectedData == nil {
		return ProvenanceAttestation{}, errors.New("provenance input is incomplete")
	}
	payload, err := json.Marshal(input.ProjectedData)
	if err != nil {
		return ProvenanceAttestation{}, fmt.Errorf("encode provenance payload: %w", err)
	}
	digest := sha256.Sum256(payload)
	attestation := ProvenanceAttestation{
		TenantID:      input.TenantID,
		SourceID:      input.SourceID,
		SourceVersion: input.SourceVersion,
		RetrievedAt:   input.RetrievedAt.UTC(),
		PayloadHash:   "sha256:" + hex.EncodeToString(digest[:]),
		SigningKeyID:  signer.KeyID(),
	}
	signingPayload, err := provenanceSigningPayload(attestation)
	if err != nil {
		return ProvenanceAttestation{}, err
	}
	signature, err := signer.Sign(ctx, signingPayload)
	if err != nil {
		return ProvenanceAttestation{}, fmt.Errorf("sign provenance: %w", err)
	}
	attestation.Signature = hex.EncodeToString(signature)
	return attestation, nil
}

// VerifyProvenance verifies the attestation against projected data.
func VerifyProvenance(
	attestation ProvenanceAttestation,
	projectedData map[string]any,
	publicKey ed25519.PublicKey,
) error {
	if len(publicKey) != ed25519.PublicKeySize ||
		attestation.Signature == "" ||
		projectedData == nil {
		return errors.New("provenance verification context is incomplete")
	}
	payload, err := json.Marshal(projectedData)
	if err != nil {
		return fmt.Errorf("encode provenance payload: %w", err)
	}
	digest := sha256.Sum256(payload)
	expectedHash := "sha256:" + hex.EncodeToString(digest[:])
	if attestation.PayloadHash != expectedHash {
		return errors.New("provenance payload hash mismatch")
	}
	signature, err := hex.DecodeString(attestation.Signature)
	if err != nil {
		return errors.New("provenance signature is malformed")
	}
	signingPayload, err := provenanceSigningPayload(attestation)
	if err != nil {
		return err
	}
	if !ed25519.Verify(publicKey, signingPayload, signature) {
		return errors.New("provenance signature verification failed")
	}
	return nil
}

func provenanceSigningPayload(attestation ProvenanceAttestation) ([]byte, error) {
	unsigned := attestation
	unsigned.Signature = ""
	payload, err := json.Marshal(unsigned)
	if err != nil {
		return nil, fmt.Errorf("encode provenance attestation: %w", err)
	}
	return payload, nil
}
