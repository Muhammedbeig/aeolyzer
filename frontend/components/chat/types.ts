export type ChatAgent = "audit" | "content"

export interface ConversationSummary {
  id: string
  agent: ChatAgent
  title: string
  starred: boolean
  created_at: string
  updated_at: string
}

export interface ChatAttachment {
  id: string
  name: string
  content_type: string
  size: number
}

export interface ChatMessage {
  id: string
  role: "user" | "assistant"
  content: string
  attachments?: ChatAttachment[]
  created_at: string
  isStreaming?: boolean
}

export interface ConversationPage {
  conversations: ConversationSummary[]
}

export interface MessagePage {
  messages: ChatMessage[]
}

export interface SendMessageResponse {
  conversation: ConversationSummary
  user_message: ChatMessage
  reply: ChatMessage
}
