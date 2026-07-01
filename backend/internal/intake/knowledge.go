package intake

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"unicode"
	"unicode/utf8"

	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
	"aeolyzer/internal/workspace"
)

const (
	maxKnowledgeBodyBytes    = 64 << 10
	maxKnowledgeSummaryRunes = 4_000
)

var (
	ErrKnowledgeApprovalRequired = errors.New("knowledge update approval is required")
	ErrInvalidKnowledgeUpdate    = errors.New("invalid knowledge update")
	ErrUnsafeKnowledgeUpdate     = errors.New("unsafe knowledge update")
)

type ValidatedKnowledgeUpdate struct {
	Document workspace.KnowledgeDocument
	Summary  string
}

func ValidateKnowledgeUpdate(
	document workspace.KnowledgeDocument,
	approved bool,
) (ValidatedKnowledgeUpdate, error) {
	if !approved {
		return ValidatedKnowledgeUpdate{}, ErrKnowledgeApprovalRequired
	}
	if !document.Section.Valid() || payloadCount(document) != 1 {
		return ValidatedKnowledgeUpdate{}, ErrInvalidKnowledgeUpdate
	}

	normalized, summary, err := normalizeKnowledgeDocument(document)
	if err != nil {
		return ValidatedKnowledgeUpdate{}, err
	}
	encoded, err := json.Marshal(normalized)
	if err != nil || len(encoded) > maxKnowledgeBodyBytes {
		return ValidatedKnowledgeUpdate{}, ErrInvalidKnowledgeUpdate
	}
	return ValidatedKnowledgeUpdate{
		Document: normalized,
		Summary:  limitRunes(summary, maxKnowledgeSummaryRunes),
	}, nil
}

func normalizeKnowledgeDocument(
	document workspace.KnowledgeDocument,
) (workspace.KnowledgeDocument, string, error) {
	normalized := workspace.KnowledgeDocument{
		Section: document.Section,
		Version: document.Version,
	}
	switch document.Section {
	case workspace.KnowledgeSectionProfile:
		if document.Profile == nil {
			return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
		}
		name, err := validateKnowledgeText(document.Profile.Name, 120)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		description, err := validateKnowledgeText(document.Profile.Description, 4_000)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		normalized.Profile = &workspace.KnowledgeProfile{
			Name:        name,
			Description: description,
		}
		return normalized, joinSummary(
			"Brand profile",
			"Name: "+name,
			"Description: "+description,
		), nil
	case workspace.KnowledgeSectionEEAT:
		if document.EEAT == nil {
			return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
		}
		guidelines, err := normalizeKnowledgeList(document.EEAT.Guidelines, 20, 500)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		normalized.EEAT = &workspace.KnowledgeEEAT{Guidelines: guidelines}
		return normalized, joinSummary("E-E-A-T guidelines", guidelines...), nil
	case workspace.KnowledgeSectionCompetitors:
		if document.Competitors == nil {
			return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
		}
		urls, err := normalizeCompetitorURLs(document.Competitors.URLs)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		normalized.Competitors = &workspace.KnowledgeCompetitors{URLs: urls}
		return normalized, joinSummary("Competitor domains", urls...), nil
	case workspace.KnowledgeSectionTopics:
		if document.Topics == nil {
			return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
		}
		topics, err := normalizeKnowledgeList(document.Topics.Topics, 30, 200)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		normalized.Topics = &workspace.KnowledgeTopics{Topics: topics}
		return normalized, joinSummary("Priority topics", topics...), nil
	case workspace.KnowledgeSectionTone:
		if document.Tone == nil || !validPrimaryTone(document.Tone.PrimaryTone) {
			return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
		}
		instructions, err := validateKnowledgeText(document.Tone.CustomInstructions, 4_000)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		normalized.Tone = &workspace.KnowledgeTone{
			PrimaryTone:        document.Tone.PrimaryTone,
			CustomInstructions: instructions,
		}
		return normalized, joinSummary(
			"Brand tone",
			"Primary tone: "+document.Tone.PrimaryTone,
			"Custom instructions: "+instructions,
		), nil
	case workspace.KnowledgeSectionMemory:
		if document.Memory == nil {
			return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
		}
		facts, err := normalizeKnowledgeList(document.Memory.Facts, 50, 500)
		if err != nil {
			return workspace.KnowledgeDocument{}, "", err
		}
		normalized.Memory = &workspace.KnowledgeMemory{Facts: facts}
		return normalized, joinSummary("Approved brand facts", facts...), nil
	default:
		return workspace.KnowledgeDocument{}, "", ErrInvalidKnowledgeUpdate
	}
}

func payloadCount(document workspace.KnowledgeDocument) int {
	count := 0
	for _, present := range []bool{
		document.Profile != nil,
		document.EEAT != nil,
		document.Competitors != nil,
		document.Topics != nil,
		document.Tone != nil,
		document.Memory != nil,
	} {
		if present {
			count++
		}
	}
	return count
}

func validateKnowledgeText(value string, maxRunes int) (string, error) {
	if !utf8.ValidString(value) {
		return "", ErrInvalidKnowledgeUpdate
	}
	value = strings.TrimSpace(value)
	if utf8.RuneCountInString(value) > maxRunes {
		return "", ErrInvalidKnowledgeUpdate
	}
	for _, character := range value {
		if unicode.IsControl(character) &&
			character != '\n' &&
			character != '\r' &&
			character != '\t' {
			return "", ErrInvalidKnowledgeUpdate
		}
	}
	input := contracts.SanitizedInput{RawText: value}
	if err := middleware.CheckForPromptInjection(input); err != nil ||
		middleware.ContainsProtectedMetadata(value) {
		return "", ErrUnsafeKnowledgeUpdate
	}
	return value, nil
}

func normalizeKnowledgeList(values []string, maxItems, maxRunes int) ([]string, error) {
	if len(values) > maxItems {
		return nil, ErrInvalidKnowledgeUpdate
	}
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		clean, err := validateKnowledgeText(value, maxRunes)
		if err != nil {
			return nil, err
		}
		if clean == "" {
			return nil, ErrInvalidKnowledgeUpdate
		}
		key := strings.ToLower(clean)
		if _, exists := seen[key]; exists {
			return nil, ErrInvalidKnowledgeUpdate
		}
		seen[key] = struct{}{}
		normalized = append(normalized, clean)
	}
	return normalized, nil
}

func normalizeCompetitorURLs(values []string) ([]string, error) {
	if len(values) > 20 {
		return nil, ErrInvalidKnowledgeUpdate
	}
	normalized := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if len(value) > 2_048 {
			return nil, ErrInvalidKnowledgeUpdate
		}
		parsed, err := url.Parse(strings.TrimSpace(value))
		if err != nil ||
			(parsed.Scheme != "http" && parsed.Scheme != "https") ||
			parsed.Hostname() == "" ||
			parsed.User != nil {
			return nil, ErrInvalidKnowledgeUpdate
		}
		origin := parsed.Scheme + "://" + strings.ToLower(parsed.Host)
		if _, exists := seen[origin]; exists {
			return nil, ErrInvalidKnowledgeUpdate
		}
		seen[origin] = struct{}{}
		normalized = append(normalized, origin)
	}
	return normalized, nil
}

func validPrimaryTone(value string) bool {
	switch value {
	case "professional_authoritative",
		"conversational_friendly",
		"academic_technical",
		"persuasive_direct":
		return true
	default:
		return false
	}
}

func joinSummary(title string, values ...string) string {
	var builder strings.Builder
	builder.WriteString(title)
	for _, value := range values {
		if value == "" {
			continue
		}
		fmt.Fprintf(&builder, "\n- %s", value)
	}
	return builder.String()
}

func limitRunes(value string, limit int) string {
	if utf8.RuneCountInString(value) <= limit {
		return value
	}
	return string([]rune(value)[:limit])
}
