package history

import (
	"crypto/sha256"
	"errors"
	"time"
)

var (
	ErrNotFound         = errors.New("conversation not found")
	ErrConflict         = errors.New("conversation conflict")
	ErrRequestPending   = errors.New("message request is already in progress")
	ErrRequestFailed    = errors.New("message request previously failed")
	ErrInvalidReference = errors.New("invalid attachment reference")
)

type Config struct {
	MaxModelEvents            int
	MaxUIEvents               int
	MaxContextAttachmentBytes int64
	MaxAttachmentBytes        int64
	Retention                 time.Duration
}

func DefaultConfig() Config {
	return Config{
		MaxModelEvents:            40,
		MaxUIEvents:               100,
		MaxContextAttachmentBytes: 30 << 20,
		MaxAttachmentBytes:        10 << 20,
		Retention:                 30 * 24 * time.Hour,
	}
}

type Conversation struct {
	AppName   string
	UserID    string
	ID        string
	Title     string
	Starred   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AttachmentInput struct {
	ID          string
	Name        string
	ContentType string
	Data        []byte
	SHA256      [sha256.Size]byte
}

type AttachmentRef struct {
	ID          string
	Name        string
	ContentType string
	Size        int64
	SHA256      [sha256.Size]byte
}

type StoredMessage struct {
	ID          string
	Author      string
	Text        string
	Attachments []AttachmentRef
	CreatedAt   time.Time
}

type MessageRequestClaim struct {
	CachedResponse []byte
	Claimed        bool
}
