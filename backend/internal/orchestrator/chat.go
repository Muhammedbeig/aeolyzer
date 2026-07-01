package orchestrator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"strings"
	"time"

	"aeolyzer/internal/extensions"
	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
	"aeolyzer/internal/interop/history"
	"aeolyzer/internal/runtime/attachments"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

const (
	auditAppName   = "aeolyzer-audit"
	contentAppName = "aeolyzer-content"
)

type ChatRunner interface {
	Run(context.Context, string, string, *genai.Content) iter.Seq2[*session.Event, error]
}

type ADKChatRunner struct {
	Runner *runner.Runner
}

func (r ADKChatRunner) Run(
	ctx context.Context,
	userID string,
	sessionID string,
	content *genai.Content,
) iter.Seq2[*session.Event, error] {
	return r.Runner.Run(
		ctx,
		userID,
		sessionID,
		content,
		agent.RunConfig{StreamingMode: agent.StreamingModeNone},
	)
}

type ConversationHistory interface {
	CreateConversation(context.Context, string, string) (history.Conversation, error)
	GetConversation(context.Context, string, string, string) (history.Conversation, error)
	ListConversations(context.Context, string, string) ([]history.Conversation, error)
	UpdateConversation(context.Context, string, string, string, *string, *bool) (history.Conversation, error)
	ListMessages(context.Context, string, string, string) ([]history.StoredMessage, error)
	Delete(context.Context, *session.DeleteRequest) error
	SaveAttachment(context.Context, string, string, string, history.AttachmentInput) (history.AttachmentRef, error)
	ClaimMessageRequest(context.Context, string, string, string, string) (history.MessageRequestClaim, error)
	CompleteMessageRequest(context.Context, string, string, string, string, []byte) error
	FailMessageRequest(context.Context, string, string, string, string) error
}

type ChatService struct {
	history ConversationHistory
	runners map[extensions.ChatAgent]ChatRunner
}

type SendChatMessageInput struct {
	Agent          extensions.ChatAgent
	UserID         string
	ConversationID string
	Text           string
	Files          []attachments.File
	IdempotencyKey string
}

func NewChatService(
	conversationHistory ConversationHistory,
	auditRunner ChatRunner,
	contentRunner ChatRunner,
) (*ChatService, error) {
	if conversationHistory == nil {
		return nil, errors.New("conversation history is nil")
	}
	if auditRunner == nil || contentRunner == nil {
		return nil, errors.New("both chat runners are required")
	}
	return &ChatService{
		history: conversationHistory,
		runners: map[extensions.ChatAgent]ChatRunner{
			extensions.ChatAgentAudit:   auditRunner,
			extensions.ChatAgentContent: contentRunner,
		},
	}, nil
}

func (s *ChatService) CreateConversation(
	ctx context.Context,
	userID string,
	chatAgent extensions.ChatAgent,
) (extensions.ConversationSummary, error) {
	appName, err := appNameForAgent(chatAgent)
	if err != nil {
		return extensions.ConversationSummary{}, err
	}
	conversation, err := s.history.CreateConversation(ctx, appName, userID)
	if err != nil {
		return extensions.ConversationSummary{}, err
	}
	return conversationSummary(conversation, chatAgent), nil
}

func (s *ChatService) ListConversations(
	ctx context.Context,
	userID string,
	chatAgent extensions.ChatAgent,
) ([]extensions.ConversationSummary, error) {
	appName, err := appNameForAgent(chatAgent)
	if err != nil {
		return nil, err
	}
	conversations, err := s.history.ListConversations(ctx, appName, userID)
	if err != nil {
		return nil, err
	}
	output := make([]extensions.ConversationSummary, 0, len(conversations))
	for _, conversation := range conversations {
		output = append(output, conversationSummary(conversation, chatAgent))
	}
	return output, nil
}

func (s *ChatService) ListMessages(
	ctx context.Context,
	userID string,
	chatAgent extensions.ChatAgent,
	conversationID string,
) ([]extensions.ChatMessage, error) {
	appName, err := appNameForAgent(chatAgent)
	if err != nil {
		return nil, err
	}
	messages, err := s.history.ListMessages(ctx, appName, userID, conversationID)
	if err != nil {
		return nil, err
	}
	return chatMessages(messages), nil
}

func (s *ChatService) UpdateConversation(
	ctx context.Context,
	userID string,
	chatAgent extensions.ChatAgent,
	conversationID string,
	title *string,
	starred *bool,
) (extensions.ConversationSummary, error) {
	appName, err := appNameForAgent(chatAgent)
	if err != nil {
		return extensions.ConversationSummary{}, err
	}
	conversation, err := s.history.UpdateConversation(
		ctx,
		appName,
		userID,
		conversationID,
		title,
		starred,
	)
	if err != nil {
		return extensions.ConversationSummary{}, err
	}
	return conversationSummary(conversation, chatAgent), nil
}

func (s *ChatService) DeleteConversation(
	ctx context.Context,
	userID string,
	chatAgent extensions.ChatAgent,
	conversationID string,
) error {
	appName, err := appNameForAgent(chatAgent)
	if err != nil {
		return err
	}
	return s.history.Delete(ctx, &session.DeleteRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: conversationID,
	})
}

func (s *ChatService) SendMessage(
	ctx context.Context,
	input SendChatMessageInput,
) (extensions.SendMessageResponse, error) {
	appName, err := appNameForAgent(input.Agent)
	if err != nil {
		return extensions.SendMessageResponse{}, err
	}
	claim, err := s.history.ClaimMessageRequest(
		ctx,
		appName,
		input.UserID,
		input.ConversationID,
		input.IdempotencyKey,
	)
	if err != nil {
		return extensions.SendMessageResponse{}, err
	}
	if len(claim.CachedResponse) > 0 {
		var cached extensions.SendMessageResponse
		if err := json.Unmarshal(claim.CachedResponse, &cached); err != nil {
			return extensions.SendMessageResponse{}, fmt.Errorf("decode cached message response: %w", err)
		}
		return cached, nil
	}

	response, err := s.runMessage(ctx, appName, input)
	if err != nil {
		failureContext, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Second)
		defer cancel()
		_ = s.history.FailMessageRequest(
			failureContext,
			appName,
			input.UserID,
			input.ConversationID,
			input.IdempotencyKey,
		)
		return extensions.SendMessageResponse{}, err
	}
	encoded, err := json.Marshal(response)
	if err != nil {
		return extensions.SendMessageResponse{}, fmt.Errorf("encode message response: %w", err)
	}
	if err := s.history.CompleteMessageRequest(
		ctx,
		appName,
		input.UserID,
		input.ConversationID,
		input.IdempotencyKey,
		encoded,
	); err != nil {
		return extensions.SendMessageResponse{}, err
	}
	return response, nil
}

func (s *ChatService) runMessage(
	ctx context.Context,
	appName string,
	input SendChatMessageInput,
) (extensions.SendMessageResponse, error) {
	conversation, err := s.history.GetConversation(
		ctx,
		appName,
		input.UserID,
		input.ConversationID,
	)
	if err != nil {
		return extensions.SendMessageResponse{}, err
	}
	refs := make([]history.AttachmentRef, 0, len(input.Files))
	parts := make([]*genai.Part, 0, len(input.Files)+1)
	if strings.TrimSpace(input.Text) != "" {
		parts = append(parts, genai.NewPartFromText(input.Text))
	}
	for _, file := range input.Files {
		ref, err := s.history.SaveAttachment(
			ctx,
			appName,
			input.UserID,
			input.ConversationID,
			history.AttachmentInput{
				Name:        file.Name,
				ContentType: file.ContentType,
				Data:        file.Data,
				SHA256:      file.SHA256,
			},
		)
		if err != nil {
			return extensions.SendMessageResponse{}, err
		}
		refs = append(refs, ref)
		parts = append(parts, genai.NewPartFromBytes(file.Data, file.ContentType))
	}
	content := &genai.Content{Role: genai.RoleUser, Parts: parts}
	runContext := history.WithAttachmentRefs(ctx, refs)
	var replyText string
	for event, runErr := range s.runners[input.Agent].Run(
		runContext,
		input.UserID,
		input.ConversationID,
		content,
	) {
		if runErr != nil {
			return extensions.SendMessageResponse{}, fmt.Errorf("run chat agent: %w", runErr)
		}
		if event == nil || !event.IsFinalResponse() {
			continue
		}
		if text := eventText(event); text != "" {
			replyText = text
		}
	}
	if replyText == "" {
		return extensions.SendMessageResponse{}, errors.New("chat agent returned no displayable response")
	}
	intent := contracts.IntentAuditSEO
	if input.Agent == extensions.ChatAgentContent {
		intent = contracts.IntentDraftArticle
	}
	replyText, err = middleware.GuardOutboundResponse(replyText, intent)
	if err != nil {
		return extensions.SendMessageResponse{}, fmt.Errorf("guard chat response: %w", err)
	}

	if conversation.Title == "New chat" {
		title := titleForMessage(input.Text, input.Files)
		conversation, err = s.history.UpdateConversation(
			ctx,
			appName,
			input.UserID,
			input.ConversationID,
			&title,
			nil,
		)
		if err != nil {
			return extensions.SendMessageResponse{}, err
		}
	} else {
		conversation, err = s.history.GetConversation(
			ctx,
			appName,
			input.UserID,
			input.ConversationID,
		)
		if err != nil {
			return extensions.SendMessageResponse{}, err
		}
	}
	messages, err := s.history.ListMessages(ctx, appName, input.UserID, input.ConversationID)
	if err != nil {
		return extensions.SendMessageResponse{}, err
	}
	userMessage, reply, err := latestTurn(chatMessages(messages))
	if err != nil {
		return extensions.SendMessageResponse{}, err
	}
	reply.Content = replyText
	return extensions.SendMessageResponse{
		Conversation: conversationSummary(conversation, input.Agent),
		UserMessage:  userMessage,
		Reply:        reply,
	}, nil
}

func appNameForAgent(chatAgent extensions.ChatAgent) (string, error) {
	switch chatAgent {
	case extensions.ChatAgentAudit:
		return auditAppName, nil
	case extensions.ChatAgentContent:
		return contentAppName, nil
	default:
		return "", errors.New("unknown chat agent")
	}
}

func conversationSummary(
	conversation history.Conversation,
	chatAgent extensions.ChatAgent,
) extensions.ConversationSummary {
	return extensions.ConversationSummary{
		ID:        conversation.ID,
		Agent:     chatAgent,
		Title:     conversation.Title,
		Starred:   conversation.Starred,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	}
}

func chatMessages(messages []history.StoredMessage) []extensions.ChatMessage {
	output := make([]extensions.ChatMessage, 0, len(messages))
	for _, message := range messages {
		role := "assistant"
		if message.Author == "user" {
			role = "user"
		}
		attachments := make([]extensions.ChatAttachment, 0, len(message.Attachments))
		for _, attachment := range message.Attachments {
			attachments = append(attachments, extensions.ChatAttachment{
				ID:          attachment.ID,
				Name:        attachment.Name,
				ContentType: attachment.ContentType,
				Size:        attachment.Size,
			})
		}
		output = append(output, extensions.ChatMessage{
			ID:          message.ID,
			Role:        role,
			Content:     message.Text,
			Attachments: attachments,
			CreatedAt:   message.CreatedAt,
		})
	}
	return output
}

func eventText(event *session.Event) string {
	if event == nil || event.Content == nil {
		return ""
	}
	var output strings.Builder
	for _, part := range event.Content.Parts {
		if part == nil || part.Thought || part.Text == "" {
			continue
		}
		if output.Len() > 0 {
			output.WriteByte('\n')
		}
		output.WriteString(part.Text)
	}
	return strings.TrimSpace(output.String())
}

func titleForMessage(text string, files []attachments.File) string {
	title := strings.Join(strings.Fields(text), " ")
	if title == "" && len(files) > 0 {
		title = files[0].Name
	}
	runes := []rune(title)
	if len(runes) > 80 {
		title = string(runes[:80])
	}
	if title == "" {
		return "New chat"
	}
	return title
}

func latestTurn(messages []extensions.ChatMessage) (extensions.ChatMessage, extensions.ChatMessage, error) {
	var userMessage extensions.ChatMessage
	var reply extensions.ChatMessage
	for index := len(messages) - 1; index >= 0; index-- {
		message := messages[index]
		if reply.ID == "" && message.Role == "assistant" {
			reply = message
			continue
		}
		if reply.ID != "" && message.Role == "user" {
			userMessage = message
			break
		}
	}
	if userMessage.ID == "" || reply.ID == "" {
		return extensions.ChatMessage{}, extensions.ChatMessage{}, errors.New("persisted chat turn is incomplete")
	}
	return userMessage, reply, nil
}
