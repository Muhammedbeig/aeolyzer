package iamcontext

import (
	"context"
	"errors"
	"testing"
	"time"
)

type issuerStub struct {
	now     time.Time
	issued  int
	revoked []string
}

func (i *issuerStub) Issue(_ context.Context, request IssuerRequest) (IssuedCredential, error) {
	i.issued++
	return IssuedCredential{
		ID:        "credential-1",
		Value:     "secret-token",
		ExpiresAt: i.now.Add(request.TTL),
	}, nil
}

func (i *issuerStub) Revoke(_ context.Context, id string) error {
	i.revoked = append(i.revoked, id)
	return nil
}

func TestBrokerBindsAndRevokesCredential(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	issuer := &issuerStub{now: now}
	broker := newTestBroker(t, issuer, &now)
	reference, err := broker.Issue(context.Background(), validCredentialRequest())
	if err != nil {
		t.Fatalf("Broker.Issue() failed: %v", err)
	}
	token, err := broker.Resolve(
		"trace-1",
		"tenant-1",
		"search-console",
		"https://connector.internal",
		reference.ReferenceID,
	)
	if err != nil {
		t.Fatalf("Broker.Resolve() failed: %v", err)
	}
	if token != "secret-token" {
		t.Fatal("Broker.Resolve() returned wrong secret")
	}
	if err := broker.Revoke(context.Background(), reference.ReferenceID); err != nil {
		t.Fatalf("Broker.Revoke() failed: %v", err)
	}
	if _, err := broker.Resolve(
		"trace-1",
		"tenant-1",
		"search-console",
		"https://connector.internal",
		reference.ReferenceID,
	); !errors.Is(err, ErrCredentialRevoked) {
		t.Fatalf("Broker.Resolve() error = %v, want %v", err, ErrCredentialRevoked)
	}
}

func TestBrokerRejectsScopeAudienceTenantAndExpiryConfusion(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	issuer := &issuerStub{now: now}
	broker := newTestBroker(t, issuer, &now)

	for name, mutate := range map[string]func(*CredentialRequest){
		"scope": func(r *CredentialRequest) {
			r.Scopes = []string{"write"}
		},
		"audience": func(r *CredentialRequest) {
			r.Audience = "https://attacker.invalid"
		},
		"ttl": func(r *CredentialRequest) {
			r.TTL = 16 * time.Minute
		},
	} {
		t.Run(name, func(t *testing.T) {
			request := validCredentialRequest()
			mutate(&request)
			if _, err := broker.Issue(context.Background(), request); !errors.Is(err, ErrCredentialPolicyDenied) {
				t.Fatalf("Broker.Issue() error = %v, want %v", err, ErrCredentialPolicyDenied)
			}
		})
	}

	reference, err := broker.Issue(context.Background(), validCredentialRequest())
	if err != nil {
		t.Fatalf("Broker.Issue() failed: %v", err)
	}
	if _, err := broker.Resolve(
		"trace-1",
		"tenant-2",
		"search-console",
		"https://connector.internal",
		reference.ReferenceID,
	); !errors.Is(err, ErrCredentialBindingMismatch) {
		t.Fatalf("Broker.Resolve() error = %v, want binding mismatch", err)
	}
	now = now.Add(11 * time.Minute)
	if _, err := broker.Resolve(
		"trace-1",
		"tenant-1",
		"search-console",
		"https://connector.internal",
		reference.ReferenceID,
	); !errors.Is(err, ErrCredentialExpired) {
		t.Fatalf("Broker.Resolve() error = %v, want expired", err)
	}
}

func newTestBroker(
	t *testing.T,
	issuer CredentialIssuer,
	now *time.Time,
) *Broker {
	t.Helper()
	broker, err := NewBroker(issuer, ScopePolicy{
		Connectors: map[string]ConnectorScopePolicy{
			"search-console": {
				Audience:      "https://connector.internal",
				AllowedScopes: []string{"read"},
				MaxTTL:        15 * time.Minute,
			},
		},
	}, func() time.Time { return *now })
	if err != nil {
		t.Fatalf("NewBroker() failed: %v", err)
	}
	return broker
}

func validCredentialRequest() CredentialRequest {
	return CredentialRequest{
		TraceID:          "trace-1",
		TenantID:         "tenant-1",
		PolicyDecisionID: "decision-1",
		ConnectorID:      "search-console",
		Audience:         "https://connector.internal",
		Scopes:           []string{"read"},
		TTL:              10 * time.Minute,
	}
}
