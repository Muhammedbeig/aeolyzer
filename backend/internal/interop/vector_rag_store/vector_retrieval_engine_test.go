package vectorragstore

import (
	"context"
	"errors"
	"math"
	"strings"
	"testing"
)

type backendStub struct {
	records []Record
}

func (b backendStub) Search(
	_ context.Context,
	_ string,
	_ []float32,
	_ int,
) ([]Record, error) {
	return b.records, nil
}

func TestEngineRejectsCrossTenantBackendResult(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	namespace, err := NamespaceForTenant(key, "tenant-a")
	if err != nil {
		t.Fatalf("NamespaceForTenant() failed: %v", err)
	}
	engine, err := NewEngine(backendStub{records: []Record{{
		ID:            "record-1",
		Namespace:     "tenant_attacker",
		ContentHash:   "sha256:" + strings.Repeat("a", 64),
		ProvenanceRef: "sha256:" + strings.Repeat("b", 64),
		Fields:        map[string]any{"title": "private"},
	}}}, key, 1536, 10)
	if err != nil {
		t.Fatalf("NewEngine() failed: %v", err)
	}
	if _, err := engine.Retrieve(context.Background(), Query{
		TenantID:           "tenant-a",
		CredentialTenantID: "tenant-a",
		Namespace:          namespace,
		Embedding:          []float32{0.1, 0.2},
		TopK:               5,
	}); err == nil {
		t.Fatal("Engine.Retrieve() accepted cross-namespace result")
	}
}

func TestEngineRejectsCredentialAndEmbeddingConfusion(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	namespace, err := NamespaceForTenant(key, "tenant-a")
	if err != nil {
		t.Fatalf("NamespaceForTenant() failed: %v", err)
	}
	engine, err := NewEngine(backendStub{}, key, 1536, 10)
	if err != nil {
		t.Fatalf("NewEngine() failed: %v", err)
	}
	query := Query{
		TenantID:           "tenant-a",
		CredentialTenantID: "tenant-b",
		Namespace:          namespace,
		Embedding:          []float32{0.1},
		TopK:               1,
	}
	if _, err := engine.Retrieve(context.Background(), query); err == nil {
		t.Fatal("Engine.Retrieve() accepted cross-tenant credential")
	}
	query.CredentialTenantID = "tenant-a"
	query.Embedding = []float32{float32(math.NaN())}
	if _, err := engine.Retrieve(context.Background(), query); err == nil {
		t.Fatal("Engine.Retrieve() accepted NaN embedding")
	}
}

func TestNamespaceBindingIsKeyedAndExact(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	namespaceA, err := NamespaceForTenant(key, "tenant-a")
	if err != nil {
		t.Fatalf("NamespaceForTenant() failed: %v", err)
	}
	namespaceB, err := NamespaceForTenant(key, "tenant-b")
	if err != nil {
		t.Fatalf("NamespaceForTenant() failed: %v", err)
	}
	if namespaceA == namespaceB {
		t.Fatal("NamespaceForTenant() produced equal cross-tenant namespaces")
	}
	if err := VerifyTenantNamespace(key, "tenant-a", namespaceB); err == nil {
		t.Fatal("VerifyTenantNamespace() accepted wrong tenant namespace")
	}
	if err := VerifyTenantNamespace(key, "tenant-a", namespaceA); err != nil {
		t.Fatalf("VerifyTenantNamespace() failed: %v", err)
	}
}

func TestEngineHonorsCancellation(t *testing.T) {
	key := []byte("01234567890123456789012345678901")
	namespace, err := NamespaceForTenant(key, "tenant-a")
	if err != nil {
		t.Fatalf("NamespaceForTenant() failed: %v", err)
	}
	engine, err := NewEngine(backendStub{}, key, 1536, 10)
	if err != nil {
		t.Fatalf("NewEngine() failed: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = engine.Retrieve(ctx, Query{
		TenantID:           "tenant-a",
		CredentialTenantID: "tenant-a",
		Namespace:          namespace,
		Embedding:          []float32{0.1},
		TopK:               1,
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Engine.Retrieve() error = %v, want %v", err, context.Canceled)
	}
}
