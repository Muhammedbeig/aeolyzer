package governanceaudit

import (
	"context"
	"strings"
	"testing"
)

func TestLedgerAppendsAndVerifiesSignedHashChain(t *testing.T) {
	signer, err := GenerateEd25519Signer("test-key")
	if err != nil {
		t.Fatalf("GenerateEd25519Signer() failed: %v", err)
	}
	ledger, err := NewLedger(signer)
	if err != nil {
		t.Fatalf("NewLedger() failed: %v", err)
	}
	for _, decisionType := range []string{"release_candidate", "release_approved"} {
		record, err := ledger.Append(context.Background(), validDecision(decisionType))
		if err != nil {
			t.Fatalf("Ledger.Append() failed: %v", err)
		}
		if record.Hash == "" || record.Signature == "" {
			t.Fatalf("Ledger.Append() returned unsigned record: %+v", record)
		}
	}
	records := ledger.Snapshot()
	if err := Verify(context.Background(), records, signer); err != nil {
		t.Fatalf("Verify() failed: %v", err)
	}
	if records[1].PreviousHash != records[0].Hash {
		t.Fatal("Ledger did not chain records")
	}
}

func TestVerifyRejectsTampering(t *testing.T) {
	signer, err := GenerateEd25519Signer("test-key")
	if err != nil {
		t.Fatalf("GenerateEd25519Signer() failed: %v", err)
	}
	ledger, err := NewLedger(signer)
	if err != nil {
		t.Fatalf("NewLedger() failed: %v", err)
	}
	if _, err := ledger.Append(context.Background(), validDecision("release_rejected")); err != nil {
		t.Fatalf("Ledger.Append() failed: %v", err)
	}
	records := ledger.Snapshot()
	records[0].Decision = "release_approved"
	if err := Verify(context.Background(), records, signer); err == nil {
		t.Fatal("Verify() accepted a tampered record")
	}
}

func TestLedgerSnapshotIsDefensiveCopy(t *testing.T) {
	signer, err := GenerateEd25519Signer("test-key")
	if err != nil {
		t.Fatalf("GenerateEd25519Signer() failed: %v", err)
	}
	ledger, err := NewLedger(signer)
	if err != nil {
		t.Fatalf("NewLedger() failed: %v", err)
	}
	if _, err := ledger.Append(context.Background(), validDecision("release_rejected")); err != nil {
		t.Fatalf("Ledger.Append() failed: %v", err)
	}
	records := ledger.Snapshot()
	records[0].EvidenceRefs[0] = strings.Repeat("f", 64)
	again := ledger.Snapshot()
	if again[0].EvidenceRefs[0] == records[0].EvidenceRefs[0] {
		t.Fatal("Snapshot() exposed mutable ledger state")
	}
}

func TestLedgerRejectsRawOrUnhashedEvidence(t *testing.T) {
	signer, err := GenerateEd25519Signer("test-key")
	if err != nil {
		t.Fatalf("GenerateEd25519Signer() failed: %v", err)
	}
	ledger, err := NewLedger(signer)
	if err != nil {
		t.Fatalf("NewLedger() failed: %v", err)
	}
	decision := validDecision("release_rejected")
	decision.EvidenceRefs = []string{"raw prompt text"}
	if _, err := ledger.Append(context.Background(), decision); err == nil {
		t.Fatal("Ledger.Append() accepted raw evidence")
	}
}

func validDecision(decisionType string) Decision {
	return Decision{
		TenantID:      "tenant-1",
		DecisionType:  decisionType,
		Decision:      "blocked",
		SafeSummary:   "Release remains blocked by required safety evidence.",
		ActorHash:     strings.Repeat("a", 64),
		PolicyVersion: "governance-v2",
		EvidenceRefs:  []string{strings.Repeat("b", 64)},
		HumanReviewed: true,
	}
}
