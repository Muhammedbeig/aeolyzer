package history

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"google.golang.org/adk/session"
)

func (s *Store) CreateConversation(
	ctx context.Context,
	appName string,
	userID string,
) (Conversation, error) {
	response, err := s.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		return Conversation{}, err
	}
	now := response.Session.LastUpdateTime()
	return Conversation{
		AppName:   appName,
		UserID:    userID,
		ID:        response.Session.ID(),
		Title:     "New chat",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *Store) GetConversation(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
) (Conversation, error) {
	if err := validateIdentity(appName, userID, sessionID); err != nil {
		return Conversation{}, err
	}
	var titleCiphertext []byte
	var starred bool
	var createdAt time.Time
	var updatedAt time.Time
	err := s.db.QueryRowContext(
		ctx,
		`SELECT title_ciphertext, starred, created_at, updated_at
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		appName,
		userID,
		sessionID,
	).Scan(&titleCiphertext, &starred, &createdAt, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Conversation{}, ErrNotFound
	}
	if err != nil {
		return Conversation{}, fmt.Errorf("get conversation metadata: %w", err)
	}
	title, err := s.decryptTitle(appName, userID, sessionID, titleCiphertext)
	if err != nil {
		return Conversation{}, err
	}
	return Conversation{
		AppName:   appName,
		UserID:    userID,
		ID:        sessionID,
		Title:     title,
		Starred:   starred,
		CreatedAt: createdAt.UTC(),
		UpdatedAt: updatedAt.UTC(),
	}, nil
}

func (s *Store) ListConversations(
	ctx context.Context,
	appName string,
	userID string,
) ([]Conversation, error) {
	if err := validateIdentity(appName, userID); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT session_id, title_ciphertext, starred, created_at, updated_at
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ?
		  ORDER BY starred DESC, updated_at DESC
		  LIMIT 100`,
		appName,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list conversation metadata: %w", err)
	}
	defer rows.Close()

	conversations := make([]Conversation, 0)
	for rows.Next() {
		var conversation Conversation
		var titleCiphertext []byte
		conversation.AppName = appName
		conversation.UserID = userID
		if err := rows.Scan(
			&conversation.ID,
			&titleCiphertext,
			&conversation.Starred,
			&conversation.CreatedAt,
			&conversation.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan conversation metadata: %w", err)
		}
		conversation.Title, err = s.decryptTitle(appName, userID, conversation.ID, titleCiphertext)
		if err != nil {
			return nil, err
		}
		conversation.CreatedAt = conversation.CreatedAt.UTC()
		conversation.UpdatedAt = conversation.UpdatedAt.UTC()
		conversations = append(conversations, conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversation metadata: %w", err)
	}
	return conversations, nil
}

func (s *Store) UpdateConversation(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	title *string,
	starred *bool,
) (Conversation, error) {
	if err := validateIdentity(appName, userID, sessionID); err != nil {
		return Conversation{}, err
	}
	if title == nil && starred == nil {
		return Conversation{}, errors.New("conversation update is empty")
	}
	sets := make([]string, 0, 3)
	arguments := make([]any, 0, 7)
	if title != nil {
		cleanTitle := normalizeTitle(*title)
		titleCiphertext, err := s.cipher.Encrypt(
			[]byte(cleanTitle),
			additionalData("title", appName, userID, sessionID, sessionID),
		)
		if err != nil {
			return Conversation{}, fmt.Errorf("encrypt conversation title: %w", err)
		}
		sets = append(sets, "title_ciphertext = ?")
		arguments = append(arguments, titleCiphertext)
	}
	if starred != nil {
		sets = append(sets, "starred = ?")
		arguments = append(arguments, *starred)
	}
	sets = append(sets, "updated_at = ?")
	arguments = append(arguments, s.now().UTC().Truncate(time.Microsecond))
	arguments = append(arguments, appName, userID, sessionID)
	result, err := s.db.ExecContext(
		ctx,
		`UPDATE aeolyzer_sessions SET `+strings.Join(sets, ", ")+`
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		arguments...,
	)
	if err != nil {
		return Conversation{}, fmt.Errorf("update conversation: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Conversation{}, fmt.Errorf("read conversation update count: %w", err)
	}
	if affected == 0 {
		return Conversation{}, ErrNotFound
	}
	return s.GetConversation(ctx, appName, userID, sessionID)
}

func (s *Store) decryptTitle(appName, userID, sessionID string, ciphertext []byte) (string, error) {
	if len(ciphertext) == 0 {
		return "New chat", nil
	}
	plaintext, err := s.cipher.Decrypt(
		ciphertext,
		additionalData("title", appName, userID, sessionID, sessionID),
	)
	if err != nil {
		return "", fmt.Errorf("decrypt conversation title: %w", err)
	}
	return normalizeTitle(string(plaintext)), nil
}

func normalizeTitle(title string) string {
	title = strings.Join(strings.Fields(title), " ")
	runes := []rune(title)
	if len(runes) > 80 {
		title = string(runes[:80])
	}
	if title == "" {
		return "New chat"
	}
	return title
}
