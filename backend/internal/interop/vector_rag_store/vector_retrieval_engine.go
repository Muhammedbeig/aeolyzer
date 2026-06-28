package vectorragstore

import (
	"context"
	"errors"
	"fmt"
	"math"

	datasecuritymesh "aeolyzer/internal/interop/data_security_mesh"
)

// Query is an already-authorized tenant-bound vector query.
type Query struct {
	TenantID           string
	CredentialTenantID string
	Namespace          string
	Embedding          []float32
	TopK               int
}

// Record is one provenance-bound retrieval result.
type Record struct {
	ID            string         `json:"id"`
	Namespace     string         `json:"namespace"`
	ContentHash   string         `json:"content_hash"`
	ProvenanceRef string         `json:"provenance_ref"`
	Score         float64        `json:"score"`
	Fields        map[string]any `json:"fields"`
}

// Backend performs storage-specific retrieval in one exact namespace.
type Backend interface {
	Search(context.Context, string, []float32, int) ([]Record, error)
}

// Engine validates tenant partitions before and after retrieval.
type Engine struct {
	backend      Backend
	namespaceKey []byte
	maxDimension int
	maxTopK      int
}

// NewEngine constructs a tenant-isolated retrieval engine.
func NewEngine(
	backend Backend,
	namespaceKey []byte,
	maxDimension, maxTopK int,
) (*Engine, error) {
	if backend == nil ||
		len(namespaceKey) < 32 ||
		maxDimension < 1 ||
		maxDimension > 8192 ||
		maxTopK < 1 ||
		maxTopK > 100 {
		return nil, errors.New("vector retrieval engine is not configured")
	}
	return &Engine{
		backend:      backend,
		namespaceKey: append([]byte(nil), namespaceKey...),
		maxDimension: maxDimension,
		maxTopK:      maxTopK,
	}, nil
}

// Retrieve performs bounded retrieval and rejects any cross-namespace result.
func (e *Engine) Retrieve(ctx context.Context, query Query) ([]Record, error) {
	if e == nil || e.backend == nil {
		return nil, errors.New("vector retrieval engine is not configured")
	}
	if err := datasecuritymesh.EnforceTenantBoundary(
		query.TenantID,
		query.CredentialTenantID,
	); err != nil {
		return nil, err
	}
	if err := VerifyTenantNamespace(
		e.namespaceKey,
		query.TenantID,
		query.Namespace,
	); err != nil {
		return nil, err
	}
	if query.TopK < 1 ||
		query.TopK > e.maxTopK ||
		len(query.Embedding) < 1 ||
		len(query.Embedding) > e.maxDimension {
		return nil, errors.New("vector query exceeds policy limits")
	}
	for _, value := range query.Embedding {
		if math.IsNaN(float64(value)) || math.IsInf(float64(value), 0) {
			return nil, errors.New("vector query contains invalid values")
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("retrieve vector evidence: %w", err)
	}
	records, err := e.backend.Search(
		ctx,
		query.Namespace,
		append([]float32(nil), query.Embedding...),
		query.TopK,
	)
	if err != nil {
		return nil, fmt.Errorf("search vector backend: %w", err)
	}
	if len(records) > query.TopK {
		return nil, errors.New("vector backend exceeded top k")
	}
	result := make([]Record, len(records))
	for i, record := range records {
		if record.ID == "" ||
			record.Namespace != query.Namespace ||
			record.ContentHash == "" ||
			record.ProvenanceRef == "" ||
			record.Fields == nil {
			return nil, errors.New("vector backend returned invalid or cross-tenant evidence")
		}
		result[i] = cloneRecord(record)
	}
	return result, nil
}

func cloneRecord(record Record) Record {
	result := record
	result.Fields = make(map[string]any, len(record.Fields))
	for key, value := range record.Fields {
		result.Fields[key] = value
	}
	return result
}
