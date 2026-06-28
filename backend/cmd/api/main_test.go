package main

import (
	"encoding/hex"
	"testing"
)

func TestSigningKeyAndTraceIDUseExpectedEntropyLengths(t *testing.T) {
	key := newSigningKey()
	if len(key) != 32 {
		t.Fatalf("newSigningKey() length = %d, want 32", len(key))
	}
	traceID := newTraceID()
	decoded, err := hex.DecodeString(traceID)
	if err != nil {
		t.Fatalf("newTraceID() returned invalid hex: %v", err)
	}
	if len(decoded) != 16 {
		t.Fatalf("newTraceID() decoded length = %d, want 16", len(decoded))
	}
}

func TestStartupContractsValidate(t *testing.T) {
	if err := validateStartupContracts(); err != nil {
		t.Fatalf("validateStartupContracts() error = %v", err)
	}
}
