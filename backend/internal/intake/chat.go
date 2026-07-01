package intake

import (
	"errors"
	"strings"
	"unicode/utf8"

	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
)

const (
	MaxChatTextRunes         = 12_000
	MaxChatAttachments       = 5
	attachmentScanChunkBytes = 16 << 10
)

var (
	ErrEmptyChatMessage   = errors.New("chat message is empty")
	ErrChatMessageTooLong = errors.New("chat message is too long")
	ErrTooManyAttachments = errors.New("too many attachments")
	ErrUnsafeChatMessage  = errors.New("chat message was blocked")
)

func ValidateChatMessage(text string, attachmentCount int) error {
	if attachmentCount < 0 || attachmentCount > MaxChatAttachments {
		return ErrTooManyAttachments
	}
	if strings.TrimSpace(text) == "" && attachmentCount == 0 {
		return ErrEmptyChatMessage
	}
	if !utf8.ValidString(text) || utf8.RuneCountInString(text) > MaxChatTextRunes {
		return ErrChatMessageTooLong
	}
	if err := middleware.CheckForPromptInjection(contracts.SanitizedInput{RawText: text}); err != nil {
		return ErrUnsafeChatMessage
	}
	return nil
}

func ValidateAttachmentContent(contentType string, data []byte) error {
	if !strings.HasPrefix(contentType, "text/") && contentType != "application/json" {
		return nil
	}
	for start := 0; start < len(data); start += attachmentScanChunkBytes {
		end := min(start+attachmentScanChunkBytes+256, len(data))
		if err := middleware.CheckForPromptInjection(contracts.SanitizedInput{
			RawText: string(data[start:end]),
		}); err != nil {
			return ErrUnsafeChatMessage
		}
	}
	return nil
}
