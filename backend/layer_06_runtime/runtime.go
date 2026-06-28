package runtime

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"aeolyzer/internal/executionauth"
)

var (
	ErrDeniedTarget     = errors.New("target is denied")
	ErrInvalidExecution = errors.New("invalid execution request")
)

type ExecutionRequest struct {
	TraceID       string
	SessionID     string
	Operation     string
	TargetURL     string
	MaxBytes      int64
	Authorization string
}

type ExecutionResult struct {
	Title                string
	Description          string
	IconURL              string
	Category             string
	CandidateCompetitors []string
}

type Adapter interface {
	Inspect(ctx context.Context, targetURL string, maxBytes int64) (ExecutionResult, error)
}

type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

type Executor struct {
	resolver   Resolver
	adapter    Adapter
	signingKey []byte
	now        func() time.Time
}

func NewExecutor(resolver Resolver, adapter Adapter, signingKey []byte, now func() time.Time) *Executor {
	return &Executor{
		resolver:   resolver,
		adapter:    adapter,
		signingKey: append([]byte(nil), signingKey...),
		now:        now,
	}
}

func (e *Executor) Execute(ctx context.Context, request ExecutionRequest) (ExecutionResult, error) {
	if e == nil || e.resolver == nil || e.adapter == nil || e.now == nil || len(e.signingKey) < 32 {
		return ExecutionResult{}, errors.New("runtime executor is not configured")
	}
	if request.TraceID == "" || request.SessionID == "" || request.Operation != "inspect_public_site" {
		return ExecutionResult{}, ErrInvalidExecution
	}
	if request.MaxBytes <= 0 || request.MaxBytes > 2<<20 {
		return ExecutionResult{}, ErrInvalidExecution
	}
	claims, err := executionauth.Verify(e.signingKey, request.Authorization, e.now())
	if err != nil ||
		claims.TraceID != request.TraceID ||
		claims.SessionID != request.SessionID ||
		claims.Operation != request.Operation ||
		claims.TargetURL != request.TargetURL ||
		claims.MaxBytes != request.MaxBytes {
		return ExecutionResult{}, ErrInvalidExecution
	}
	if err := validatePublicTarget(ctx, e.resolver, request.TargetURL); err != nil {
		return ExecutionResult{}, err
	}

	return e.adapter.Inspect(ctx, request.TargetURL, request.MaxBytes)
}

func validatePublicTarget(ctx context.Context, resolver Resolver, rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Hostname() == "" {
		return ErrDeniedTarget
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrDeniedTarget
	}

	host := strings.ToLower(parsed.Hostname())
	if host == "localhost" || strings.HasSuffix(host, ".localhost") {
		return ErrDeniedTarget
	}

	addresses, err := resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf("resolve target: %w", err)
	}
	if len(addresses) == 0 {
		return ErrDeniedTarget
	}
	for _, address := range addresses {
		if !isPublicIP(address.IP) {
			return ErrDeniedTarget
		}
	}

	return nil
}

func isPublicIP(ip net.IP) bool {
	if ip == nil || ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return false
	}

	return true
}
