"use client"

import { useState, useEffect } from "react"
import { Copy, ThumbsUp, ThumbsDown, RotateCcw } from "lucide-react"
import { AeolyzerLogoAnimated } from "./aeolyzer-logo"
import { cn } from "@/lib/utils"

interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  isStreaming?: boolean
}

interface AeolyzerMessageProps {
  message: Message
}

export function AeolyzerMessage({ message }: AeolyzerMessageProps) {
  const [displayedContent, setDisplayedContent] = useState("")
  const [isComplete, setIsComplete] = useState(false)

  useEffect(() => {
    if (message.role === "assistant" && message.isStreaming) {
      setDisplayedContent("")
      setIsComplete(false)
      let currentIndex = 0
      const content = message.content
      
      const interval = setInterval(() => {
        if (currentIndex < content.length) {
          // Add characters in small chunks for smooth animation
          const chunkSize = Math.floor(Math.random() * 3) + 1
          const nextChunk = content.slice(currentIndex, currentIndex + chunkSize)
          setDisplayedContent(prev => prev + nextChunk)
          currentIndex += chunkSize
        } else {
          clearInterval(interval)
          setIsComplete(true)
        }
      }, 15) // Fast typing speed like Claude

      return () => clearInterval(interval)
    } else {
      setDisplayedContent(message.content)
      setIsComplete(true)
    }
  }, [message.content, message.isStreaming, message.role])

  if (message.role === "user") {
    return (
      <div className="flex justify-end mb-6 animate-fade-in-up">
        <div 
          className="max-w-[80%] px-4 py-3 rounded-2xl"
          style={{ backgroundColor: "#393836" }}
        >
          <p className="text-[15px] leading-relaxed" style={{ color: "#ececec" }}>
            {message.content}
          </p>
        </div>
      </div>
    )
  }

  return (
    <div className="mb-6 animate-fade-in-up">
      {/* Thinking indicator when starting */}
      {message.isStreaming && displayedContent.length === 0 && (
        <div className="flex items-center gap-3 mb-4">
          <AeolyzerLogoAnimated size={28} />
          <span className="text-sm" style={{ color: "#a3a29e" }}>Thinking...</span>
        </div>
      )}

      {/* Message content */}
      <div className="max-w-none">
        <div className="prose prose-invert max-w-none">
          {renderFormattedContent(displayedContent)}
        </div>
        
        {/* Typing cursor */}
        {message.isStreaming && !isComplete && displayedContent.length > 0 && (
          <span className="inline-block w-2 h-5 ml-0.5 bg-[#e07b53] animate-pulse" />
        )}
      </div>

      {/* Action buttons - show after complete */}
      {isComplete && displayedContent.length > 0 && (
        <div className="flex items-center gap-1 mt-4">
          <button 
            className="p-2 rounded-lg transition-colors hover:bg-[#252422]"
            style={{ color: "#a3a29e" }}
          >
            <Copy size={16} />
          </button>
          <button 
            className="p-2 rounded-lg transition-colors hover:bg-[#252422]"
            style={{ color: "#a3a29e" }}
          >
            <ThumbsUp size={16} />
          </button>
          <button 
            className="p-2 rounded-lg transition-colors hover:bg-[#252422]"
            style={{ color: "#a3a29e" }}
          >
            <ThumbsDown size={16} />
          </button>
          <button 
            className="p-2 rounded-lg transition-colors hover:bg-[#252422]"
            style={{ color: "#a3a29e" }}
          >
            <RotateCcw size={16} />
          </button>
        </div>
      )}
    </div>
  )
}

function renderFormattedContent(content: string) {
  // Split content by double newlines for paragraphs
  const paragraphs = content.split(/\n\n/)
  
  return paragraphs.map((paragraph, pIndex) => {
    // Check for headers
    if (paragraph.startsWith("**") && paragraph.endsWith("**")) {
      const headerText = paragraph.slice(2, -2)
      return (
        <h3 
          key={pIndex} 
          className="text-lg font-semibold mt-6 mb-3 first:mt-0"
          style={{ color: "#ececec" }}
        >
          {headerText}
        </h3>
      )
    }

    // Check for section headers like "Writing & Communication"
    if (paragraph.includes("**") && paragraph.includes(":")) {
      const parts = paragraph.split("**")
      return (
        <div key={pIndex} className="mt-4 first:mt-0">
          {parts.map((part, i) => {
            if (i % 2 === 1) {
              return (
                <h4 
                  key={i}
                  className="text-base font-semibold mb-2"
                  style={{ color: "#ececec" }}
                >
                  {part}
                </h4>
              )
            }
            return part && (
              <p 
                key={i}
                className="text-[15px] leading-relaxed mb-3"
                style={{ color: "#ececec" }}
              >
                {part}
              </p>
            )
          })}
        </div>
      )
    }

    // Regular paragraph with possible bold text
    const formattedParagraph = paragraph.split("**").map((part, i) => {
      if (i % 2 === 1) {
        return <strong key={i} className="font-semibold">{part}</strong>
      }
      return part
    })

    return (
      <p 
        key={pIndex}
        className="text-[15px] leading-relaxed mb-3 last:mb-0"
        style={{ color: "#ececec" }}
      >
        {formattedParagraph}
      </p>
    )
  })
}

// Loading indicator component
export function AeolyzerThinkingIndicator() {
  return (
    <div className="flex items-start gap-3 mb-6">
      <AeolyzerLogoAnimated size={28} />
    </div>
  )
}
