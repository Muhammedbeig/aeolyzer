package intake

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"aeolyzer/internal/executionauth"
)

const (
	maxBrandLength      = 120
	maxCountryLength    = 80
	maxLanguageLength   = 80
	maxCompetitorCount  = 5
	maxCompetitorLength = 180
)

var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrInvalidURL     = errors.New("invalid website URL")
	ErrInvalidProfile = errors.New("invalid onboarding profile")
)

type AccountType string

const (
	AccountTypeBrand  AccountType = "brand"
	AccountTypeAgency AccountType = "agency"
)

type Reach string

const (
	ReachGlobal     Reach = "global"
	ReachPrimary    Reach = "primary-market"
	ReachNationwide Reach = "nationwide"
	ReachRegional   Reach = "regional"
	ReachLocal      Reach = "local"
)

type SiteInspectionInput struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

type SiteInspectionDecision struct {
	TraceID       string
	SessionID     string
	CanonicalURL  string
	Operation     string
	Authorization string
}

type OnboardingInput struct {
	SessionID   string      `json:"session_id"`
	AccountType AccountType `json:"account_type"`
	Domain      string      `json:"domain"`
	BrandName   string      `json:"brand_name"`
	Reach       Reach       `json:"reach"`
	CountryCode string      `json:"country_code"`
	CountryName string      `json:"country_name"`
	Language    string      `json:"language"`
	Competitors []string    `json:"competitors"`
}

type OnboardingDecision struct {
	TraceID string
	Profile ProjectProfile
}

type ProjectProfile struct {
	SessionID   string      `json:"session_id"`
	AccountType AccountType `json:"account_type"`
	Domain      string      `json:"domain"`
	BrandName   string      `json:"brand_name"`
	Reach       Reach       `json:"reach"`
	CountryCode string      `json:"country_code"`
	CountryName string      `json:"country_name"`
	Language    string      `json:"language"`
	Competitors []string    `json:"competitors"`
}

type Service struct {
	newTraceID func() string
	signingKey []byte
	now        func() time.Time
}

func NewService(newTraceID func() string, signingKey []byte, now func() time.Time) *Service {
	return &Service{
		newTraceID: newTraceID,
		signingKey: append([]byte(nil), signingKey...),
		now:        now,
	}
}

func (s *Service) InspectSite(input SiteInspectionInput) (SiteInspectionDecision, error) {
	if s == nil || s.newTraceID == nil || s.now == nil || len(s.signingKey) < 32 {
		return SiteInspectionDecision{}, errors.New("intake service is not configured")
	}

	sessionID, err := cleanText("session_id", input.SessionID, 80, true)
	if err != nil {
		return SiteInspectionDecision{}, err
	}

	canonicalURL, err := canonicalizePublicURL(input.URL)
	if err != nil {
		return SiteInspectionDecision{}, err
	}

	traceID := s.newTraceID()
	authorization, err := executionauth.Sign(s.signingKey, executionauth.Claims{
		TraceID:   traceID,
		SessionID: sessionID,
		Operation: "inspect_public_site",
		TargetURL: canonicalURL,
		MaxBytes:  2 << 20,
		ExpiresAt: s.now().Add(time.Minute).Unix(),
	})
	if err != nil {
		return SiteInspectionDecision{}, errors.New("create execution authorization")
	}

	return SiteInspectionDecision{
		TraceID:       traceID,
		SessionID:     sessionID,
		CanonicalURL:  canonicalURL,
		Operation:     "inspect_public_site",
		Authorization: authorization,
	}, nil
}

func (s *Service) CompleteOnboarding(input OnboardingInput) (OnboardingDecision, error) {
	if s == nil || s.newTraceID == nil {
		return OnboardingDecision{}, errors.New("intake service is not configured")
	}
	if input.AccountType != AccountTypeBrand && input.AccountType != AccountTypeAgency {
		return OnboardingDecision{}, fmt.Errorf("%w: unsupported account type", ErrInvalidProfile)
	}
	if !validReach(input.Reach) {
		return OnboardingDecision{}, fmt.Errorf("%w: unsupported reach", ErrInvalidProfile)
	}

	sessionID, err := cleanText("session_id", input.SessionID, 80, true)
	if err != nil {
		return OnboardingDecision{}, err
	}
	domain, err := canonicalizePublicURL(input.Domain)
	if err != nil {
		return OnboardingDecision{}, err
	}
	brandName, err := cleanText("brand_name", input.BrandName, maxBrandLength, true)
	if err != nil {
		return OnboardingDecision{}, err
	}
	countryCode, err := cleanText("country_code", strings.ToUpper(input.CountryCode), 2, true)
	if err != nil || len(countryCode) != 2 {
		return OnboardingDecision{}, fmt.Errorf("%w: invalid country code", ErrInvalidProfile)
	}
	countryName, err := cleanText("country_name", input.CountryName, maxCountryLength, true)
	if err != nil {
		return OnboardingDecision{}, err
	}
	language, err := cleanText("language", input.Language, maxLanguageLength, true)
	if err != nil {
		return OnboardingDecision{}, err
	}
	competitors, err := cleanCompetitors(input.Competitors, domain)
	if err != nil {
		return OnboardingDecision{}, err
	}

	return OnboardingDecision{
		TraceID: s.newTraceID(),
		Profile: ProjectProfile{
			SessionID:   sessionID,
			AccountType: input.AccountType,
			Domain:      domain,
			BrandName:   brandName,
			Reach:       input.Reach,
			CountryCode: countryCode,
			CountryName: countryName,
			Language:    language,
			Competitors: competitors,
		},
	}, nil
}

func canonicalizePublicURL(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", ErrInvalidURL
	}
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}

	parsed, err := url.Parse(value)
	if err != nil || parsed.Hostname() == "" {
		return "", ErrInvalidURL
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", ErrInvalidURL
	}
	if parsed.User != nil {
		return "", ErrInvalidURL
	}

	parsed.Fragment = ""
	parsed.RawQuery = ""
	parsed.Path = strings.TrimRight(parsed.EscapedPath(), "/")
	if parsed.Path == "" {
		parsed.Path = "/"
	}

	return parsed.String(), nil
}

func cleanCompetitors(values []string, ownDomain string) ([]string, error) {
	if len(values) > maxCompetitorCount {
		return nil, fmt.Errorf("%w: too many competitors", ErrInvalidProfile)
	}

	ownURL, _ := url.Parse(ownDomain)
	ownHost := strings.TrimPrefix(strings.ToLower(ownURL.Hostname()), "www.")
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		cleaned, err := cleanText("competitor", value, maxCompetitorLength, false)
		if err != nil {
			return nil, err
		}
		if cleaned == "" {
			continue
		}
		key := strings.ToLower(cleaned)
		if key == ownHost || key == strings.TrimPrefix(ownDomain, "https://") {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, cleaned)
	}

	return result, nil
}

func cleanText(field, value string, limit int, required bool) (string, error) {
	cleaned := strings.TrimSpace(strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, value))
	if required && cleaned == "" {
		return "", fmt.Errorf("%w: %s is required", ErrInvalidInput, field)
	}
	if !utf8.ValidString(cleaned) || utf8.RuneCountInString(cleaned) > limit {
		return "", fmt.Errorf("%w: %s exceeds its limit", ErrInvalidInput, field)
	}
	return cleaned, nil
}

func validReach(reach Reach) bool {
	switch reach {
	case ReachGlobal, ReachPrimary, ReachNationwide, ReachRegional, ReachLocal:
		return true
	default:
		return false
	}
}
