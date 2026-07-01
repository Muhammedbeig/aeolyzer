package extensions

import "aeolyzer/internal/workspace"

type KnowledgeSection = workspace.KnowledgeSection

const (
	KnowledgeSectionProfile     = workspace.KnowledgeSectionProfile
	KnowledgeSectionEEAT        = workspace.KnowledgeSectionEEAT
	KnowledgeSectionCompetitors = workspace.KnowledgeSectionCompetitors
	KnowledgeSectionTopics      = workspace.KnowledgeSectionTopics
	KnowledgeSectionTone        = workspace.KnowledgeSectionTone
	KnowledgeSectionMemory      = workspace.KnowledgeSectionMemory
)

type KnowledgeProfile = workspace.KnowledgeProfile
type KnowledgeEEAT = workspace.KnowledgeEEAT
type KnowledgeCompetitors = workspace.KnowledgeCompetitors
type KnowledgeTopics = workspace.KnowledgeTopics
type KnowledgeTone = workspace.KnowledgeTone
type KnowledgeMemory = workspace.KnowledgeMemory
type KnowledgeDocument = workspace.KnowledgeDocument

func EmptyKnowledgeDocument(section KnowledgeSection) KnowledgeDocument {
	return workspace.EmptyKnowledgeDocument(section)
}
