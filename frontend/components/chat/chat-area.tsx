"use client"

import { useRef, useEffect } from "react"
import { AeolyzerMessage, AeolyzerThinkingIndicator } from "./chat-message"
import { AeolyzerWelcome } from "./welcome-screen"
import { ChevronDown } from "lucide-react"
import type { ChatMessage } from "./types"

interface ChatAreaProps {
  messages: ChatMessage[]
  isGenerating: boolean
  chatTitle?: string
  onSend: (message: string, files?: File[]) => void
}

export function AeolyzerChatArea({ messages, isGenerating, chatTitle, onSend }: ChatAreaProps) {
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  // Auto-scroll to bottom when new messages arrive or content updates
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [messages, isGenerating])

  // If no messages, show welcome screen
  if (messages.length === 0) {
    return (
      <div className="flex-1 flex flex-col bg-background">
        <AeolyzerWelcome onSend={onSend} isGenerating={isGenerating} />
      </div>
    )
  }

  return (
    <div className="flex-1 flex flex-col min-h-0 bg-background">
      {/* Header with chat title */}
      {chatTitle && (
        <div className="flex items-center justify-between px-6 py-3 border-b-[0.5px] flex-shrink-0 border-black/10 dark:border-white/10">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium text-foreground">{chatTitle}</span>
            <ChevronDown size={14} className="text-muted-foreground" />
          </div>
          <button 
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors hover:bg-muted bg-foreground text-background"
          >
            Share
          </button>
        </div>
      )}

      {/* Messages area - scrollable */}
      <div 
        ref={containerRef}
        className="custom-scrollbar flex-1 overflow-y-auto px-4 py-6"
      >
        <div className="max-w-3xl mx-auto">
          {messages.map((message) => (
            <AeolyzerMessage key={message.id} message={message} />
          ))}
          
          {/* Show thinking indicator when generating and last message is user */}
          {isGenerating && messages[messages.length - 1]?.role === "user" && (
            <AeolyzerThinkingIndicator />
          )}
          
          <div ref={messagesEndRef} />
        </div>
      </div>
    </div>
  )
}
