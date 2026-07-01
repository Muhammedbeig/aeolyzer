import type {
  ChatAgent,
  ConversationPage,
  ConversationSummary,
  MessagePage,
  SendMessageResponse,
} from "@/components/chat/types"

const API_URL =
  process.env.NEXT_PUBLIC_AEOLYZER_API_URL?.replace(/\/$/, "") ??
  "http://localhost:8080"

interface APIErrorBody {
  error?: {
    code?: string
    message?: string
  }
}

export class AeolyzerAPIError extends Error {
  readonly status: number
  readonly code: string

  constructor(status: number, code: string, message: string) {
    super(message)
    this.name = "AeolyzerAPIError"
    this.status = status
    this.code = code
  }
}

export async function createConversation(
  agent: ChatAgent,
): Promise<ConversationSummary> {
  return request<ConversationSummary>("/v1/conversations", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ agent }),
  })
}

export async function listConversations(
  agent: ChatAgent,
  signal?: AbortSignal,
): Promise<ConversationSummary[]> {
  const page = await request<ConversationPage>(
    `/v1/conversations?agent=${encodeURIComponent(agent)}`,
    { signal },
  )
  return page.conversations
}

export async function listMessages(
  agent: ChatAgent,
  conversationID: string,
  signal?: AbortSignal,
) {
  const page = await request<MessagePage>(
    `/v1/conversations/${encodeURIComponent(conversationID)}/messages?agent=${encodeURIComponent(agent)}`,
    { signal },
  )
  return page.messages
}

export async function sendMessage(
  agent: ChatAgent,
  conversationID: string,
  text: string,
  files: File[],
): Promise<SendMessageResponse> {
  const body = new FormData()
  body.append("agent", agent)
  if (text) {
    body.append("text", text)
  }
  for (const file of files) {
    body.append("attachments", file)
  }
  return request<SendMessageResponse>(
    `/v1/conversations/${encodeURIComponent(conversationID)}/messages`,
    {
      method: "POST",
      headers: { "Idempotency-Key": crypto.randomUUID() },
      body,
    },
  )
}

export async function updateConversation(
  agent: ChatAgent,
  conversationID: string,
  update: { title?: string; starred?: boolean },
): Promise<ConversationSummary> {
  return request<ConversationSummary>(
    `/v1/conversations/${encodeURIComponent(conversationID)}`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ agent, ...update }),
    },
  )
}

async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    credentials: "include",
    cache: "no-store",
  })
  if (!response.ok) {
    let body: APIErrorBody | undefined
    try {
      body = (await response.json()) as APIErrorBody
    } catch {
      body = undefined
    }
    throw new AeolyzerAPIError(
      response.status,
      body?.error?.code ?? "request_failed",
      body?.error?.message ?? "AEOlyzer could not complete this request.",
    )
  }
  return (await response.json()) as T
}
