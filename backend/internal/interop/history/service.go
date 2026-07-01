package history

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	mysql "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"google.golang.org/adk/session"
)

func (s *Store) Create(ctx context.Context, request *session.CreateRequest) (*session.CreateResponse, error) {
	if request == nil {
		return nil, errors.New("create session request is nil")
	}
	sessionID := request.SessionID
	if sessionID == "" {
		sessionID = uuid.NewString()
	}
	if err := validateIdentity(request.AppName, request.UserID, sessionID); err != nil {
		return nil, err
	}
	stateDelta, err := persistentStateDelta(request.State)
	if err != nil {
		return nil, err
	}
	stateBytes, err := marshalState(stateDelta)
	if err != nil {
		return nil, err
	}
	stateCiphertext, err := s.cipher.Encrypt(
		stateBytes,
		additionalData("state", request.AppName, request.UserID, sessionID, sessionID),
	)
	if err != nil {
		return nil, fmt.Errorf("encrypt session state: %w", err)
	}
	now := s.now().UTC().Truncate(time.Microsecond)
	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO aeolyzer_sessions
			(app_name, user_id, session_id, state_ciphertext, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		request.AppName,
		request.UserID,
		sessionID,
		stateCiphertext,
		now,
		now,
	)
	if isDuplicateKey(err) {
		return nil, ErrConflict
	}
	if err != nil {
		return nil, fmt.Errorf("create conversation: %w", err)
	}
	return &session.CreateResponse{
		Session: &localSession{
			appName:   request.AppName,
			userID:    request.UserID,
			id:        sessionID,
			state:     stateDelta,
			updatedAt: now,
		},
	}, nil
}

func (s *Store) Get(ctx context.Context, request *session.GetRequest) (*session.GetResponse, error) {
	if request == nil {
		return nil, errors.New("get session request is nil")
	}
	if err := validateIdentity(request.AppName, request.UserID, request.SessionID); err != nil {
		return nil, err
	}
	var stateCiphertext []byte
	var updatedAt time.Time
	err := s.db.QueryRowContext(
		ctx,
		`SELECT state_ciphertext, updated_at
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		request.AppName,
		request.UserID,
		request.SessionID,
	).Scan(&stateCiphertext, &updatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get conversation: %w", err)
	}
	state, err := s.decryptState(request.AppName, request.UserID, request.SessionID, stateCiphertext)
	if err != nil {
		return nil, err
	}

	limit := request.NumRecentEvents
	if limit < 1 || limit > s.config.MaxModelEvents {
		limit = s.config.MaxModelEvents
	}
	events, err := s.readEvents(ctx, request.AppName, request.UserID, request.SessionID, limit, request.After)
	if err != nil {
		return nil, err
	}
	if err := s.hydrateAttachments(ctx, request.AppName, request.UserID, request.SessionID, events); err != nil {
		return nil, err
	}
	return &session.GetResponse{
		Session: &localSession{
			appName:   request.AppName,
			userID:    request.UserID,
			id:        request.SessionID,
			state:     state,
			events:    events,
			updatedAt: updatedAt.UTC(),
		},
	}, nil
}

func (s *Store) List(ctx context.Context, request *session.ListRequest) (*session.ListResponse, error) {
	if request == nil {
		return nil, errors.New("list sessions request is nil")
	}
	if err := validateIdentity(request.AppName, request.UserID); err != nil {
		return nil, err
	}
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT session_id, state_ciphertext, updated_at
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ?
		  ORDER BY updated_at DESC
		  LIMIT 100`,
		request.AppName,
		request.UserID,
	)
	if err != nil {
		return nil, fmt.Errorf("list conversations: %w", err)
	}
	defer rows.Close()

	sessions := make([]session.Session, 0)
	for rows.Next() {
		var sessionID string
		var stateCiphertext []byte
		var updatedAt time.Time
		if err := rows.Scan(&sessionID, &stateCiphertext, &updatedAt); err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}
		state, err := s.decryptState(request.AppName, request.UserID, sessionID, stateCiphertext)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &localSession{
			appName:   request.AppName,
			userID:    request.UserID,
			id:        sessionID,
			state:     state,
			updatedAt: updatedAt.UTC(),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversations: %w", err)
	}
	return &session.ListResponse{Sessions: sessions}, nil
}

func (s *Store) Delete(ctx context.Context, request *session.DeleteRequest) error {
	if request == nil {
		return errors.New("delete session request is nil")
	}
	if err := validateIdentity(request.AppName, request.UserID, request.SessionID); err != nil {
		return err
	}
	result, err := s.db.ExecContext(
		ctx,
		`DELETE FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		request.AppName,
		request.UserID,
		request.SessionID,
	)
	if err != nil {
		return fmt.Errorf("delete conversation: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read deleted conversation count: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AppendEvent(ctx context.Context, current session.Session, event *session.Event) error {
	if current == nil {
		return errors.New("session is nil")
	}
	if event == nil {
		return errors.New("event is nil")
	}
	if event.Partial {
		return nil
	}
	local, ok := current.(*localSession)
	if !ok {
		return fmt.Errorf("unexpected session type %T", current)
	}
	if err := validateIdentity(local.appName, local.userID, local.id); err != nil {
		return err
	}
	if event.ID == "" {
		event.ID = uuid.NewString()
	}
	if !identityPattern.MatchString(event.ID) {
		return errors.New("invalid event identity")
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = s.now()
	}
	event.Timestamp = event.Timestamp.UTC().Truncate(time.Microsecond)

	var refs []AttachmentRef
	if event.Author == "user" {
		refs = attachmentRefsFromContext(ctx)
	}
	memoryEvent, storedEvent, err := prepareEventsForAppend(event, refs)
	if err != nil {
		return err
	}
	eventBytes, err := marshalEvent(storedEvent)
	if err != nil {
		return err
	}
	eventCiphertext, err := s.cipher.Encrypt(
		eventBytes,
		additionalData("event", local.appName, local.userID, local.id, event.ID),
	)
	if err != nil {
		return fmt.Errorf("encrypt session event: %w", err)
	}

	transaction, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin event transaction: %w", err)
	}
	defer transaction.Rollback()

	var sequence uint64
	var stateCiphertext []byte
	err = transaction.QueryRowContext(
		ctx,
		`SELECT next_sequence, state_ciphertext
		   FROM aeolyzer_sessions
		  WHERE app_name = ? AND user_id = ? AND session_id = ?
		  FOR UPDATE`,
		local.appName,
		local.userID,
		local.id,
	).Scan(&sequence, &stateCiphertext)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("lock conversation: %w", err)
	}

	state, err := s.decryptState(local.appName, local.userID, local.id, stateCiphertext)
	if err != nil {
		return err
	}
	state = mergeState(state, storedEvent.Actions.StateDelta)
	stateBytes, err := marshalState(state)
	if err != nil {
		return err
	}
	stateCiphertext, err = s.cipher.Encrypt(
		stateBytes,
		additionalData("state", local.appName, local.userID, local.id, local.id),
	)
	if err != nil {
		return fmt.Errorf("encrypt session state: %w", err)
	}

	_, err = transaction.ExecContext(
		ctx,
		`INSERT INTO aeolyzer_events
			(app_name, user_id, session_id, event_id, sequence_number, event_ciphertext, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		local.appName,
		local.userID,
		local.id,
		event.ID,
		sequence,
		eventCiphertext,
		event.Timestamp,
	)
	if isDuplicateKey(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("append conversation event: %w", err)
	}
	_, err = transaction.ExecContext(
		ctx,
		`UPDATE aeolyzer_sessions
		    SET next_sequence = ?, state_ciphertext = ?, updated_at = ?
		  WHERE app_name = ? AND user_id = ? AND session_id = ?`,
		sequence+1,
		stateCiphertext,
		event.Timestamp,
		local.appName,
		local.userID,
		local.id,
	)
	if err != nil {
		return fmt.Errorf("update conversation after event: %w", err)
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("commit conversation event: %w", err)
	}
	local.appendEvent(memoryEvent)
	return nil
}

func (s *Store) decryptState(appName, userID, sessionID string, ciphertext []byte) (map[string]any, error) {
	plaintext, err := s.cipher.Decrypt(
		ciphertext,
		additionalData("state", appName, userID, sessionID, sessionID),
	)
	if err != nil {
		return nil, fmt.Errorf("decrypt session state: %w", err)
	}
	return unmarshalState(plaintext)
}

func isDuplicateKey(err error) bool {
	var mysqlError *mysql.MySQLError
	return errors.As(err, &mysqlError) && mysqlError.Number == 1062
}
