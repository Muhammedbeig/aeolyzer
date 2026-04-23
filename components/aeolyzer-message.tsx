"use client" // -> tells Next.js this code must run in the user's browser, allowing us to use interactive tools like useState

import { useState, useEffect } from "react" // -> brings in React's memory (useState) and side-effect tool (useEffect)
import { Copy, ThumbsUp, ThumbsDown, RotateCcw } from "lucide-react" // -> brings in icons for the action buttons
import { AeolyzerLogoAnimated } from "./aeolyzer-logo"
import { cn } from "@/lib/utils"

// -> defines what a chat Message object looks like
// - `id` -> the unique label for this exact message
// - `role` -> who sent it (either the human "user" or the AI "assistant")
// - `content` -> the actual text of the message
// - `isStreaming` -> optional true/false switch indicating if the AI is currently "typing" this message out
interface Message {
  id: string
  role: "user" | "assistant"
  content: string
  isStreaming?: boolean
}

// -> defines the shape of the data (props) the AeolyzerMessage block expects to receive
interface AeolyzerMessageProps {
  message: Message
}

// -> This builds the visual bubble for a single message in the chat
export function AeolyzerMessage({ message }: AeolyzerMessageProps) {
  // - `const [displayedContent, setDisplayedContent]` -> memory to track the text currently visible on screen (used to create the typing effect)
  const [displayedContent, setDisplayedContent] = useState("")
  
  // - `const [isComplete, setIsComplete]` -> memory to track whether the AI has finished typing the whole message
  const [isComplete, setIsComplete] = useState(false)

  // -> runs this logic whenever the message changes to handle the fake typing effect
  useEffect(() => {
    // -> If this is an AI message and it wants to fake a "typing" animation...
    if (message.role === "assistant" && message.isStreaming) {
      setDisplayedContent("") // -> start with a blank screen
      setIsComplete(false) // -> mark that we aren't done yet
      let currentIndex = 0
      const content = message.content
      
      /*
      Let's break this down very slowly, word by word 👇
      
      🧩 Line 1:
      const interval = setInterval(() => { ... }, 15)
      
      🔸 `setInterval`
      - A JavaScript tool that runs a block of code over and over again
      👉 Think: "A metronome tapping every 15 milliseconds, telling the browser to add more letters."
      
      🧩 Line 2:
      const chunkSize = Math.floor(Math.random() * 3) + 1
      
      🔸 `Math.random() * 3`
      - Generates a random decimal between 0 and 2.99
      🔸 `Math.floor(...) + 1`
      - Chops off the decimal and adds 1 (resulting in 1, 2, or 3)
      👉 Why? Real humans don't type at a perfectly even robotic speed. This makes the text appear in slightly random chunks of letters, mimicking a real Claude response.
      
      🧩 Line 3:
      const nextChunk = content.slice(currentIndex, currentIndex + chunkSize)
      
      🔸 `.slice(...)`
      - Cuts out a piece of the message string
      - Starts at our current spot in the sentence (`currentIndex`)
      - Ends a few letters later (`currentIndex + chunkSize`)
      
      👉 Example (if `chunkSize = 2` and `content = "Hello world"`):
      | Step | currentIndex | slice() | nextChunk | New Screen Text |
      |---|---|---|---|---|
      | 1 | 0 | slice(0, 2) | "He" | "He" |
      | 2 | 2 | slice(2, 4) | "ll" | "Hell" |
      | 3 | 4 | slice(4, 7) | "o w" (chunk=3) | "Hello w" |
      
      🧩 Final Result
      The text slowly bleeds onto the screen letter by letter, giving a smooth typing effect.
      */
      const interval = setInterval(() => {
        // -> if we haven't typed the whole message yet...
        if (currentIndex < content.length) {
          const chunkSize = Math.floor(Math.random() * 3) + 1
          const nextChunk = content.slice(currentIndex, currentIndex + chunkSize)
          
          // -> add them to what's currently on the screen
          setDisplayedContent(prev => prev + nextChunk)
          currentIndex += chunkSize // -> update our place in the book
        } else {
          // -> Once we reach the end of the message:
          clearInterval(interval) // -> stop the repeating timer
          setIsComplete(true) // -> mark the message as fully finished
        }
      }, 15) // Fast typing speed like Claude, runs every 15 milliseconds

      // -> Cleanup function: if the component shuts down early, stop the timer so it doesn't run forever in the background
      return () => clearInterval(interval)
    } else {
      // -> If it's the human's message, or a completely finished previous message, just show it all instantly
      setDisplayedContent(message.content)
      setIsComplete(true)
    }
  }, [message.content, message.isStreaming, message.role]) // -> only rerun this if the words themselves change

  // -> How to draw a message sent by the human
  if (message.role === "user") {
    return (
      <div className="flex justify-end mb-6 animate-fade-in-up">
        {/* -> A rounded dark-gray box pushed to the right side of the screen */}
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

  // -> How to draw a message sent by the AI
  return (
    <div className="mb-6 animate-fade-in-up">
      {/* Thinking indicator when starting */}
      {/* -> If it's marked as streaming, but we haven't typed any letters yet... */}
      {message.isStreaming && displayedContent.length === 0 && (
        <div className="flex items-center gap-3 mb-4">
          <AeolyzerLogoAnimated size={28} /> {/* -> ...show the spinning sunburst */}
          <span className="text-sm" style={{ color: "#a3a29e" }}>Thinking...</span>
        </div>
      )}

      {/* Message content */}
      <div className="max-w-none">
        <div className="prose prose-invert max-w-none">
          {/* -> takes the raw text and makes parts of it bold or turns them into sections */}
          {renderFormattedContent(displayedContent)}
        </div>
        
        {/* Typing cursor */}
        {/* -> If we are currently typing... */}
        {message.isStreaming && !isComplete && displayedContent.length > 0 && (
          // -> ...show a little blinking orange block at the end
          <span className="inline-block w-2 h-5 ml-0.5 bg-[#e07b53] animate-pulse" />
        )}
      </div>

      {/* Action buttons - show after complete */}
      {/* -> Only show the Copy and Thumbs Up buttons if the AI has finished typing */}
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

// -> This helper function chops up the message text and styles it so it isn't just one giant wall of raw text
function renderFormattedContent(content: string) {
  // Split content by double newlines for paragraphs
  // -> slices the huge block of text into separate chunks wherever the AI pressed 'Enter' twice
  const paragraphs = content.split(/\n\n/)
  
  // -> loops over every single paragraph chunk
  return paragraphs.map((paragraph, pIndex) => {
    // Check for headers
    // -> If the paragraph starts with ** and ends with **
    if (paragraph.startsWith("**") && paragraph.endsWith("**")) {
      const headerText = paragraph.slice(2, -2) // -> cut the ** symbols off
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
    // -> If the paragraph contains ** and a colon
    if (paragraph.includes("**") && paragraph.includes(":")) {
      const parts = paragraph.split("**")
      return (
        <div key={pIndex} className="mt-4 first:mt-0">
          {parts.map((part, i) => {
            // -> Every other separated piece gets marked as a bold sub-header
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
    // -> slices the text wherever there is a ** symbol, turning alternating chunks strong and bold
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
// -> The mini visual that shows up when the AI first starts thinking
export function AeolyzerThinkingIndicator() {
  return (
    <div className="flex items-start gap-3 mb-6">
      <AeolyzerLogoAnimated size={28} />
    </div>
  )
}
