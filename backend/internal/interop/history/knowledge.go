package history

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	maxKnowledgeBodyBytes    = 64 << 10
	maxKnowledgeSummaryBytes = 16 << 10
	maxAgentContextRunes     = 16_000
)

func (s *Store) GetKnowledge(
	ctx context.Context,
	userID string,
	section string,
) (KnowledgeRecord, error) {
	if err := validateIdentity(userID, section); err != nil {
		return KnowledgeRecord{}, err
	}
	var version uint64
	var bodyCiphertext []byte
	var updatedAt time.Time
	err := s.db.QueryRowContext(
		ctx,
		`SELECT version, body_ciphertext, updated_at
		   FROM aeolyzer_knowledge
		  WHERE user_id = ? AND section = ?`,
		userID,
		section,
	).Scan(&version, &bodyCiphertext, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return KnowledgeRecord{Section: section}, nil
	}
	if err != nil {
		return KnowledgeRecord{}, fmt.Errorf("get knowledge section: %w", err)
	}
	body, err := s.cipher.Decrypt(
		bodyCiphertext,
		knowledgeAdditionalData("knowledge-body", userID, section, version),
	)
	if err != nil {
		return KnowledgeRecord{}, fmt.Errorf("decrypt knowledge section: %w", err)
	}
	if len(body) > maxKnowledgeBodyBytes {
		return KnowledgeRecord{}, errors.New("knowledge section exceeds storage limit")
	}
	return KnowledgeRecord{
		Section:   section,
		Version:   version,
		Body:      body,
		UpdatedAt: updatedAt.UTC(),
	}, nil
}

func (s *Store) UpdateKnowledge(
	ctx context.Context,
	userID string,
	section string,
	expectedVersion uint64,
	body []byte,
	summary string,
) (KnowledgeRecord, error) {
	if err := validateIdentity(userID, section); err != nil {
		return KnowledgeRecord{}, err
	}
	if len(body) == 0 || len(body) > maxKnowledgeBodyBytes ||
		len(summary) > maxKnowledgeSummaryBytes ||
		!utf8.ValidString(summary) {
		return KnowledgeRecord{}, errors.New("invalid knowledge payload")
	}
	nextVersion := expectedVersion + 1
	bodyCiphertext, err := s.cipher.Encrypt(
		body,
		knowledgeAdditionalData("knowledge-body", userID, section, nextVersion),
	)
	if err != nil {
		return KnowledgeRecord{}, fmt.Errorf("encrypt knowledge section: %w", err)
	}
	summaryCiphertext, err := s.cipher.Encrypt(
		[]byte(summary),
		knowledgeAdditionalData("knowledge-summary", userID, section, nextVersion),
	)
	if err != nil {
		return KnowledgeRecord{}, fmt.Errorf("encrypt knowledge summary: %w", err)
	}
	now := s.now().UTC().Truncate(time.Microsecond)
	if expectedVersion == 0 {
		_, err = s.db.ExecContext(
			ctx,
			`INSERT INTO aeolyzer_knowledge
				(user_id, section, version, body_ciphertext, summary_ciphertext, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			userID,
			section,
			nextVersion,
			bodyCiphertext,
			summaryCiphertext,
			now,
			now,
		)
		if isDuplicateKey(err) {
			return KnowledgeRecord{}, ErrConflict
		}
		if err != nil {
			return KnowledgeRecord{}, fmt.Errorf("create knowledge section: %w", err)
		}
	} else {
		result, err := s.db.ExecContext(
			ctx,
			`UPDATE aeolyzer_knowledge
			    SET version = ?, body_ciphertext = ?, summary_ciphertext = ?, updated_at = ?
			  WHERE user_id = ? AND section = ? AND version = ?`,
			nextVersion,
			bodyCiphertext,
			summaryCiphertext,
			now,
			userID,
			section,
			expectedVersion,
		)
		if err != nil {
			return KnowledgeRecord{}, fmt.Errorf("update knowledge section: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return KnowledgeRecord{}, fmt.Errorf("read knowledge update count: %w", err)
		}
		if affected == 0 {
			return KnowledgeRecord{}, ErrConflict
		}
	}
	return KnowledgeRecord{
		Section:   section,
		Version:   nextVersion,
		Body:      append([]byte(nil), body...),
		UpdatedAt: now,
	}, nil
}

func (s *Store) AgentKnowledgeContext(
	ctx context.Context,
	userID string,
) (string, error) {
	if err := validateIdentity(userID); err != nil {
		return "", err
	}
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT section, version, summary_ciphertext
		   FROM aeolyzer_knowledge
		  WHERE user_id = ?
		  ORDER BY section`,
		userID,
	)
	if err != nil {
		return "", fmt.Errorf("list knowledge summaries: %w", err)
	}
	defer rows.Close()

	var contextBuilder strings.Builder
	for rows.Next() {
		var section string
		var version uint64
		var summaryCiphertext []byte
		if err := rows.Scan(&section, &version, &summaryCiphertext); err != nil {
			return "", fmt.Errorf("scan knowledge summary: %w", err)
		}
		summary, err := s.cipher.Decrypt(
			summaryCiphertext,
			knowledgeAdditionalData("knowledge-summary", userID, section, version),
		)
		if err != nil {
			return "", fmt.Errorf("decrypt knowledge summary: %w", err)
		}
		if len(summary) > maxKnowledgeSummaryBytes || !utf8.Valid(summary) {
			return "", errors.New("knowledge summary integrity check failed")
		}
		if contextBuilder.Len() > 0 {
			contextBuilder.WriteString("\n\n")
		}
		contextBuilder.Write(summary)
	}
	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("iterate knowledge summaries: %w", err)
	}
	return limitContextRunes(contextBuilder.String(), maxAgentContextRunes), nil
}

func knowledgeAdditionalData(kind, userID, section string, version uint64) []byte {
	return additionalData(
		kind,
		"aeolyzer-knowledge",
		userID,
		section,
		strconv.FormatUint(version, 10),
	)
}

func limitContextRunes(value string, limit int) string {
	if utf8.RuneCountInString(value) <= limit {
		return value
	}
	return string([]rune(value)[:limit])
}
