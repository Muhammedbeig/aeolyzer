package history

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (s *Store) ClaimMessageRequest(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	idempotencyKey string,
) (MessageRequestClaim, error) {
	if err := validateIdentity(appName, userID, sessionID); err != nil {
		return MessageRequestClaim{}, err
	}
	if len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
		return MessageRequestClaim{}, errors.New("invalid idempotency key")
	}
	requestHash := sha256.Sum256([]byte(idempotencyKey))
	now := s.now().UTC().Truncate(time.Microsecond)
	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO aeolyzer_message_requests
			(app_name, user_id, session_id, request_hash, status, created_at, updated_at)
		 SELECT ?, ?, ?, ?, 'pending', ?, ?
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		appName,
		userID,
		sessionID,
		requestHash[:],
		now,
		now,
		appName,
		userID,
		sessionID,
	)
	if err == nil {
		affected, rowsErr := result.RowsAffected()
		if rowsErr != nil {
			return MessageRequestClaim{}, fmt.Errorf("read message request insert count: %w", rowsErr)
		}
		if affected != 1 {
			return MessageRequestClaim{}, ErrNotFound
		}
		return MessageRequestClaim{Claimed: true}, nil
	}
	if !isDuplicateKey(err) {
		return MessageRequestClaim{}, fmt.Errorf("claim message request: %w", err)
	}

	var status string
	var responseCiphertext []byte
	err = s.db.QueryRowContext(
		ctx,
		`SELECT status, response_ciphertext
		   FROM aeolyzer_message_requests
		  WHERE app_name = ? AND user_id = ? AND session_id = ? AND request_hash = ?`,
		appName,
		userID,
		sessionID,
		requestHash[:],
	).Scan(&status, &responseCiphertext)
	if errors.Is(err, sql.ErrNoRows) {
		return MessageRequestClaim{}, ErrNotFound
	}
	if err != nil {
		return MessageRequestClaim{}, fmt.Errorf("read message request: %w", err)
	}
	switch status {
	case "completed":
		response, err := s.cipher.Decrypt(
			responseCiphertext,
			additionalData("request", appName, userID, sessionID, fmt.Sprintf("%x", requestHash[:])),
		)
		if err != nil {
			return MessageRequestClaim{}, fmt.Errorf("decrypt message response: %w", err)
		}
		return MessageRequestClaim{CachedResponse: response}, nil
	case "pending":
		return MessageRequestClaim{}, ErrRequestPending
	case "failed":
		return MessageRequestClaim{}, ErrRequestFailed
	default:
		return MessageRequestClaim{}, errors.New("invalid message request state")
	}
}

func (s *Store) CompleteMessageRequest(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	idempotencyKey string,
	response []byte,
) error {
	return s.finishMessageRequest(ctx, appName, userID, sessionID, idempotencyKey, "completed", response)
}

func (s *Store) FailMessageRequest(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	idempotencyKey string,
) error {
	return s.finishMessageRequest(ctx, appName, userID, sessionID, idempotencyKey, "failed", nil)
}

func (s *Store) finishMessageRequest(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	idempotencyKey string,
	status string,
	response []byte,
) error {
	if err := validateIdentity(appName, userID, sessionID); err != nil {
		return err
	}
	if len(idempotencyKey) < 16 || len(idempotencyKey) > 128 {
		return errors.New("invalid idempotency key")
	}
	requestHash := sha256.Sum256([]byte(idempotencyKey))
	var responseCiphertext []byte
	var err error
	if len(response) > 0 {
		responseCiphertext, err = s.cipher.Encrypt(
			response,
			additionalData("request", appName, userID, sessionID, fmt.Sprintf("%x", requestHash[:])),
		)
		if err != nil {
			return fmt.Errorf("encrypt message response: %w", err)
		}
	}
	result, err := s.db.ExecContext(
		ctx,
		`UPDATE aeolyzer_message_requests
		    SET status = ?, response_ciphertext = ?, updated_at = ?
		  WHERE app_name = ? AND user_id = ? AND session_id = ? AND request_hash = ? AND status = 'pending'`,
		status,
		responseCiphertext,
		s.now().UTC().Truncate(time.Microsecond),
		appName,
		userID,
		sessionID,
		requestHash[:],
	)
	if err != nil {
		return fmt.Errorf("finish message request: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read message request update count: %w", err)
	}
	if affected != 1 {
		return ErrConflict
	}
	return nil
}
