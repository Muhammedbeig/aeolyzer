package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGuestIdentityRoundTrip(t *testing.T) {
	t.Parallel()

	identity, err := NewGuestIdentity(bytes.Repeat([]byte{1}, 32), false)
	if err != nil {
		t.Fatal(err)
	}
	firstResponse := httptest.NewRecorder()
	firstRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	firstID := identity.Resolve(firstResponse, firstRequest)
	cookies := firstResponse.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("Resolve() cookies = %d, want 1", len(cookies))
	}

	secondResponse := httptest.NewRecorder()
	secondRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	secondRequest.AddCookie(cookies[0])
	secondID := identity.Resolve(secondResponse, secondRequest)
	if secondID != firstID {
		t.Fatalf("Resolve() ID = %q, want %q", secondID, firstID)
	}
	if len(secondResponse.Result().Cookies()) != 0 {
		t.Fatal("Resolve() reset a valid cookie")
	}
}

func TestGuestIdentityRejectsTamperedCookie(t *testing.T) {
	t.Parallel()

	identity, err := NewGuestIdentity(bytes.Repeat([]byte{2}, 32), true)
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.AddCookie(&http.Cookie{Name: guestCookieName, Value: "tampered"})
	userID := identity.Resolve(response, request)
	if userID == "" {
		t.Fatal("Resolve() did not create a replacement identity")
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].Secure || !cookies[0].HttpOnly {
		t.Fatal("Resolve() did not issue a secure replacement cookie")
	}
}
