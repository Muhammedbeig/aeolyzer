package history

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"mime"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func (s *Store) SaveAttachment(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	input AttachmentInput,
) (AttachmentRef, error) {
	if err := validateIdentity(appName, userID, sessionID); err != nil {
		return AttachmentRef{}, err
	}
	if input.ID == "" {
		input.ID = uuid.NewString()
	}
	if !identityPattern.MatchString(input.ID) || input.Name == "" || len(input.Name) > 255 {
		return AttachmentRef{}, ErrInvalidReference
	}
	if len(input.Data) == 0 || int64(len(input.Data)) > s.config.MaxAttachmentBytes {
		return AttachmentRef{}, ErrInvalidReference
	}
	mediaType, _, err := mime.ParseMediaType(input.ContentType)
	if err != nil || mediaType == "" || len(mediaType) > 128 {
		return AttachmentRef{}, ErrInvalidReference
	}
	digest := sha256.Sum256(input.Data)
	if !bytes.Equal(digest[:], input.SHA256[:]) {
		return AttachmentRef{}, ErrInvalidReference
	}

	nameCiphertext, err := s.cipher.Encrypt(
		[]byte(input.Name),
		additionalData("attachment-name", appName, userID, sessionID, input.ID),
	)
	if err != nil {
		return AttachmentRef{}, fmt.Errorf("encrypt attachment name: %w", err)
	}
	dataCiphertext, err := s.cipher.Encrypt(
		input.Data,
		additionalData("attachment-data", appName, userID, sessionID, input.ID),
	)
	if err != nil {
		return AttachmentRef{}, fmt.Errorf("encrypt attachment data: %w", err)
	}
	result, err := s.db.ExecContext(
		ctx,
		`INSERT INTO aeolyzer_attachments
			(app_name, user_id, session_id, attachment_id, name_ciphertext, content_type, byte_size, sha256, data_ciphertext, created_at)
		 SELECT ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		appName,
		userID,
		sessionID,
		input.ID,
		nameCiphertext,
		mediaType,
		len(input.Data),
		digest[:],
		dataCiphertext,
		s.now().UTC().Truncate(time.Microsecond),
		appName,
		userID,
		sessionID,
	)
	if isDuplicateKey(err) {
		return AttachmentRef{}, ErrConflict
	}
	if err != nil {
		return AttachmentRef{}, fmt.Errorf("store attachment: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return AttachmentRef{}, fmt.Errorf("read attachment insert count: %w", err)
	}
	if affected != 1 {
		return AttachmentRef{}, ErrNotFound
	}
	return AttachmentRef{
		ID:          input.ID,
		Name:        input.Name,
		ContentType: mediaType,
		Size:        int64(len(input.Data)),
		SHA256:      digest,
	}, nil
}

func (s *Store) hydrateAttachments(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	events []*session.Event,
) error {
	remaining := s.config.MaxContextAttachmentBytes
	for eventIndex := len(events) - 1; eventIndex >= 0; eventIndex-- {
		event := events[eventIndex]
		if event.Content == nil {
			continue
		}
		for partIndex := len(event.Content.Parts) - 1; partIndex >= 0; partIndex-- {
			part := event.Content.Parts[partIndex]
			if part == nil {
				continue
			}
			attachmentID, ok := attachmentIDFromPart(part)
			part.PartMetadata = nil
			if !ok {
				continue
			}
			ref, data, err := s.loadAttachment(ctx, appName, userID, sessionID, attachmentID, true)
			if err != nil {
				return err
			}
			if ref.Size > remaining {
				event.Content.Parts[partIndex] = omittedAttachmentPart()
				continue
			}
			remaining -= ref.Size
			part.InlineData = &genai.Blob{
				Data:     data,
				MIMEType: ref.ContentType,
			}
		}
	}
	return nil
}

func (s *Store) loadAttachment(
	ctx context.Context,
	appName string,
	userID string,
	sessionID string,
	attachmentID string,
	withData bool,
) (AttachmentRef, []byte, error) {
	if err := validateIdentity(appName, userID, sessionID, attachmentID); err != nil {
		return AttachmentRef{}, nil, err
	}
	var nameCiphertext []byte
	var contentType string
	var size int64
	var storedHash []byte
	var dataCiphertext []byte
	query := `SELECT name_ciphertext, content_type, byte_size, sha256, data_ciphertext
	            FROM aeolyzer_attachments
	           WHERE app_name = ? AND user_id = ? AND session_id = ? AND attachment_id = ?`
	if !withData {
		query = `SELECT name_ciphertext, content_type, byte_size, sha256, NULL
		           FROM aeolyzer_attachments
		          WHERE app_name = ? AND user_id = ? AND session_id = ? AND attachment_id = ?`
	}
	err := s.db.QueryRowContext(ctx, query, appName, userID, sessionID, attachmentID).Scan(
		&nameCiphertext,
		&contentType,
		&size,
		&storedHash,
		&dataCiphertext,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return AttachmentRef{}, nil, ErrNotFound
	}
	if err != nil {
		return AttachmentRef{}, nil, fmt.Errorf("load attachment: %w", err)
	}
	if len(storedHash) != sha256.Size || size < 1 || size > s.config.MaxAttachmentBytes {
		return AttachmentRef{}, nil, ErrInvalidReference
	}
	name, err := s.cipher.Decrypt(
		nameCiphertext,
		additionalData("attachment-name", appName, userID, sessionID, attachmentID),
	)
	if err != nil {
		return AttachmentRef{}, nil, fmt.Errorf("decrypt attachment name: %w", err)
	}
	var digest [sha256.Size]byte
	copy(digest[:], storedHash)
	ref := AttachmentRef{
		ID:          attachmentID,
		Name:        string(name),
		ContentType: contentType,
		Size:        size,
		SHA256:      digest,
	}
	if !withData {
		return ref, nil, nil
	}
	data, err := s.cipher.Decrypt(
		dataCiphertext,
		additionalData("attachment-data", appName, userID, sessionID, attachmentID),
	)
	if err != nil {
		return AttachmentRef{}, nil, fmt.Errorf("decrypt attachment data: %w", err)
	}
	actualDigest := sha256.Sum256(data)
	if int64(len(data)) != size || !bytes.Equal(actualDigest[:], storedHash) {
		return AttachmentRef{}, nil, ErrInvalidReference
	}
	return ref, data, nil
}

func attachmentIDFromPart(part *genai.Part) (string, bool) {
	if part == nil || len(part.PartMetadata) != 1 {
		return "", false
	}
	value, ok := part.PartMetadata[attachmentMetadataKey]
	if !ok {
		return "", false
	}
	attachmentID, ok := value.(string)
	return attachmentID, ok && identityPattern.MatchString(attachmentID)
}
