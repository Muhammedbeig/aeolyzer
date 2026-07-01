package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"aeolyzer/internal/extensions"
	"aeolyzer/internal/intake"
)

func (s *ChatService) GetKnowledge(
	ctx context.Context,
	userID string,
	section extensions.KnowledgeSection,
) (extensions.KnowledgeDocument, error) {
	if !section.Valid() {
		return extensions.KnowledgeDocument{}, intake.ErrInvalidKnowledgeUpdate
	}
	record, err := s.knowledge.GetKnowledge(ctx, userID, string(section))
	if err != nil {
		return extensions.KnowledgeDocument{}, err
	}
	if record.Version == 0 {
		return extensions.EmptyKnowledgeDocument(section), nil
	}
	var document extensions.KnowledgeDocument
	if err := json.Unmarshal(record.Body, &document); err != nil {
		return extensions.KnowledgeDocument{}, fmt.Errorf("decode knowledge section: %w", err)
	}
	if document.Section != section {
		return extensions.KnowledgeDocument{}, errors.New("knowledge section integrity check failed")
	}
	document.Version = record.Version
	document.UpdatedAt = &record.UpdatedAt
	return document, nil
}

func (s *ChatService) UpdateKnowledge(
	ctx context.Context,
	userID string,
	update intake.ValidatedKnowledgeUpdate,
) (extensions.KnowledgeDocument, error) {
	document := update.Document
	expectedVersion := document.Version
	document.Version = 0
	document.UpdatedAt = nil
	body, err := json.Marshal(document)
	if err != nil {
		return extensions.KnowledgeDocument{}, fmt.Errorf("encode knowledge section: %w", err)
	}
	record, err := s.knowledge.UpdateKnowledge(
		ctx,
		userID,
		string(document.Section),
		expectedVersion,
		body,
		update.Summary,
	)
	if err != nil {
		return extensions.KnowledgeDocument{}, err
	}
	document.Version = record.Version
	document.UpdatedAt = &record.UpdatedAt
	return document, nil
}
