"use client"

import { useState, useEffect } from "react"
import { AeolyzerLogo } from "./aeolyzer-logo"
import { AeolyzerChatInput } from "./aeolyzer-chat-input"

interface WelcomeProps {
  userName?: string
  onSend: (message: string) => void
  isGenerating: boolean
}

export function AeolyzerWelcome({ userName = "Muhammad", onSend, isGenerating }: WelcomeProps) {
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
    <div className="flex flex-col items-center justify-center flex-1 px-4">
      {/* Free plan badge */}
      <div className="mb-6">
        <span 
          className="text-sm px-3 py-1.5 rounded-full bg-welcome-bg text-welcome-text"
        >
          Free plan · <span className="text-accent cursor-pointer hover:underline">Upgrade</span>
        </span>
      </div>

      {/* Greeting */}
      <div className="flex items-center gap-4 mb-4" style={{ minHeight: "64px" }}>
        <AeolyzerLogo size={52} />
        {mounted && greeting && (
          <h1 
            className="font-light"
            style={{ 
              color: "#b8977e", 
              fontFamily: "var(--font-display), 'Rokkitt', Georgia, serif",
              fontSize: "3rem",
              lineHeight: "1.1"
            }}
          >
            {greeting}, {userName}
          </h1>
        )}
      </div>

      {/* Centered input with quick actions */}
      <div className="w-full max-w-3xl">
        <AeolyzerChatInput 
          onSend={onSend}
          isGenerating={isGenerating}
          placeholder="How can I help you today?"
          showQuickActions={true}
        />
      </div>
    </div>
  )
}
