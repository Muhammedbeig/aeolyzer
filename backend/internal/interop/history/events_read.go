package history

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/adk/session"
)

func (s *Store) readEvents(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	limit int,
	after time.Time,
) ([]*session.Event, error) {
	query := `SELECT event_id, event_ciphertext, created_at
	            FROM (
	                  SELECT event_id, event_ciphertext, created_at, sequence_number
	                    FROM aeolyzer_events
	                   WHERE app_name = ? AND user_id = ? AND session_id = ?
	                   ORDER BY sequence_number DESC
	                   LIMIT ?
	                 ) AS recent
	           ORDER BY sequence_number ASC`
	arguments := []any{appName, userID, sessionID, limit}
	if !after.IsZero() {
		query = `SELECT event_id, event_ciphertext, created_at
		           FROM (
		                 SELECT event_id, event_ciphertext, created_at, sequence_number
		                   FROM aeolyzer_events
		                  WHERE app_name = ? AND user_id = ? AND session_id = ? AND created_at >= ?
		                  ORDER BY sequence_number DESC
		                  LIMIT ?
		                ) AS recent
		          ORDER BY sequence_number ASC`
		arguments = []any{appName, userID, sessionID, after.UTC(), limit}
	}
	rows, err := s.db.QueryContext(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("read conversation events: %w", err)
	}
	defer rows.Close()

	events := make([]*session.Event, 0, limit)
	for rows.Next() {
		var eventID string
		var ciphertext []byte
		var createdAt time.Time
		if err := rows.Scan(&eventID, &ciphertext, &createdAt); err != nil {
			return nil, fmt.Errorf("scan conversation event: %w", err)
		}
		plaintext, err := s.cipher.Decrypt(
			ciphertext,
			additionalData("event", appName, userID, sessionID, eventID),
		)
		if err != nil {
			return nil, fmt.Errorf("decrypt conversation event: %w", err)
		}
		event, err := unmarshalEvent(plaintext)
		if err != nil {
			return nil, err
		}
		if event.ID != eventID || !event.Timestamp.Equal(createdAt.UTC()) {
			return nil, errors.New("conversation event integrity check failed")
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversation events: %w", err)
	}
	return events, nil
}

func (s *Store) ListMessages(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
) ([]StoredMessage, error) {
	if err := validateIdentity(appName, userID, sessionID); err != nil {
		return nil, err
	}
	events, err := s.readEvents(ctx, appName, userID, sessionID, s.config.MaxUIEvents, time.Time{})
	if err != nil {
		return nil, err
	}
	messages := make([]StoredMessage, 0, len(events))
	for _, event := range events {
		if event.Content == nil || event.Partial {
			continue
		}
		message := StoredMessage{
			ID:        event.ID,
			Author:    event.Author,
			CreatedAt: event.Timestamp.UTC(),
		}
		var text strings.Builder
		for _, part := range event.Content.Parts {
			if part == nil || part.Thought {
				continue
			}
			if part.Text != "" {
				if text.Len() > 0 {
					text.WriteByte('\n')
				}
				text.WriteString(part.Text)
			}
			attachmentID, ok := attachmentIDFromPart(part)
			if !ok {
				continue
			}
			ref, _, err := s.loadAttachment(ctx, appName, userID, sessionID, attachmentID, false)
			if err != nil {
				return nil, err
			}
			message.Attachments = append(message.Attachments, ref)
		}
		message.Text = text.String()
		if message.Text != "" || len(message.Attachments) > 0 {
			messages = append(messages, message)
		}
	}
	return messages, nil
}
