package executionauth

import (
	"errors"
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	t.Parallel()

	key := []byte("01234567890123456789012345678901")
	now := time.Unix(100, 0)
	claims := Claims{
		TraceID:   "trace-1",
		SessionID: "guest-1",
		Operation: "inspect_public_site",
		TargetURL: "https://example.com/",
		MaxBytes:  1024,
		ExpiresAt: now.Add(time.Minute).Unix(),
	}

	token, err := Sign(key, claims)
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	got, err := Verify(key, token, now)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if got != claims {
		t.Fatalf("Verify() = %#v, want %#v", got, claims)
	}
}

func TestVerifyRejectsExpiredToken(t *testing.T) {
	t.Parallel()

	key := []byte("01234567890123456789012345678901")
	token, err := Sign(key, Claims{
		TraceID:   "trace-1",
		SessionID: "guest-1",
		Operation: "inspect_public_site",
		TargetURL: "https://example.com/",
		MaxBytes:  1024,
		ExpiresAt: 100,
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	_, err = Verify(key, token, time.Unix(100, 0))
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Verify() error = %v, want %v", err, ErrInvalidToken)
	}
}
