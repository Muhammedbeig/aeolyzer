"use client" // -> tells Next.js this code must run in the user's browser, allowing us to use interactive features like tracking scrolling

import { useRef, useEffect } from "react" // -> brings in React tools: useRef (for grabbing HTML elements) and useEffect (for running side effects like scrolling)
import { AeolyzerMessage, AeolyzerThinkingIndicator } from "./aeolyzer-message"
import { AeolyzerWelcome } from "./aeolyzer-welcome"
import { ChevronDown } from "lucide-react"

// -> defines exactly what a single chat message should look like
// - `id` -> a unique piece of text to identify this exact message
// - `role` -> tells us who sent it (either the human "user" or the AI "assistant")
// - `content` -> the actual text of the message
// - `isStreaming` -> optional true/false setting that tells us if the AI is currently typing this out
interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  isStreaming?: boolean
}

// -> defines the shape of the data (props) the AeolyzerChatArea block expects to receive
// - `messages` -> a list of Message objects that we defined above
// - `isGenerating` -> true/false switch to tell if the AI is currently thinking
// - `chatTitle` -> optional text to show at the top of the screen (like "New Chat")
// - `onSend` -> the function to run when the user types something new
interface ChatAreaProps {
  messages: Message[]
  isGenerating: boolean
  chatTitle?: string
  onSend: (message: string) => void
}

// -> This is the main scrolling area that holds either the Welcome screen or the list of chat messages
export function AeolyzerChatArea({ messages, isGenerating, chatTitle, onSend }: ChatAreaProps) {
  // - `const messagesEndRef` -> creates a blank sticky note we will attach to the very bottom of our message list (so we can scroll to it later)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  
  // - `const containerRef` -> creates a sticky note attached to the main scrolling box itself
  const containerRef = useRef<HTMLDivElement>(null)

  // -> useEffect runs this auto-scroll logic every time our 'messages' list changes or 'isGenerating' changes
  // Imagine reading a book that automatically slides the page up for you when you reach the bottom
  useEffect(() => {
    // -> if our sticky note actually found the bottom of the list...
    if (messagesEndRef.current) {
      // -> ...tell the browser to smoothly scroll down so that element is visible
      messagesEndRef.current.scrollIntoView({ behavior: "smooth" })
    }
  }, [messages, isGenerating])

  // -> If there are zero messages in the chat...
  if (messages.length === 0) {
    // -> ...show the Welcome screen instead of an empty box
    return (
      <div className="flex-1 flex flex-col" style={{ backgroundColor: "#2b2a27" }}>
        <AeolyzerWelcome onSend={onSend} isGenerating={isGenerating} />
      </div>
    )
  }

  // -> If there are messages, we draw the full chat interface
  return (
    <div className="flex-1 flex flex-col min-h-0" style={{ backgroundColor: "#2b2a27" }}>
      {/* Header with chat title */}
      {/* -> The `&&` means "only draw this box if chatTitle actually exists" */}
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
        ref={containerRef} // -> attaches our finding tool to this scrolling box
        className="flex-1 overflow-y-auto px-4 py-6"
        style={{ scrollbarWidth: "thin", scrollbarColor: "#4a4945 transparent" }}
      >
        <div className="max-w-3xl mx-auto">
          {/* -> loops through our list of messages, turning each one into a visual AeolyzerMessage box */}
          {messages.map((message) => (
            <AeolyzerMessage key={message.id} message={message} />
          ))}
          
          {/* Show thinking indicator when generating and last message is user */}
          {/* -> If the AI is busy AND the very last message in the list came from the human... */}
          {isGenerating && messages[messages.length - 1]?.role === "user" && (
            // -> ...show a little animated loading icon
            <AeolyzerThinkingIndicator />
          )}
          
          {/* -> This empty div is the invisible target we placed our sticky note on. We always scroll down to reach this! */}
          <div ref={messagesEndRef} />
        </div>
      </div>
    </div>
  )
}
