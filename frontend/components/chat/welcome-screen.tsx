"use client"

import { useState, useEffect } from "react"
import { AeolyzerLogo } from "@/components/ui/aeolyzer-logo"
import { AeolyzerChatInput } from "./chat-input"
import { cn } from "@/lib/utils"

interface WelcomeProps {
  userName?: string
  onSend: (message: string) => void
  isGenerating: boolean
  title?: string
  placeholder?: string
  showContentOptions?: boolean
}

export function AeolyzerWelcome({ 
  userName = "Muhammad", 
  onSend, 
  isGenerating,
  title,
  placeholder = "How can I help you today?",
  showContentOptions
}: WelcomeProps) {
  const [greeting, setGreeting] = useState<string | null>(null)
  const [mounted, setMounted] = useState(false)

  // Set greeting on client side only to avoid hydration mismatch
  useEffect(() => {
    setMounted(true)
    const hour = new Date().getHours()
    if (hour < 12) {
      setGreeting("Good morning")
    } else if (hour < 18) {
      setGreeting("Good afternoon")
    } else {
      setGreeting("Good evening")
    }
  }, [])

  return (
    <div className="flex flex-col items-center sm:justify-center flex-1 px-4 pb-4 sm:pb-0">
      
      {/* Centered content wrapper for mobile */}
      <div className="flex-1 sm:flex-none flex flex-col items-center justify-center w-full">
        {/* Greeting */}
        <div className="flex flex-col sm:flex-row items-center gap-4 mb-4" style={{ minHeight: "64px" }}>
          <AeolyzerLogo size={52} />
          {mounted && (
            <h1 
              className="text-center font-outfit tracking-tight font-light text-foreground/80 text-5xl leading-tight"
            >
              {title || (greeting ? `${greeting}, ${userName}` : "")}
            </h1>
          )}
        </div>
      </div>

      {/* Input pinned to bottom on mobile, centered on desktop */}
      <div className="w-full max-w-3xl mt-auto sm:mt-0">
        <AeolyzerChatInput 
          onSend={onSend}
          isGenerating={isGenerating}
          placeholder={placeholder}
          showQuickActions={!showContentOptions}
          showContentOptions={showContentOptions}
        />
      </div>
    </div>
  )
}
