"use client" // -> tells Next.js this code must run in the user's browser, allowing us to use interactive features like useState and useEffect

import { useState, useEffect } from "react" // -> brings in core React tools for memory (useState) and side-effects (useEffect)
import { AeolyzerLogo } from "./aeolyzer-logo"
import { AeolyzerChatInput } from "./aeolyzer-chat-input"

// -> defines the shape of the data (props) the AeolyzerWelcome block expects to receive
// - `userName` -> optional text holding the user's name
// - `onSend` -> a function it can call when the user types a message and hits enter
// - `isGenerating` -> true/false switch to tell if the AI is currently thinking/typing
interface WelcomeProps {
  userName?: string
  onSend: (message: string) => void
  isGenerating: boolean
}

// -> This is the main Welcome screen component shown when a new chat starts
export function AeolyzerWelcome({ userName = "Muhammad", onSend, isGenerating }: WelcomeProps) {
  // - `const [greeting, setGreeting]` -> sets up a piece of memory to store our "Good morning/afternoon/evening" message
  const [greeting, setGreeting] = useState<string | null>(null)
  
  // - `const [mounted, setMounted]` -> sets up memory to track if the component has finished loading onto the page
  const [mounted, setMounted] = useState(false)

  // -> useEffect runs a piece of code after the component first appears on screen
  // We calculate the time here instead of instantly to prevent a "hydration mismatch" (where the server's timezone doesn't match the user's browser)
  useEffect(() => {
    setMounted(true) // -> records that we are safely loaded in the browser
    const hour = new Date().getHours() // -> checks the current hour of the day (0-23)
    
    // -> decides which greeting to use based on the clock
    if (hour < 12) {
      setGreeting("Good morning")
    } else if (hour < 18) {
      setGreeting("Good afternoon")
    } else {
      setGreeting("Good evening")
    }
  }, []) // -> the empty track [] means this code only runs exactly once when the component appears

  return (
    // -> A container that centers everything in the middle of the screen
    <div className="flex flex-col items-center justify-center flex-1 px-4">
      {/* Free plan badge */}
      <div className="mb-6">
        <span 
          className="text-sm px-3 py-1.5 rounded-full"
          style={{ backgroundColor: "#343330", color: "#8b8b87" }}
        >
          Free plan · <span style={{ color: "#e07b53" }} className="cursor-pointer hover:underline">Upgrade</span>
        </span>
      </div>

      {/* Greeting Section */}
      <div className="flex items-center gap-4 mb-4" style={{ minHeight: "64px" }}>
        {/* -> shows the sunburst logo we commented in aeolyzer-logo */}
        <AeolyzerLogo size={52} />
        
        {/* -> Wait until 'mounted' and 'greeting' are ready before showing the text */}
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

      {/* Centered chat input with quick action buttons */}
      <div className="w-full max-w-3xl">
        <AeolyzerChatInput 
          onSend={onSend}
          isGenerating={isGenerating}
          placeholder="How can I help you today?"
          showQuickActions={true} // -> turns on the suggested actions (like "Help me write...")
        />
      </div>
    </div>
  )
}
