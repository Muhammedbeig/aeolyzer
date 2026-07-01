package httpapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	guestCookieName = "aeolyzer_guest"
	guestCookieDays = 30
)

type GuestIdentity struct {
	signingKey []byte
	secure     bool
	now        func() time.Time
}

func NewGuestIdentity(dataKey []byte, secure bool) (*GuestIdentity, error) {
	if len(dataKey) != 32 {
		return nil, errors.New("guest identity key must be 32 bytes")
	}
	mac := hmac.New(sha256.New, dataKey)
	_, _ = mac.Write([]byte("aeolyzer-guest-cookie-v1"))
	return &GuestIdentity{
		signingKey: mac.Sum(nil),
		secure:     secure,
		now:        time.Now,
	}, nil
}

func (i *GuestIdentity) Resolve(response http.ResponseWriter, request *http.Request) string {
	if cookie, err := request.Cookie(guestCookieName); err == nil {
		if userID, ok := i.verify(cookie.Value); ok {
			return userID
		}
	}
	userID := uuid.NewString()
	http.SetCookie(response, &http.Cookie{
		Name:     guestCookieName,
		Value:    i.sign(userID),
		Path:     "/",
		MaxAge:   guestCookieDays * 24 * 60 * 60,
		Expires:  i.now().UTC().Add(guestCookieDays * 24 * time.Hour),
		HttpOnly: true,
		Secure:   i.secure,
		SameSite: http.SameSiteStrictMode,
	})
	return userID
}

func (i *GuestIdentity) sign(userID string) string {
	mac := hmac.New(sha256.New, i.signingKey)
	_, _ = mac.Write([]byte(userID))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return userID + "." + signature
}

func (i *GuestIdentity) verify(value string) (string, bool) {
	userID, signature, ok := strings.Cut(value, ".")
	if !ok {
		return "", false
	}
	if _, err := uuid.Parse(userID); err != nil {
		return "", false
	}
	provided, err := base64.RawURLEncoding.DecodeString(signature)
	if err != nil {
		return "", false
	}
	mac := hmac.New(sha256.New, i.signingKey)
	_, _ = mac.Write([]byte(userID))
	return userID, hmac.Equal(provided, mac.Sum(nil))
}
