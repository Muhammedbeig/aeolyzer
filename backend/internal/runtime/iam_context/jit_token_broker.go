package iamcontext

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

var (
	// ErrCredentialPolicyDenied indicates that a JIT request exceeds policy.
	ErrCredentialPolicyDenied = errors.New("credential request denied by policy")
	// ErrCredentialExpired indicates that a credential is no longer usable.
	ErrCredentialExpired = errors.New("credential is expired")
	// ErrCredentialRevoked indicates that a credential was explicitly revoked.
	ErrCredentialRevoked = errors.New("credential is revoked")
	// ErrCredentialBindingMismatch indicates trace, tenant, audience, or connector
	// confusion.
	ErrCredentialBindingMismatch = errors.New("credential binding mismatch")
)

// CredentialRequest is an already-authorized JIT credential request.
type CredentialRequest struct {
	TraceID          string        `json:"trace_id"`
	TenantID         string        `json:"tenant_id"`
	PolicyDecisionID string        `json:"policy_decision_id"`
	ConnectorID      string        `json:"connector_id"`
	Audience         string        `json:"audience"`
	Scopes           []string      `json:"scopes"`
	TTL              time.Duration `json:"ttl"`
}

// IssuerRequest is the minimum request sent to an external credential issuer.
type IssuerRequest struct {
	TenantID    string
	ConnectorID string
	Audience    string
	Scopes      []string
	TTL         time.Duration
}

// IssuedCredential contains secret material returned by a credential provider.
// Value is intentionally excluded from JSON.
type IssuedCredential struct {
	ID        string
	Value     string `json:"-"`
	ExpiresAt time.Time
}

// CredentialIssuer mints and revokes provider credentials.
type CredentialIssuer interface {
	Issue(context.Context, IssuerRequest) (IssuedCredential, error)
	Revoke(context.Context, string) error
}

// ConnectorScopePolicy defines the exact audience, scopes, and TTL for one
// connector.
type ConnectorScopePolicy struct {
	Audience      string
	AllowedScopes []string
	MaxTTL        time.Duration
}

// ScopePolicy is a fail-closed connector credential allowlist.
type ScopePolicy struct {
	Connectors map[string]ConnectorScopePolicy
}

// TokenReference is safe to pass across runtime components.
type TokenReference struct {
	ReferenceID string    `json:"reference_id"`
	ExpiresAt   time.Time `json:"expires_at"`
	ConnectorID string    `json:"connector_id"`
	Audience    string    `json:"audience"`
	Scopes      []string  `json:"scopes"`
}

type activeCredential struct {
	issued      IssuedCredential
	traceID     string
	tenantID    string
	connectorID string
	audience    string
	scopes      []string
	revoked     bool
}

// Broker issues, binds, resolves, and revokes short-lived credentials.
type Broker struct {
	mu          sync.Mutex
	issuer      CredentialIssuer
	policy      ScopePolicy
	now         func() time.Time
	credentials map[string]*activeCredential
}

// NewBroker constructs a fail-closed JIT credential broker.
func NewBroker(
	issuer CredentialIssuer,
	policy ScopePolicy,
	now func() time.Time,
) (*Broker, error) {
	if issuer == nil || now == nil || len(policy.Connectors) == 0 {
		return nil, errors.New("jit credential broker is not configured")
	}
	for connectorID, connector := range policy.Connectors {
		if connectorID == "" ||
			connector.Audience == "" ||
			len(connector.AllowedScopes) == 0 ||
			connector.MaxTTL <= 0 ||
			connector.MaxTTL > time.Hour {
			return nil, errors.New("jit credential policy is invalid")
		}
	}
	return &Broker{
		issuer:      issuer,
		policy:      cloneScopePolicy(policy),
		now:         now,
		credentials: make(map[string]*activeCredential),
	}, nil
}

// Issue validates authorization bindings and mints a short-lived credential.
func (b *Broker) Issue(
	ctx context.Context,
	request CredentialRequest,
) (TokenReference, error) {
	if b == nil || b.issuer == nil {
		return TokenReference{}, errors.New("jit credential broker is not configured")
	}
	policy, scopes, err := b.validateRequest(request)
	if err != nil {
		return TokenReference{}, err
	}
	if err := ctx.Err(); err != nil {
		return TokenReference{}, fmt.Errorf("issue jit credential: %w", err)
	}

	issued, err := b.issuer.Issue(ctx, IssuerRequest{
		TenantID:    request.TenantID,
		ConnectorID: request.ConnectorID,
		Audience:    request.Audience,
		Scopes:      scopes,
		TTL:         request.TTL,
	})
	if err != nil {
		return TokenReference{}, fmt.Errorf("issue provider credential: %w", err)
	}
	now := b.now()
	if issued.ID == "" ||
		issued.Value == "" ||
		!issued.ExpiresAt.After(now) ||
		issued.ExpiresAt.After(now.Add(request.TTL)) ||
		issued.ExpiresAt.After(now.Add(policy.MaxTTL)) {
		if issued.ID != "" {
			_ = b.issuer.Revoke(ctx, issued.ID)
		}
		return TokenReference{}, errors.New("credential issuer returned invalid bounds")
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.credentials[issued.ID]; exists {
		_ = b.issuer.Revoke(ctx, issued.ID)
		return TokenReference{}, errors.New("credential issuer reused an active id")
	}
	b.credentials[issued.ID] = &activeCredential{
		issued:      issued,
		traceID:     request.TraceID,
		tenantID:    request.TenantID,
		connectorID: request.ConnectorID,
		audience:    request.Audience,
		scopes:      append([]string(nil), scopes...),
	}
	return TokenReference{
		ReferenceID: issued.ID,
		ExpiresAt:   issued.ExpiresAt,
		ConnectorID: request.ConnectorID,
		Audience:    request.Audience,
		Scopes:      append([]string(nil), scopes...),
	}, nil
}

// Resolve returns secret material only after exact binding validation. Callers
// must attach it directly to the approved connector and must not log it.
func (b *Broker) Resolve(
	traceID, tenantID, connectorID, audience, referenceID string,
) (string, error) {
	if b == nil {
		return "", errors.New("jit credential broker is not configured")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	credential, found := b.credentials[referenceID]
	if !found {
		return "", ErrCredentialBindingMismatch
	}
	if credential.revoked {
		return "", ErrCredentialRevoked
	}
	if !b.now().Before(credential.issued.ExpiresAt) {
		return "", ErrCredentialExpired
	}
	if credential.traceID != traceID ||
		credential.tenantID != tenantID ||
		credential.connectorID != connectorID ||
		credential.audience != audience {
		return "", ErrCredentialBindingMismatch
	}
	return credential.issued.Value, nil
}

// Revoke revokes one credential and makes local denial immediate.
func (b *Broker) Revoke(ctx context.Context, referenceID string) error {
	if b == nil {
		return errors.New("jit credential broker is not configured")
	}
	b.mu.Lock()
	credential, found := b.credentials[referenceID]
	if !found {
		b.mu.Unlock()
		return ErrCredentialBindingMismatch
	}
	if credential.revoked {
		b.mu.Unlock()
		return nil
	}
	credential.revoked = true
	b.mu.Unlock()

	if err := b.issuer.Revoke(ctx, credential.issued.ID); err != nil {
		return fmt.Errorf("revoke provider credential: %w", err)
	}
	return nil
}

// RevokeTrace revokes every credential bound to a trace.
func (b *Broker) RevokeTrace(ctx context.Context, traceID string) error {
	if traceID == "" {
		return errors.New("trace id is required")
	}
	b.mu.Lock()
	var ids []string
	for id, credential := range b.credentials {
		if credential.traceID == traceID && !credential.revoked {
			credential.revoked = true
			ids = append(ids, id)
		}
	}
	b.mu.Unlock()

	var firstErr error
	for _, id := range ids {
		if err := b.issuer.Revoke(ctx, id); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("revoke trace credential: %w", err)
		}
	}
	return firstErr
}

func (b *Broker) validateRequest(
	request CredentialRequest,
) (ConnectorScopePolicy, []string, error) {
	if request.TraceID == "" ||
		request.TenantID == "" ||
		request.PolicyDecisionID == "" ||
		request.ConnectorID == "" ||
		request.Audience == "" ||
		len(request.Scopes) == 0 ||
		request.TTL <= 0 {
		return ConnectorScopePolicy{}, nil, ErrCredentialPolicyDenied
	}
	policy, found := b.policy.Connectors[request.ConnectorID]
	if !found ||
		request.Audience != policy.Audience ||
		request.TTL > policy.MaxTTL {
		return ConnectorScopePolicy{}, nil, ErrCredentialPolicyDenied
	}
	allowed := make(map[string]struct{}, len(policy.AllowedScopes))
	for _, scope := range policy.AllowedScopes {
		allowed[scope] = struct{}{}
	}
	unique := make(map[string]struct{}, len(request.Scopes))
	scopes := make([]string, 0, len(request.Scopes))
	for _, scope := range request.Scopes {
		if _, ok := allowed[scope]; !ok {
			return ConnectorScopePolicy{}, nil, ErrCredentialPolicyDenied
		}
		if _, duplicate := unique[scope]; duplicate {
			continue
		}
		unique[scope] = struct{}{}
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)
	return policy, scopes, nil
}

func cloneScopePolicy(policy ScopePolicy) ScopePolicy {
	result := ScopePolicy{Connectors: make(map[string]ConnectorScopePolicy, len(policy.Connectors))}
	for id, connector := range policy.Connectors {
		copy := connector
		copy.AllowedScopes = append([]string(nil), connector.AllowedScopes...)
		result.Connectors[id] = copy
	}
	return result
}
