package intake

import (
	"errors"
	"testing"
	"time"
)

func TestInspectSite(t *testing.T) {
	t.Parallel()

	service := newTestService()
	tests := []struct {
		name    string
		rawURL  string
		wantURL string
		wantErr error
	}{
		{name: "adds scheme", rawURL: "example.com", wantURL: "https://example.com/"},
		{name: "removes query and fragment", rawURL: "https://example.com/page?q=secret#part", wantURL: "https://example.com/page"},
		{name: "rejects credentials", rawURL: "https://user:pass@example.com", wantErr: ErrInvalidURL},
		{name: "rejects unsupported scheme", rawURL: "file:///tmp/site", wantErr: ErrInvalidURL},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := service.InspectSite(SiteInspectionInput{SessionID: "guest", URL: tt.rawURL})
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("InspectSite() error = %v, want %v", err, tt.wantErr)
			}
			if got.CanonicalURL != tt.wantURL {
				t.Fatalf("InspectSite() URL = %q, want %q", got.CanonicalURL, tt.wantURL)
			}
		})
	}
}

func TestCompleteOnboardingRejectsUnknownReach(t *testing.T) {
	t.Parallel()

	service := newTestService()
	_, err := service.CompleteOnboarding(OnboardingInput{
		SessionID:   "guest",
		AccountType: AccountTypeBrand,
		Domain:      "example.com",
		BrandName:   "Example",
		Reach:       Reach("everywhere"),
		CountryCode: "PK",
		CountryName: "Pakistan",
		Language:    "English (UK)",
	})
	if !errors.Is(err, ErrInvalidProfile) {
		t.Fatalf("CompleteOnboarding() error = %v, want %v", err, ErrInvalidProfile)
	}
}

func newTestService() *Service {
	return NewService(
		func() string { return "trace-1" },
		[]byte("01234567890123456789012345678901"),
		func() time.Time { return time.Unix(100, 0) },
	)
}
