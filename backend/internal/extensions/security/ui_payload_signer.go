package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"aeolyzer/internal/extensions"
)

// SignFrame signs a complete, expiring A2UI frame.
func SignFrame(
	key []byte,
	frame extensions.A2UIFrame,
	now time.Time,
) (extensions.A2UIFrame, error) {
	if len(key) < 32 || now.IsZero() {
		return extensions.A2UIFrame{}, errors.New("ui signing context is invalid")
	}
	if err := validateFrameExpiry(frame.ExpiresAt, now); err != nil {
		return extensions.A2UIFrame{}, err
	}
	frame.Signature = ""
	payload, err := json.Marshal(frame)
	if err != nil {
		return extensions.A2UIFrame{}, fmt.Errorf("encode ui frame: %w", err)
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(payload)
	frame.Signature = base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return frame, nil
}

// VerifyFrame verifies signature and expiry.
func VerifyFrame(key []byte, frame extensions.A2UIFrame, now time.Time) error {
	if len(key) < 32 || frame.Signature == "" || now.IsZero() {
		return errors.New("ui verification context is invalid")
	}
	if err := validateFrameExpiry(frame.ExpiresAt, now); err != nil {
		return err
	}
	provided, err := base64.RawURLEncoding.DecodeString(frame.Signature)
	if err != nil {
		return errors.New("ui frame signature is malformed")
	}
	unsigned := frame
	unsigned.Signature = ""
	payload, err := json.Marshal(unsigned)
	if err != nil {
		return fmt.Errorf("encode ui frame: %w", err)
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(payload)
	if !hmac.Equal(provided, mac.Sum(nil)) {
		return errors.New("ui frame signature verification failed")
	}
	return nil
}

func validateFrameExpiry(raw string, now time.Time) error {
	expiresAt, err := time.Parse(time.RFC3339, raw)
	if err != nil ||
		!now.Before(expiresAt) ||
		expiresAt.After(now.Add(15*time.Minute)) {
		return errors.New("ui frame expiry is invalid")
	}
	return nil
}
