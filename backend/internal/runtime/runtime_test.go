package runtime

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"aeolyzer/internal/executionauth"
)

type stubResolver struct {
	addresses []net.IPAddr
}

func (r stubResolver) LookupIPAddr(context.Context, string) ([]net.IPAddr, error) {
	return r.addresses, nil
}

type stubAdapter struct{}

func (stubAdapter) Inspect(context.Context, string, int64) (ExecutionResult, error) {
	return ExecutionResult{Title: "Example"}, nil
}

func TestExecutorDeniesPrivateTargets(t *testing.T) {
	// PERFORMANCE: Enable parallel test execution to amortize suite runtime.
	t.Parallel()

	// STATE MANAGEMENT: Use a fixed 256-bit symmetric key for deterministic test evaluation.
	key := []byte("01234567890123456789012345678901")
	now := time.Unix(100, 0)
	authorization, err := executionauth.Sign(key, executionauth.Claims{
		TraceID:   "trace-1",
		SessionID: "guest",
		Operation: "inspect_public_site",
		TargetURL: "http://example.test/",
		MaxBytes:  1024,
		ExpiresAt: now.Add(time.Minute).Unix(),
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	executor := NewExecutor(stubResolver{
		addresses: []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}},
	}, stubAdapter{}, key, func() time.Time { return now })
	_, err = executor.Execute(context.Background(), ExecutionRequest{
		TraceID:       "trace-1",
		SessionID:     "guest",
		Operation:     "inspect_public_site",
		TargetURL:     "http://example.test/",
		MaxBytes:      1024,
		Authorization: authorization,
	})
	if !errors.Is(err, ErrDeniedTarget) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrDeniedTarget)
	}
}

func TestExecutorRejectsTamperedAuthorization(t *testing.T) {
	// PERFORMANCE: Run test concurrently with other stateless validations.
	t.Parallel()

	// STATE MANAGEMENT: Use a fixed 256-bit symmetric key for deterministic test evaluation.
	key := []byte("01234567890123456789012345678901")
	executor := NewExecutor(
		stubResolver{addresses: []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}},
		stubAdapter{},
		key,
		func() time.Time { return time.Unix(100, 0) },
	)
	_, err := executor.Execute(context.Background(), ExecutionRequest{
		TraceID:       "trace-1",
		SessionID:     "guest",
		Operation:     "inspect_public_site",
		TargetURL:     "https://example.com/",
		MaxBytes:      1024,
		Authorization: "tampered",
	})
	if !errors.Is(err, ErrInvalidExecution) {
		t.Fatalf("Execute() error = %v, want %v", err, ErrInvalidExecution)
	}
}
