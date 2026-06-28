// Package governanceaudit records signed, append-only governance decisions.
package governanceaudit

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"
)

var hashPattern = regexp.MustCompile(`^(?:sha256:)?[a-f0-9]{64}$`)

// Signer signs a record hash using a key held outside the ledger.
type Signer interface {
	KeyID() string
	Sign(context.Context, []byte) ([]byte, error)
}

// Verifier resolves a public key for one signing key ID.
type Verifier interface {
	PublicKey(context.Context, string) (ed25519.PublicKey, error)
}

// Decision is sanitized input for one governance record.
type Decision struct {
	TenantID      string   `json:"tenant_id"`
	DecisionType  string   `json:"decision_type"`
	Decision      string   `json:"decision"`
	SafeSummary   string   `json:"safe_summary"`
	ActorHash     string   `json:"actor_hash"`
	PolicyVersion string   `json:"policy_version"`
	EvidenceRefs  []string `json:"evidence_refs"`
	HumanReviewed bool     `json:"human_reviewed"`
}

// Record is one immutable hash-chained and signed decision.
type Record struct {
	Sequence      uint64    `json:"sequence"`
	Timestamp     time.Time `json:"timestamp"`
	TenantID      string    `json:"tenant_id"`
	DecisionType  string    `json:"decision_type"`
	Decision      string    `json:"decision"`
	SafeSummary   string    `json:"safe_summary"`
	ActorHash     string    `json:"actor_hash"`
	PolicyVersion string    `json:"policy_version"`
	EvidenceRefs  []string  `json:"evidence_refs"`
	HumanReviewed bool      `json:"human_reviewed"`
	PreviousHash  string    `json:"previous_hash"`
	Hash          string    `json:"hash"`
	SigningKeyID  string    `json:"signing_key_id"`
	Signature     string    `json:"signature"`
}

// Ledger is an in-memory append-only ledger. A production deployment must
// persist records to an immutable store after Append succeeds.
type Ledger struct {
	mu      sync.RWMutex
	signer  Signer
	now     func() time.Time
	records []Record
}

// NewLedger creates a ledger with a required external signer.
func NewLedger(signer Signer) (*Ledger, error) {
	if signer == nil || signer.KeyID() == "" {
		return nil, errors.New("governance signer and key id are required")
	}
	return &Ledger{
		signer: signer,
		now: func() time.Time {
			return time.Now().UTC()
		},
	}, nil
}

// Append validates, hashes, signs, and atomically appends one record.
func (l *Ledger) Append(ctx context.Context, decision Decision) (Record, error) {
	if l == nil || l.signer == nil {
		return Record{}, errors.New("governance ledger is not configured")
	}
	if err := validateDecision(decision); err != nil {
		return Record{}, err
	}
	if err := ctx.Err(); err != nil {
		return Record{}, fmt.Errorf("append governance record: %w", err)
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	record := Record{
		Sequence:      uint64(len(l.records) + 1),
		Timestamp:     l.now(),
		TenantID:      decision.TenantID,
		DecisionType:  decision.DecisionType,
		Decision:      decision.Decision,
		SafeSummary:   decision.SafeSummary,
		ActorHash:     decision.ActorHash,
		PolicyVersion: decision.PolicyVersion,
		EvidenceRefs:  append([]string(nil), decision.EvidenceRefs...),
		HumanReviewed: decision.HumanReviewed,
		SigningKeyID:  l.signer.KeyID(),
	}
	if len(l.records) > 0 {
		record.PreviousHash = l.records[len(l.records)-1].Hash
	}
	hash, err := hashRecord(record)
	if err != nil {
		return Record{}, err
	}
	record.Hash = hash
	signature, err := l.signer.Sign(ctx, []byte(hash))
	if err != nil {
		return Record{}, fmt.Errorf("sign governance record: %w", err)
	}
	if len(signature) == 0 {
		return Record{}, errors.New("governance signer returned an empty signature")
	}
	record.Signature = hex.EncodeToString(signature)
	l.records = append(l.records, record)
	return cloneRecord(record), nil
}

// Snapshot returns a defensive copy of all records.
func (l *Ledger) Snapshot() []Record {
	if l == nil {
		return nil
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	result := make([]Record, len(l.records))
	for i, record := range l.records {
		result[i] = cloneRecord(record)
	}
	return result
}

// Verify validates ordering, timestamps, hash chaining, record hashes, and
// signatures for a snapshot.
func Verify(ctx context.Context, records []Record, verifier Verifier) error {
	if verifier == nil {
		return errors.New("governance verifier is required")
	}
	var previousHash string
	var previousTime time.Time
	for i, record := range records {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("verify governance ledger: %w", err)
		}
		if record.Sequence != uint64(i+1) {
			return fmt.Errorf("record %d has invalid sequence", i+1)
		}
		if record.Timestamp.IsZero() || (!previousTime.IsZero() && record.Timestamp.Before(previousTime)) {
			return fmt.Errorf("record %d has invalid timestamp ordering", i+1)
		}
		if record.PreviousHash != previousHash {
			return fmt.Errorf("record %d has invalid previous hash", i+1)
		}
		expected, err := hashRecord(record)
		if err != nil {
			return err
		}
		if record.Hash != expected {
			return fmt.Errorf("record %d hash mismatch", i+1)
		}
		signature, err := hex.DecodeString(record.Signature)
		if err != nil {
			return fmt.Errorf("record %d signature encoding: %w", i+1, err)
		}
		publicKey, err := verifier.PublicKey(ctx, record.SigningKeyID)
		if err != nil {
			return fmt.Errorf("resolve record %d public key: %w", i+1, err)
		}
		if !ed25519.Verify(publicKey, []byte(record.Hash), signature) {
			return fmt.Errorf("record %d signature verification failed", i+1)
		}
		previousHash = record.Hash
		previousTime = record.Timestamp
	}
	return nil
}

// Ed25519Signer is a local signer suitable for tests and isolated deployments.
// Production should use a KMS- or HSM-backed Signer implementation.
type Ed25519Signer struct {
	keyID      string
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewEd25519Signer creates a signer from an existing private key.
func NewEd25519Signer(keyID string, privateKey ed25519.PrivateKey) (*Ed25519Signer, error) {
	if keyID == "" {
		return nil, errors.New("signing key id is required")
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid ed25519 private key")
	}
	publicKey := privateKey.Public().(ed25519.PublicKey)
	return &Ed25519Signer{
		keyID:      keyID,
		privateKey: append(ed25519.PrivateKey(nil), privateKey...),
		publicKey:  append(ed25519.PublicKey(nil), publicKey...),
	}, nil
}

// GenerateEd25519Signer creates an ephemeral signer. It must not be used as a
// substitute for production key management.
func GenerateEd25519Signer(keyID string) (*Ed25519Signer, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate ed25519 key: %w", err)
	}
	return NewEd25519Signer(keyID, privateKey)
}

// KeyID returns the non-secret signing-key identifier.
func (s *Ed25519Signer) KeyID() string {
	if s == nil {
		return ""
	}
	return s.keyID
}

// Sign signs one record hash.
func (s *Ed25519Signer) Sign(ctx context.Context, digest []byte) ([]byte, error) {
	if s == nil || len(s.privateKey) != ed25519.PrivateKeySize {
		return nil, errors.New("ed25519 signer is not configured")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return ed25519.Sign(s.privateKey, digest), nil
}

// PublicKey resolves this signer's public key.
func (s *Ed25519Signer) PublicKey(
	ctx context.Context,
	keyID string,
) (ed25519.PublicKey, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if s == nil || keyID != s.keyID {
		return nil, errors.New("signing key is unknown")
	}
	return append(ed25519.PublicKey(nil), s.publicKey...), nil
}

func validateDecision(decision Decision) error {
	if decision.TenantID == "" ||
		decision.DecisionType == "" ||
		decision.Decision == "" ||
		decision.SafeSummary == "" ||
		decision.ActorHash == "" ||
		decision.PolicyVersion == "" {
		return errors.New("governance decision fields are required")
	}
	if len(decision.SafeSummary) > 2000 {
		return errors.New("governance safe summary exceeds limit")
	}
	if !hashPattern.MatchString(decision.ActorHash) {
		return errors.New("governance actor hash is invalid")
	}
	if len(decision.EvidenceRefs) == 0 || len(decision.EvidenceRefs) > 64 {
		return errors.New("governance evidence references are required and bounded")
	}
	for _, evidence := range decision.EvidenceRefs {
		if !hashPattern.MatchString(evidence) {
			return errors.New("governance evidence reference is invalid")
		}
	}
	return nil
}

func hashRecord(record Record) (string, error) {
	hashable := struct {
		Sequence      uint64    `json:"sequence"`
		Timestamp     time.Time `json:"timestamp"`
		TenantID      string    `json:"tenant_id"`
		DecisionType  string    `json:"decision_type"`
		Decision      string    `json:"decision"`
		SafeSummary   string    `json:"safe_summary"`
		ActorHash     string    `json:"actor_hash"`
		PolicyVersion string    `json:"policy_version"`
		EvidenceRefs  []string  `json:"evidence_refs"`
		HumanReviewed bool      `json:"human_reviewed"`
		PreviousHash  string    `json:"previous_hash"`
		SigningKeyID  string    `json:"signing_key_id"`
	}{
		Sequence:      record.Sequence,
		Timestamp:     record.Timestamp.UTC(),
		TenantID:      record.TenantID,
		DecisionType:  record.DecisionType,
		Decision:      record.Decision,
		SafeSummary:   record.SafeSummary,
		ActorHash:     record.ActorHash,
		PolicyVersion: record.PolicyVersion,
		EvidenceRefs:  record.EvidenceRefs,
		HumanReviewed: record.HumanReviewed,
		PreviousHash:  record.PreviousHash,
		SigningKeyID:  record.SigningKeyID,
	}
	data, err := json.Marshal(hashable)
	if err != nil {
		return "", fmt.Errorf("encode governance record: %w", err)
	}
	digest := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

func cloneRecord(record Record) Record {
	result := record
	result.EvidenceRefs = append([]string(nil), record.EvidenceRefs...)
	return result
}
