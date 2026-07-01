import type {
  ChatAgent,
  ContentType,
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
  contentType?: ContentType,
): Promise<ConversationSummary> {
  return aeolyzerRequest<ConversationSummary>("/v1/conversations", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      agent,
      ...(agent === "content" ? { content_type: contentType ?? "article" } : {}),
    }),
  })
}

export async function listConversations(
  agent: ChatAgent,
  signal?: AbortSignal,
): Promise<ConversationSummary[]> {
  const page = await aeolyzerRequest<ConversationPage>(
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
  const page = await aeolyzerRequest<MessagePage>(
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
  contentType?: ContentType,
): Promise<SendMessageResponse> {
  const body = new FormData()
  body.append("agent", agent)
  if (text) {
    body.append("text", text)
  }
  if (agent === "content") {
    body.append("content_type", contentType ?? "article")
  }
  for (const file of files) {
    body.append("attachments", file)
  }
  return aeolyzerRequest<SendMessageResponse>(
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
  return aeolyzerRequest<ConversationSummary>(
    `/v1/conversations/${encodeURIComponent(conversationID)}`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ agent, ...update }),
    },
  )
}

export async function aeolyzerRequest<T>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
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
