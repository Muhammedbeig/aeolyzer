package executionauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var ErrInvalidToken = errors.New("invalid execution authorization")

type Claims struct {
	TraceID   string `json:"trace_id"`
	SessionID string `json:"session_id"`
	Operation string `json:"operation"`
	TargetURL string `json:"target_url"`
	MaxBytes  int64  `json:"max_bytes"`
	ExpiresAt int64  `json:"expires_at"`
}

func Sign(key []byte, claims Claims) (string, error) {
	// Enforcing strict 256-bit key length guarantees cryptographic margin for HMAC-SHA256.
	// Eager claim validation short-circuits to avoid burning CPU cycles on malformed inputs.
	if len(key) < 32 || !validClaims(claims) {
		return "", ErrInvalidToken
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", ErrInvalidToken
	}
	// RawURLEncoding drops padding characters, minimizing wire-size and parsing complexity.
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := signatureFor(key, encodedPayload)

	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

func Verify(key []byte, token string, now time.Time) (Claims, error) {
	if len(key) < 32 {
		return Claims{}, ErrInvalidToken
	}

	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, ErrInvalidToken
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	// Constant-time comparison (hmac.Equal) neutralizes timing side-channels during signature verification.
	if err != nil || !hmac.Equal(signature, signatureFor(key, parts[0])) {
		return Claims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil || !validClaims(claims) {
		return Claims{}, ErrInvalidToken
	}
	// Time check relies on explicit, deterministic clock injection (now) to support testability and edge-case replay.
	if now.Unix() >= claims.ExpiresAt {
		return Claims{}, ErrInvalidToken
	}

	return claims, nil
}

func signatureFor(key []byte, payload string) []byte {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(payload))
	return mac.Sum(nil)
}

func validClaims(claims Claims) bool {
	return claims.TraceID != "" &&
		claims.SessionID != "" &&
		claims.Operation != "" &&
		claims.TargetURL != "" &&
		claims.MaxBytes > 0 &&
		claims.ExpiresAt > 0
}
