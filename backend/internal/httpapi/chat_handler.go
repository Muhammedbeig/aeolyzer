package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
	"unicode"

	"aeolyzer/internal/extensions"
	"aeolyzer/internal/intake"
	"aeolyzer/internal/interop/history"
	"aeolyzer/internal/observability"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/runtime/attachments"
)

const (
	maxMultipartOverhead = 1 << 20
	chatRequestTimeout   = 60 * time.Second
)

type ChatAPI interface {
	CreateConversation(context.Context, string, extensions.ChatAgent) (extensions.ConversationSummary, error)
	ListConversations(context.Context, string, extensions.ChatAgent) ([]extensions.ConversationSummary, error)
	ListMessages(context.Context, string, extensions.ChatAgent, string) ([]extensions.ChatMessage, error)
	UpdateConversation(context.Context, string, extensions.ChatAgent, string, *string, *bool) (extensions.ConversationSummary, error)
	DeleteConversation(context.Context, string, extensions.ChatAgent, string) error
	SendMessage(context.Context, orchestrator.SendChatMessageInput) (extensions.SendMessageResponse, error)
}

type ChatHandler struct {
	chat          ChatAPI
	files         *attachments.Processor
	identity      *GuestIdentity
	events        *observability.Sink
	logger        *slog.Logger
	allowedOrigin string
	now           func() time.Time
}

type createConversationRequest struct {
	Agent extensions.ChatAgent `json:"agent"`
}

type updateConversationRequest struct {
	Agent   extensions.ChatAgent `json:"agent"`
	Title   *string              `json:"title,omitempty"`
	Starred *bool                `json:"starred,omitempty"`
}

type guestContextKey struct{}

func NewChatHandler(
	chat ChatAPI,
	files *attachments.Processor,
	identity *GuestIdentity,
	events *observability.Sink,
	logger *slog.Logger,
	allowedOrigin string,
) (*ChatHandler, error) {
	if chat == nil {
		return nil, errors.New("chat api is nil")
	}
	if files == nil {
		return nil, errors.New("attachment processor is nil")
	}
	if identity == nil {
		return nil, errors.New("guest identity is nil")
	}
	return &ChatHandler{
		chat:          chat,
		files:         files,
		identity:      identity,
		events:        events,
		logger:        loggerOrDefault(logger),
		allowedOrigin: strings.TrimRight(allowedOrigin, "/"),
		now:           time.Now,
	}, nil
}

func (h *ChatHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/conversations", h.createConversation)
	mux.HandleFunc("GET /v1/conversations", h.listConversations)
	mux.HandleFunc("GET /v1/conversations/{conversationID}/messages", h.listMessages)
	mux.HandleFunc("POST /v1/conversations/{conversationID}/messages", h.sendMessage)
	mux.HandleFunc("PATCH /v1/conversations/{conversationID}", h.updateConversation)
	mux.HandleFunc("DELETE /v1/conversations/{conversationID}", h.deleteConversation)
	return h.withMiddleware(mux)
}

func (h *ChatHandler) createConversation(response http.ResponseWriter, request *http.Request) {
	var input createConversationRequest
	if err := decodeJSON(response, request, &input); err != nil || !input.Agent.Valid() {
		writeError(response, http.StatusBadRequest, "invalid_request", "Choose a valid agent and try again.")
		return
	}
	conversation, err := h.chat.CreateConversation(request.Context(), guestID(request), input.Agent)
	if err != nil {
		h.handleError(response, request, err, "conversation_create")
		return
	}
	h.record("conversation_create", "succeeded")
	writeJSON(response, http.StatusCreated, conversation)
}

func (h *ChatHandler) listConversations(response http.ResponseWriter, request *http.Request) {
	chatAgent, err := queryAgent(request)
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid_request", "Choose a valid agent and try again.")
		return
	}
	conversations, err := h.chat.ListConversations(request.Context(), guestID(request), chatAgent)
	if err != nil {
		h.handleError(response, request, err, "conversation_list")
		return
	}
	writeJSON(response, http.StatusOK, extensions.ConversationPage{Conversations: conversations})
}

func (h *ChatHandler) listMessages(response http.ResponseWriter, request *http.Request) {
	chatAgent, err := queryAgent(request)
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid_request", "Choose a valid agent and try again.")
		return
	}
	messages, err := h.chat.ListMessages(
		request.Context(),
		guestID(request),
		chatAgent,
		request.PathValue("conversationID"),
	)
	if err != nil {
		h.handleError(response, request, err, "message_list")
		return
	}
	writeJSON(response, http.StatusOK, extensions.MessagePage{Messages: messages})
}

func (h *ChatHandler) updateConversation(response http.ResponseWriter, request *http.Request) {
	var input updateConversationRequest
	if err := decodeJSON(response, request, &input); err != nil ||
		!input.Agent.Valid() ||
		!validTitle(input.Title) ||
		(input.Title == nil && input.Starred == nil) {
		writeError(response, http.StatusBadRequest, "invalid_request", "Check the conversation update and try again.")
		return
	}
	conversation, err := h.chat.UpdateConversation(
		request.Context(),
		guestID(request),
		input.Agent,
		request.PathValue("conversationID"),
		input.Title,
		input.Starred,
	)
	if err != nil {
		h.handleError(response, request, err, "conversation_update")
		return
	}
	h.record("conversation_update", "succeeded")
	writeJSON(response, http.StatusOK, conversation)
}

func (h *ChatHandler) deleteConversation(response http.ResponseWriter, request *http.Request) {
	chatAgent, err := queryAgent(request)
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid_request", "Choose a valid agent and try again.")
		return
	}
	err = h.chat.DeleteConversation(
		request.Context(),
		guestID(request),
		chatAgent,
		request.PathValue("conversationID"),
	)
	if err != nil {
		h.handleError(response, request, err, "conversation_delete")
		return
	}
	h.record("conversation_delete", "succeeded")
	response.WriteHeader(http.StatusNoContent)
}

func (h *ChatHandler) sendMessage(response http.ResponseWriter, request *http.Request) {
	input, err := h.readMessage(response, request)
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid_message", "Check the message and attachments, then try again.")
		return
	}
	input.UserID = guestID(request)
	input.ConversationID = request.PathValue("conversationID")
	input.IdempotencyKey = request.Header.Get("Idempotency-Key")
	result, err := h.chat.SendMessage(request.Context(), input)
	if err != nil {
		h.handleError(response, request, err, "message_send")
		return
	}
	h.record("message_send", "succeeded")
	writeJSON(response, http.StatusOK, result)
}

func (h *ChatHandler) readMessage(
	response http.ResponseWriter,
	request *http.Request,
) (orchestrator.SendChatMessageInput, error) {
	request.Body = http.MaxBytesReader(
		response,
		request.Body,
		attachments.DefaultMaxTotalBytes+maxMultipartOverhead,
	)
	reader, err := request.MultipartReader()
	if err != nil {
		return orchestrator.SendChatMessageInput{}, err
	}
	var input orchestrator.SendChatMessageInput
	var agentSeen bool
	var textSeen bool
	var totalFileBytes int64
	for partCount := 0; ; partCount++ {
		if partCount > intake.MaxChatAttachments+2 {
			return orchestrator.SendChatMessageInput{}, errors.New("too many multipart fields")
		}
		part, err := reader.NextPart()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return orchestrator.SendChatMessageInput{}, err
		}
		switch part.FormName() {
		case "agent":
			if agentSeen || part.FileName() != "" {
				_ = part.Close()
				return orchestrator.SendChatMessageInput{}, errors.New("invalid agent field")
			}
			value, err := readPart(part, 32)
			if err != nil {
				return orchestrator.SendChatMessageInput{}, err
			}
			input.Agent = extensions.ChatAgent(value)
			agentSeen = true
		case "text":
			if textSeen || part.FileName() != "" {
				_ = part.Close()
				return orchestrator.SendChatMessageInput{}, errors.New("invalid text field")
			}
			value, err := readPart(part, int64(intake.MaxChatTextRunes*4))
			if err != nil {
				return orchestrator.SendChatMessageInput{}, err
			}
			input.Text = value
			textSeen = true
		case "attachments":
			if part.FileName() == "" || len(input.Files) >= intake.MaxChatAttachments {
				_ = part.Close()
				return orchestrator.SendChatMessageInput{}, errors.New("invalid attachment field")
			}
			data, err := readPartBytes(part, h.files.MaxFileBytes)
			if err != nil {
				return orchestrator.SendChatMessageInput{}, err
			}
			totalFileBytes += int64(len(data))
			if totalFileBytes > attachments.DefaultMaxTotalBytes {
				return orchestrator.SendChatMessageInput{}, attachments.ErrFileTooLarge
			}
			file, err := h.files.Process(part.FileName(), data)
			if err != nil {
				return orchestrator.SendChatMessageInput{}, err
			}
			if err := intake.ValidateAttachmentContent(file.ContentType, file.Data); err != nil {
				return orchestrator.SendChatMessageInput{}, err
			}
			input.Files = append(input.Files, file)
		default:
			_ = part.Close()
			return orchestrator.SendChatMessageInput{}, errors.New("unknown multipart field")
		}
	}
	if !agentSeen || !input.Agent.Valid() {
		return orchestrator.SendChatMessageInput{}, errors.New("invalid agent")
	}
	if err := intake.ValidateChatMessage(input.Text, len(input.Files)); err != nil {
		return orchestrator.SendChatMessageInput{}, err
	}
	if key := request.Header.Get("Idempotency-Key"); len(key) < 16 || len(key) > 128 {
		return orchestrator.SendChatMessageInput{}, errors.New("invalid idempotency key")
	}
	return input, nil
}

func (h *ChatHandler) handleError(
	response http.ResponseWriter,
	request *http.Request,
	err error,
	eventType string,
) {
	switch {
	case errors.Is(err, history.ErrNotFound):
		writeError(response, http.StatusNotFound, "not_found", "This conversation was not found.")
	case errors.Is(err, history.ErrRequestPending):
		writeError(response, http.StatusConflict, "request_in_progress", "This message is already being processed.")
	case errors.Is(err, history.ErrRequestFailed):
		writeError(response, http.StatusConflict, "request_failed", "The previous attempt failed. Send the message again.")
	case errors.Is(err, history.ErrConflict):
		writeError(response, http.StatusConflict, "conflict", "The conversation changed. Refresh and try again.")
	case errors.Is(err, context.DeadlineExceeded):
		writeError(response, http.StatusGatewayTimeout, "timeout", "The agent took too long to respond. Try again.")
	default:
		h.logger.ErrorContext(request.Context(), "chat request failed", "event_type", eventType, "error", err)
		writeError(response, http.StatusBadGateway, "agent_unavailable", "The agent could not complete that request.")
	}
	h.record(eventType, "failed")
}

func (h *ChatHandler) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("Cache-Control", "no-store")
		response.Header().Set("Referrer-Policy", "no-referrer")

		origin := strings.TrimRight(request.Header.Get("Origin"), "/")
		if origin != "" && origin != h.allowedOrigin {
			writeError(response, http.StatusForbidden, "origin_denied", "This origin is not allowed.")
			return
		}
		if origin == h.allowedOrigin && origin != "" {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Access-Control-Allow-Credentials", "true")
			response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Idempotency-Key")
			response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			response.Header().Set("Vary", "Origin")
		}
		if request.Method == http.MethodOptions {
			if origin == "" || origin != h.allowedOrigin {
				writeError(response, http.StatusForbidden, "origin_denied", "This origin is not allowed.")
				return
			}
			response.WriteHeader(http.StatusNoContent)
			return
		}

		userID := h.identity.Resolve(response, request)
		requestContext, cancel := context.WithTimeout(request.Context(), chatRequestTimeout)
		defer cancel()
		request = request.WithContext(context.WithValue(requestContext, guestContextKey{}, userID))
		next.ServeHTTP(response, request)
	})
}

func (h *ChatHandler) record(eventType, outcome string) {
	if h.events == nil {
		return
	}
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return
	}
	h.events.Record(observability.Event{
		TraceID:   hex.EncodeToString(random[:]),
		EventType: eventType,
		Outcome:   outcome,
		At:        h.now().UTC(),
	})
}

func queryAgent(request *http.Request) (extensions.ChatAgent, error) {
	chatAgent := extensions.ChatAgent(request.URL.Query().Get("agent"))
	if !chatAgent.Valid() {
		return "", errors.New("invalid chat agent")
	}
	return chatAgent, nil
}

func guestID(request *http.Request) string {
	userID, _ := request.Context().Value(guestContextKey{}).(string)
	return userID
}

func readPart(part *multipart.Part, limit int64) (string, error) {
	data, err := readPartBytes(part, limit)
	return string(data), err
}

func readPartBytes(part *multipart.Part, limit int64) ([]byte, error) {
	defer part.Close()
	data, err := io.ReadAll(io.LimitReader(part, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, errors.New("multipart field exceeds limit")
	}
	return data, nil
}

func validTitle(title *string) bool {
	if title == nil {
		return true
	}
	if len([]rune(*title)) > 80 {
		return false
	}
	for _, character := range *title {
		if unicode.IsControl(character) && character != '\t' {
			return false
		}
	}
	return true
}
