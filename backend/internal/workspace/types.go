package workspace

import "time"

type ContentType string

const (
	ContentTypeArticle            ContentType = "article"
	ContentTypeBlogPost           ContentType = "blog_post"
	ContentTypeLinkedInPost       ContentType = "linkedin_post"
	ContentTypeYouTubeDescription ContentType = "youtube_description"
	ContentTypeProductDescription ContentType = "product_description"
)

func (t ContentType) Valid() bool {
	switch t {
	case ContentTypeArticle,
		ContentTypeBlogPost,
		ContentTypeLinkedInPost,
		ContentTypeYouTubeDescription,
		ContentTypeProductDescription:
		return true
	default:
		return false
	}
}

type KnowledgeSection string

const (
	KnowledgeSectionProfile     KnowledgeSection = "profile"
	KnowledgeSectionEEAT        KnowledgeSection = "eeat"
	KnowledgeSectionCompetitors KnowledgeSection = "competitors"
	KnowledgeSectionTopics      KnowledgeSection = "topics"
	KnowledgeSectionTone        KnowledgeSection = "tone"
	KnowledgeSectionMemory      KnowledgeSection = "memory"
)

func (s KnowledgeSection) Valid() bool {
	switch s {
	case KnowledgeSectionProfile,
		KnowledgeSectionEEAT,
		KnowledgeSectionCompetitors,
		KnowledgeSectionTopics,
		KnowledgeSectionTone,
		KnowledgeSectionMemory:
		return true
	default:
		return false
	}
}

type KnowledgeProfile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type KnowledgeEEAT struct {
	Guidelines []string `json:"guidelines"`
}

type KnowledgeCompetitors struct {
	URLs []string `json:"urls"`
}

type KnowledgeTopics struct {
	Topics []string `json:"topics"`
}

type KnowledgeTone struct {
	PrimaryTone        string `json:"primary_tone"`
	CustomInstructions string `json:"custom_instructions"`
}

type KnowledgeMemory struct {
	Facts []string `json:"facts"`
}

// KnowledgeDocument contains one section payload and its optimistic-lock version.
type KnowledgeDocument struct {
	Section     KnowledgeSection      `json:"section"`
	Version     uint64                `json:"version"`
	Profile     *KnowledgeProfile     `json:"profile,omitempty"`
	EEAT        *KnowledgeEEAT        `json:"eeat,omitempty"`
	Competitors *KnowledgeCompetitors `json:"competitors,omitempty"`
	Topics      *KnowledgeTopics      `json:"topics,omitempty"`
	Tone        *KnowledgeTone        `json:"tone,omitempty"`
	Memory      *KnowledgeMemory      `json:"memory,omitempty"`
	UpdatedAt   *time.Time            `json:"updated_at,omitempty"`
}

func EmptyKnowledgeDocument(section KnowledgeSection) KnowledgeDocument {
	document := KnowledgeDocument{Section: section}
	switch section {
	case KnowledgeSectionProfile:
		document.Profile = &KnowledgeProfile{}
	case KnowledgeSectionEEAT:
		document.EEAT = &KnowledgeEEAT{Guidelines: []string{}}
	case KnowledgeSectionCompetitors:
		document.Competitors = &KnowledgeCompetitors{URLs: []string{}}
	case KnowledgeSectionTopics:
		document.Topics = &KnowledgeTopics{Topics: []string{}}
	case KnowledgeSectionTone:
		document.Tone = &KnowledgeTone{PrimaryTone: "professional_authoritative"}
	case KnowledgeSectionMemory:
		document.Memory = &KnowledgeMemory{Facts: []string{}}
	}
	return document
}
