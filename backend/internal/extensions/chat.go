package extensions

import "time"

type ChatAgent string

const (
	ChatAgentAudit   ChatAgent = "audit"
	ChatAgentContent ChatAgent = "content"
)

func (a ChatAgent) Valid() bool {
	return a == ChatAgentAudit || a == ChatAgentContent
}

type ConversationSummary struct {
	ID          string      `json:"id"`
	Agent       ChatAgent   `json:"agent"`
	ContentType ContentType `json:"content_type,omitempty"`
	Title       string      `json:"title"`
	Starred     bool        `json:"starred"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ChatAttachment struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

type ChatMessage struct {
	ID          string           `json:"id"`
	Role        string           `json:"role"`
	Content     string           `json:"content"`
	Attachments []ChatAttachment `json:"attachments,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
}

type ConversationPage struct {
	Conversations []ConversationSummary `json:"conversations"`
}

type MessagePage struct {
	Messages []ChatMessage `json:"messages"`
}

type SendMessageResponse struct {
	Conversation ConversationSummary `json:"conversation"`
	UserMessage  ChatMessage         `json:"user_message"`
	Reply        ChatMessage         `json:"reply"`
}
