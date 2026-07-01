package history

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strings"

	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

const (
	attachmentMetadataKey = "aeolyzer_attachment_id"
	maxStoredEventBytes   = 2 << 20
	maxStoredStateBytes   = 64 << 10
)

type attachmentContextKey struct{}

func WithAttachmentRefs(ctx context.Context, refs []AttachmentRef) context.Context {
	copied := make([]AttachmentRef, len(refs))
	copy(copied, refs)
	return context.WithValue(ctx, attachmentContextKey{}, copied)
}

func attachmentRefsFromContext(ctx context.Context) []AttachmentRef {
	refs, _ := ctx.Value(attachmentContextKey{}).([]AttachmentRef)
	return refs
}

func prepareEventsForAppend(event *session.Event, refs []AttachmentRef) (*session.Event, *session.Event, error) {
	memoryEvent, err := cloneEvent(event)
	if err != nil {
		return nil, nil, err
	}
	storedEvent, err := cloneEvent(event)
	if err != nil {
		return nil, nil, err
	}
	sanitizeProtectedModelParts(memoryEvent)
	sanitizeProtectedModelParts(storedEvent)
	storedEvent.Actions.StateDelta, err = persistentStateDelta(storedEvent.Actions.StateDelta)
	if err != nil {
		return nil, nil, err
	}

	refIndex := 0
	if storedEvent.Content != nil {
		for _, part := range storedEvent.Content.Parts {
			if part == nil {
				continue
			}
			part.PartMetadata = nil
			if part.InlineData == nil {
				continue
			}
			if refIndex >= len(refs) {
				return nil, nil, ErrInvalidReference
			}
			ref := refs[refIndex]
			if ref.ID == "" || ref.ContentType != part.InlineData.MIMEType || int64(len(part.InlineData.Data)) != ref.Size {
				return nil, nil, ErrInvalidReference
			}
			part.InlineData = nil
			part.PartMetadata = map[string]any{attachmentMetadataKey: ref.ID}
			refIndex++
		}
	}
	if refIndex != len(refs) {
		return nil, nil, ErrInvalidReference
	}
	return memoryEvent, storedEvent, nil
}

func sanitizeProtectedModelParts(event *session.Event) {
	if event == nil || event.Content == nil {
		return
	}
	parts := event.Content.Parts[:0]
	for _, part := range event.Content.Parts {
		if part == nil || part.Thought {
			continue
		}
		part.ThoughtSignature = nil
		parts = append(parts, part)
	}
	event.Content.Parts = parts
}

func cloneEvent(event *session.Event) (*session.Event, error) {
	encoded, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("encode session event: %w", err)
	}
	if len(encoded) > maxStoredEventBytes {
		return nil, errors.New("session event exceeds storage limit")
	}
	var cloned session.Event
	if err := json.Unmarshal(encoded, &cloned); err != nil {
		return nil, fmt.Errorf("decode session event: %w", err)
	}
	return &cloned, nil
}

func marshalEvent(event *session.Event) ([]byte, error) {
	encoded, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("encode session event: %w", err)
	}
	if len(encoded) > maxStoredEventBytes {
		return nil, errors.New("session event exceeds storage limit")
	}
	return encoded, nil
}

func unmarshalEvent(encoded []byte) (*session.Event, error) {
	var event session.Event
	if err := json.Unmarshal(encoded, &event); err != nil {
		return nil, fmt.Errorf("decode session event: %w", err)
	}
	return &event, nil
}

func marshalState(state map[string]any) ([]byte, error) {
	encoded, err := json.Marshal(state)
	if err != nil {
		return nil, fmt.Errorf("encode session state: %w", err)
	}
	if len(encoded) > maxStoredStateBytes {
		return nil, errors.New("session state exceeds storage limit")
	}
	return encoded, nil
}

func unmarshalState(encoded []byte) (map[string]any, error) {
	state := make(map[string]any)
	if err := json.Unmarshal(encoded, &state); err != nil {
		return nil, fmt.Errorf("decode session state: %w", err)
	}
	return state, nil
}

func persistentStateDelta(delta map[string]any) (map[string]any, error) {
	output := make(map[string]any)
	for key, value := range delta {
		if strings.HasPrefix(key, session.KeyPrefixTemp) {
			continue
		}
		if strings.HasPrefix(key, session.KeyPrefixApp) || strings.HasPrefix(key, session.KeyPrefixUser) {
			return nil, errors.New("shared adk state scopes are not enabled")
		}
		if !isPersistentSessionStateKey(key) {
			return nil, errors.New("invalid session state key")
		}
		output[key] = value
	}
	return output, nil
}

func isPersistentSessionStateKey(key string) bool {
	return key != "" && len(key) <= 128 &&
		!strings.HasPrefix(key, session.KeyPrefixTemp) &&
		!strings.HasPrefix(key, session.KeyPrefixApp) &&
		!strings.HasPrefix(key, session.KeyPrefixUser)
}

func mergeState(state map[string]any, delta map[string]any) map[string]any {
	merged := maps.Clone(state)
	maps.Copy(merged, delta)
	return merged
}

func omittedAttachmentPart() *genai.Part {
	return genai.NewPartFromText("[An earlier attachment was omitted from the active context because the attachment context limit was reached.]")
}
