package history

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

var identityPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

type Store struct {
	db     *sql.DB
	cipher *Cipher
	config Config
	now    func() time.Time
}

func Open(ctx context.Context, dsn string, key []byte, config Config) (*Store, error) {
	parsed, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse database dsn: %w", err)
	}
	if err := validateTransport(parsed); err != nil {
		return nil, err
	}
	parsed.ParseTime = true
	parsed.Loc = time.UTC
	if parsed.Timeout == 0 {
		parsed.Timeout = 5 * time.Second
	}
	if parsed.ReadTimeout == 0 {
		parsed.ReadTimeout = 10 * time.Second
	}
	if parsed.WriteTimeout == 0 {
		parsed.WriteTimeout = 10 * time.Second
	}
	db, err := sql.Open("mysql", parsed.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(30 * time.Minute)

	cipher, err := NewCipher(key)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	store, err := NewStore(db, cipher, config)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	if err := Migrate(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func NewStore(db *sql.DB, cipher *Cipher, config Config) (*Store, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}
	if cipher == nil {
		return nil, errors.New("cipher is nil")
	}
	defaults := DefaultConfig()
	if config.MaxModelEvents < 1 {
		config.MaxModelEvents = defaults.MaxModelEvents
	}
	if config.MaxUIEvents < 1 {
		config.MaxUIEvents = defaults.MaxUIEvents
	}
	if config.MaxContextAttachmentBytes < 1 {
		config.MaxContextAttachmentBytes = defaults.MaxContextAttachmentBytes
	}
	if config.MaxAttachmentBytes < 1 {
		config.MaxAttachmentBytes = defaults.MaxAttachmentBytes
	}
	if config.Retention <= 0 {
		config.Retention = defaults.Retention
	}
	return &Store{
		db:     db,
		cipher: cipher,
		config: config,
		now:    time.Now,
	}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) PurgeExpired(ctx context.Context) (int64, error) {
	if s == nil || s.db == nil {
		return 0, errors.New("history store is nil")
	}
	result, err := s.db.ExecContext(
		ctx,
		`DELETE FROM aeolyzer_sessions WHERE updated_at < ?`,
		s.now().UTC().Add(-s.config.Retention),
	)
	if err != nil {
		return 0, fmt.Errorf("purge expired conversations: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("read purge count: %w", err)
	}
	return count, nil
}

func validateTransport(config *mysql.Config) error {
	if config == nil {
		return errors.New("database config is nil")
	}
	if config.Net != "tcp" || config.TLSConfig != "" {
		return nil
	}
	host, _, err := net.SplitHostPort(config.Addr)
	if err != nil {
		host = config.Addr
	}
	host = strings.Trim(host, "[]")
	if strings.EqualFold(host, "localhost") {
		return nil
	}
	ip := net.ParseIP(host)
	if ip != nil && ip.IsLoopback() {
		return nil
	}
	return errors.New("remote database connections require tls")
}

func validateIdentity(values ...string) error {
	for _, value := range values {
		if !identityPattern.MatchString(value) {
			return errors.New("invalid conversation identity")
		}
	}
	return nil
}

func additionalData(kind, appName, userID, sessionID, objectID string) []byte {
	return []byte(strings.Join([]string{kind, appName, userID, sessionID, objectID}, "\x00"))
}
