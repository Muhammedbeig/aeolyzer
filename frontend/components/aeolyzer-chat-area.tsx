"use client"

import { useRef, useEffect } from "react"
import { AeolyzerMessage, AeolyzerThinkingIndicator } from "./aeolyzer-message"
import { AeolyzerWelcome } from "./aeolyzer-welcome"
import { ChevronDown } from "lucide-react"

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  isStreaming?: boolean
}

interface ChatAreaProps {
  messages: Message[]
  isGenerating: boolean
  chatTitle?: string
  onSend: (message: string) => void
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
      <div className="flex-1 flex flex-col" style={{ backgroundColor: "#2b2a27" }}>
        <AeolyzerWelcome onSend={onSend} isGenerating={isGenerating} />
      </div>
    )
  }

  return (
    <div className="flex-1 flex flex-col min-h-0" style={{ backgroundColor: "#2b2a27" }}>
      {/* Header with chat title */}
      {chatTitle && (
        <div className="flex items-center justify-between px-6 py-3 border-b flex-shrink-0" style={{ borderColor: "#4a4945" }}>
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium" style={{ color: "#ececec" }}>{chatTitle}</span>
            <ChevronDown size={14} style={{ color: "#a3a29e" }} />
          </div>
          <button 
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors hover:bg-[#d4d4d4]"
            style={{ backgroundColor: "#ececec", color: "#2b2a27" }}
          >
            Share
          </button>
        </div>
      )}

      {/* Messages area - scrollable */}
      <div 
        ref={containerRef}
        className="flex-1 overflow-y-auto px-4 py-6"
        style={{ scrollbarWidth: "thin", scrollbarColor: "#4a4945 transparent" }}
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
