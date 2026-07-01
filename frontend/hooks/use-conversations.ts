"use client"

import { useCallback, useEffect, useMemo, useState } from "react"
import type {
  ChatAgent,
  ChatMessage,
  ContentType,
  ConversationSummary,
} from "@/components/chat/types"
import {
  AeolyzerAPIError,
  createConversation,
  listConversations,
  listMessages,
  sendMessage,
  updateConversation,
} from "@/lib/aeolyzer-api"

type AgentState<T> = Record<ChatAgent, T>

const EMPTY_CONVERSATIONS: AgentState<ConversationSummary[]> = {
  audit: [],
  content: [],
}

const EMPTY_MESSAGES: AgentState<ChatMessage[]> = {
  audit: [],
  content: [],
}

const EMPTY_ACTIVE_IDS: AgentState<string | undefined> = {
  audit: undefined,
  content: undefined,
}

const EMPTY_LOADING: AgentState<boolean> = {
  audit: false,
  content: false,
}

export function useConversations(initialAgent: ChatAgent = "audit") {
  const [agent, setAgent] = useState<ChatAgent>(initialAgent)
  const [contentType, setContentType] = useState<ContentType>("article")
  const [conversations, setConversations] =
    useState<AgentState<ConversationSummary[]>>(EMPTY_CONVERSATIONS)
  const [messages, setMessages] =
    useState<AgentState<ChatMessage[]>>(EMPTY_MESSAGES)
  const [activeIDs, setActiveIDs] =
    useState<AgentState<string | undefined>>(EMPTY_ACTIVE_IDS)
  const [isGenerating, setIsGenerating] =
    useState<AgentState<boolean>>(EMPTY_LOADING)

  useEffect(() => {
    const controller = new AbortController()
    Promise.all(
      (["audit", "content"] as const).map(async (chatAgent) => {
        const items = await listConversations(chatAgent, controller.signal)
        return [chatAgent, items] as const
      }),
    )
      .then((results) => {
        setConversations((current) => {
          const next = { ...current }
          for (const [chatAgent, items] of results) {
            next[chatAgent] = items
          }
          return next
        })
      })
      .catch((error: unknown) => {
        if (!(error instanceof DOMException && error.name === "AbortError")) {
          console.error("Could not load conversation history", error)
        }
      })
    return () => controller.abort()
  }, [])

  const beginNewConversation = useCallback(() => {
    setActiveIDs((current) => ({ ...current, [agent]: undefined }))
    setMessages((current) => ({ ...current, [agent]: [] }))
    setIsGenerating((current) => ({ ...current, [agent]: false }))
    if (agent === "content") {
      setContentType("article")
    }
  }, [agent])

  const selectConversation = useCallback(
    async (conversation: ConversationSummary) => {
      setAgent(conversation.agent)
      setActiveIDs((current) => ({
        ...current,
        [conversation.agent]: conversation.id,
      }))
      setMessages((current) => ({ ...current, [conversation.agent]: [] }))
      if (conversation.agent === "content") {
        setContentType(conversation.content_type ?? "article")
      }
      try {
        const items = await listMessages(conversation.agent, conversation.id)
        setMessages((current) => ({
          ...current,
          [conversation.agent]: items,
        }))
      } catch (error) {
        setMessages((current) => ({
          ...current,
          [conversation.agent]: [
            errorMessage(error, `history-error-${conversation.id}`),
          ],
        }))
      }
    },
    [],
  )

  const submitMessage = useCallback(
    async (
      text: string,
      files: File[] = [],
      requestedContentType?: ContentType,
    ) => {
      if (isGenerating[agent]) {
        return
      }
      const pendingID = `pending-${crypto.randomUUID()}`
      const optimisticMessage: ChatMessage = {
        id: pendingID,
        role: "user",
        content: text,
        attachments: files.map((file, index) => ({
          id: `${pendingID}-${index}`,
          name: file.name,
          content_type: file.type || "application/octet-stream",
          size: file.size,
        })),
        created_at: new Date().toISOString(),
      }
      setMessages((current) => ({
        ...current,
        [agent]: [...current[agent], optimisticMessage],
      }))
      setIsGenerating((current) => ({ ...current, [agent]: true }))

      try {
        const activeContentType =
          agent === "content"
            ? requestedContentType ?? contentType
            : undefined
        let conversationID = activeIDs[agent]
        if (!conversationID) {
          const created = await createConversation(agent, activeContentType)
          conversationID = created.id
          setActiveIDs((current) => ({ ...current, [agent]: created.id }))
          setConversations((current) => ({
            ...current,
            [agent]: upsertConversation(current[agent], created),
          }))
        }
        const result = await sendMessage(
          agent,
          conversationID,
          text,
          files,
          activeContentType,
        )
        setMessages((current) => {
          const withoutPending = current[agent].filter(
            (message) => message.id !== pendingID,
          )
          return {
            ...current,
            [agent]: [
              ...withoutPending,
              result.user_message,
              result.reply,
            ],
          }
        })
        setConversations((current) => ({
          ...current,
          [agent]: upsertConversation(
            current[agent],
            result.conversation,
          ),
        }))
      } catch (error) {
        setMessages((current) => ({
          ...current,
          [agent]: [
            ...current[agent],
            errorMessage(error, `error-${crypto.randomUUID()}`),
          ],
        }))
      } finally {
        setIsGenerating((current) => ({ ...current, [agent]: false }))
      }
    },
    [activeIDs, agent, contentType, isGenerating],
  )

  const toggleStar = useCallback(
    async (conversation: ConversationSummary) => {
      try {
        const updated = await updateConversation(
          conversation.agent,
          conversation.id,
          { starred: !conversation.starred },
        )
        setConversations((current) => ({
          ...current,
          [conversation.agent]: upsertConversation(
            current[conversation.agent],
            updated,
          ),
        }))
      } catch (error) {
        console.error("Could not update conversation", error)
      }
    },
    [],
  )

  const currentConversation = useMemo(
    () =>
      conversations[agent].find(
        (conversation) => conversation.id === activeIDs[agent],
      ),
    [activeIDs, agent, conversations],
  )

  return {
    agent,
    setAgent,
    contentType,
    setContentType,
    conversations,
    allConversations: [
      ...conversations.audit,
      ...conversations.content,
    ],
    currentConversations: conversations[agent],
    currentConversation,
    activeConversationID: activeIDs[agent],
    messages: messages[agent],
    isGenerating: isGenerating[agent],
    beginNewConversation,
    selectConversation,
    submitMessage,
    toggleStar,
  }
}

function upsertConversation(
  conversations: ConversationSummary[],
  value: ConversationSummary,
) {
  const next = [
    value,
    ...conversations.filter(
      (conversation) => conversation.id !== value.id,
    ),
  ]
  return next.sort((left, right) => {
    if (left.starred !== right.starred) {
      return left.starred ? -1 : 1
    }
    return Date.parse(right.updated_at) - Date.parse(left.updated_at)
  })
}

function errorMessage(error: unknown, id: string): ChatMessage {
  const content =
    error instanceof AeolyzerAPIError
      ? error.message
      : "AEOlyzer could not complete that request. Please try again."
  return {
    id,
    role: "assistant",
    content,
    created_at: new Date().toISOString(),
  }
}
