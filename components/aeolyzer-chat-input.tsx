"use client" // -> tells Next.js this code must run in the user's browser, allowing us to use interactive tools like text inputs

import { useState, useRef, useEffect } from "react" // -> brings in core React tools for memory and side effects
import { Plus, ChevronDown, PenLine, GraduationCap, Code2, Armchair, HardDrive } from "lucide-react" // -> brings in icons from the lucide package
import { cn } from "@/lib/utils"

// -> defines the shape of the data (props) the AeolyzerChatInput block expects
interface ChatInputProps {
  onSend: (message: string) => void // -> a function it calls when the user types a message and hits enter
  isGenerating: boolean // -> true/false switch to disable the text box if the AI is still typing
  placeholder?: string // -> the grey hint text inside an empty input box
  showQuickActions?: boolean // -> true/false switch to show those little suggested action buttons below the box
}

// -> a built-in list of the suggested actions users can click (like "Write", "Learn", etc)
const quickActions = [
  { icon: PenLine, label: "Write" },
  { icon: GraduationCap, label: "Learn" },
  { icon: Code2, label: "Code" },
  { icon: Armchair, label: "Life stuff" },
  { icon: HardDrive, label: "From Drive", hasIcon: true },
]

// -> This is the main text box at the bottom where users type their messages
export function AeolyzerChatInput({ onSend, isGenerating, placeholder = "How can I help you today?", showQuickActions = false }: ChatInputProps) {
  // - `const [message, setMessage]` -> React memory to hold whatever the user has currently typed into the box
  const [message, setMessage] = useState("")
  
  // - `const [showModelDropdown, setShowModelDropdown]` -> React memory to toggle the "Models" menu open and closed
  const [showModelDropdown, setShowModelDropdown] = useState(false)
  
  // - `const [selectedModel, setSelectedModel]` -> React memory to remember which AI brain they chose (e.g., Sonnet 4.6)
  const [selectedModel, setSelectedModel] = useState("Sonnet 4.6")
  
  // - `const [extendedThinking, setExtendedThinking]` -> React memory for the "Think longer" toggle switch
  const [extendedThinking, setExtendedThinking] = useState(true)
  
  // - `const [isFocused, setIsFocused]` -> React memory to remember if the user has clicked into the text box (to highlight it)
  const [isFocused, setIsFocused] = useState(false)
  
  // - `const textareaRef` -> a sticky note attached to the actual `<textarea>` HTML element so we can measure how tall it needs to be
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  // -> useEffect runs this logic every time the user types a new letter
  // Why this matters: Text boxes usually stay one height and scroll inside. This makes the box physically grow taller as you type multiple lines.
  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "24px" // -> squeeze it down to default size first to get an accurate measurement
      // -> measure how tall the text actually is inside, but cap it at a maximum of 200px tall
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 200)}px`
    }
  }, [message])

  // -> the function that runs when they actually hit the Send button
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault() // -> stops the browser from doing its default behavior (refreshing the page)
    
    // -> `message.trim()` cuts off any accidental blank spaces at the start or end
    // -> If the message actually has letters in it AND the AI isn't currently busy typing...
    if (message.trim() && !isGenerating) {
      onSend(message.trim()) // -> fire off the message to the parent component
      setMessage("") // -> wipe the text box clean
      if (textareaRef.current) {
        textareaRef.current.style.height = "24px" // -> shrink the box back down to normal
      }
    }
  }

  // -> handles what happens when they press keys on the keyboard inside the text box
  const handleKeyDown = (e: React.KeyboardEvent) => {
    // -> if they pressed 'Enter' but did NOT hold the 'Shift' key (which would just make a new line)
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault() // -> stop the 'Enter' key from just adding a blank line in the text box
      handleSubmit(e as any) // -> pretend they clicked the send button
    }
  }

  // -> the list of AI brains available in the dropdown
  const models = [
    { name: "Opus 4.6", description: "Most capable for ambitious work", upgrade: true },
    { name: "Sonnet 4.6", description: "Most efficient for everyday tasks", selected: true },
    { name: "Haiku 4.5", description: "Fastest for quick answers" },
  ]

  return (
    <div className="w-full max-w-2xl mx-auto">
      {/* -> The form tag lets us use handleSubmit whenever a "submit" button inside it is clicked */}
      <form onSubmit={handleSubmit}>
        <div 
          className="relative rounded-xl transition-all"
          style={{ 
            backgroundColor: "#393836",
            // -> if they clicked inside the box, add a subtle glowing white ring around it
            boxShadow: isFocused ? "inset 0 0 0 1px rgba(255,255,255,0.1)" : "none"
          }}
        >
          {/* Input area - taller */}
          <div className="px-4 pt-5 pb-4">
            <textarea
              ref={textareaRef} // -> attaches our measuring sticky note here
              value={message} // -> ties what's in the box directly to our 'message' memory
              onChange={(e) => setMessage(e.target.value)} // -> updates memory whenever they type
              onKeyDown={handleKeyDown} // -> checks for the Enter key
              onFocus={() => setIsFocused(true)} // -> marks the box as active
              onBlur={() => setIsFocused(false)} // -> marks the box as inactive
              placeholder={placeholder}
              disabled={isGenerating} // -> blocks typing if the AI is generating
              className={cn(
                "w-full bg-transparent resize-none outline-none text-[15px] leading-7",
                "placeholder:font-light placeholder:tracking-wide"
              )}
              style={{ 
                color: "#ececec",
                minHeight: "32px",
                maxHeight: "200px",
                caretColor: "#ececec" // -> makes the blinking typing cursor line white
              }}
              rows={1}
            />
          </div>

          {/* Bottom toolbar */}
          <div className="flex items-center justify-between px-3 pb-3">
            {/* Left side - Add button */}
            <button
              type="button"
              className="p-2 rounded-md transition-colors hover:bg-[#2f2e2b]"
              style={{ color: "#a3a29e" }}
            >
              <Plus size={20} strokeWidth={1.5} />
            </button>

            {/* Right side - Model selector and voice */}
            <div className="flex items-center gap-2">
              {/* Model dropdown */}
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setShowModelDropdown(!showModelDropdown)} // -> flips the menu open/closed
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm transition-colors hover:bg-[#2f2e2b]"
                  style={{ color: "#ececec" }}
                >
                  <span className="font-medium">{selectedModel}</span>
                  {/* -> if the extended thinking switch is on, show a little label for it */}
                  {extendedThinking && (
                    <span style={{ color: "#a3a29e" }}>Extended</span>
                  )}
                  <ChevronDown size={14} style={{ color: "#a3a29e" }} />
                </button>

                {/* Dropdown menu */}
                {/* -> if the menu is supposed to be open, draw it on the screen */}
                {showModelDropdown && (
                  <div 
                    className="absolute bottom-full right-0 mb-2 w-72 rounded-lg shadow-xl overflow-hidden z-50"
                    style={{ backgroundColor: "#393836", border: "1px solid #4a4945" }}
                  >
                    <div className="p-2">
                      {/* -> loop over our list of AI models and draw a button for each one */}
                      {models.map((model) => (
                        <button
                          key={model.name}
                          type="button"
                          onClick={() => {
                            setSelectedModel(model.name) // -> remember what they picked
                            setShowModelDropdown(false) // -> close the menu
                          }}
                          // -> apply a grey background color automatically to the currently selected item
                          className={cn(
                            "w-full flex items-center justify-between px-3 py-2.5 rounded-md text-left transition-colors",
                            selectedModel === model.name ? "bg-[#2f2e2b]" : "hover:bg-[#2f2e2b]"
                          )}
                        >
                          <div>
                            <p className="text-sm font-medium" style={{ color: "#ececec" }}>{model.name}</p>
                            <p className="text-xs mt-0.5" style={{ color: "#a3a29e" }}>{model.description}</p>
                          </div>
                          
                          {/* -> if this model requires a paid upgrade, show an orange badge */}
                          {model.upgrade && (
                            <span className="text-xs px-2 py-0.5 rounded" style={{ backgroundColor: "#e07b53", color: "white" }}>
                              Upgrade
                            </span>
                          )}
                          
                          {/* -> if this is the currently selected model, draw a small orange checkmark */}
                          {selectedModel === model.name && !model.upgrade && (
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#e07b53" strokeWidth="2">
                              <polyline points="20,6 9,17 4,12" />
                            </svg>
                          )}
                        </button>
                      ))}
                    </div>

                    {/* Extended thinking toggle */}
                    <div className="px-3 py-3 border-t" style={{ borderColor: "#4a4945" }}>
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium" style={{ color: "#ececec" }}>Extended thinking</p>
                          <p className="text-xs" style={{ color: "#a3a29e" }}>Think longer for complex tasks</p>
                        </div>
                        {/* -> The physical toggle switch */}
                        <button
                          type="button"
                          onClick={() => setExtendedThinking(!extendedThinking)} // -> flips it on and off
                          className={cn(
                            "w-11 h-6 rounded-full transition-colors relative",
                            extendedThinking ? "bg-[#e07b53]" : "bg-[#4a4945]" // -> orange if on, grey if off
                          )}
                        >
                          <span
                            className={cn(
                              "absolute top-1 w-4 h-4 rounded-full bg-white transition-transform",
                              extendedThinking ? "right-1" : "left-1" // -> slides the white circle left or right
                            )}
                          />
                        </button>
                      </div>
                    </div>

                    {/* More models */}
                    <button
                      type="button"
                      className="w-full flex items-center justify-between px-5 py-3 border-t transition-colors hover:bg-[#2f2e2b]"
                      style={{ borderColor: "#4a4945", color: "#ececec" }}
                    >
                      <span className="text-sm">More models</span>
                      <ChevronDown size={14} className="-rotate-90" style={{ color: "#a3a29e" }} />
                    </button>
                  </div>
                )}
              </div>

              {/* Voice/Send button */}
              <button
                // -> if there's text typed, act as a submit button. Otherwise, it's just a normal button (voice input)
                type={message.trim() ? "submit" : "button"}
                className="p-2 rounded-md transition-colors hover:bg-[#2f2e2b]"
                style={{ color: "#a3a29e" }}
              >
                {message.trim() ? (
                  // -> Show arrow icon (Send) if text exists
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                  </svg>
                ) : (
                  // -> Show microphone/soundwave icon if empty
                  <div className="flex items-center gap-0.5">
                    <div className="w-0.5 h-3 rounded-full bg-[#a3a29e]" />
                    <div className="w-0.5 h-4 rounded-full bg-[#a3a29e]" />
                    <div className="w-0.5 h-2 rounded-full bg-[#a3a29e]" />
                    <div className="w-0.5 h-5 rounded-full bg-[#a3a29e]" />
                    <div className="w-0.5 h-3 rounded-full bg-[#a3a29e]" />
                  </div>
                )}
              </button>
            </div>
          </div>
        </div>
      </form>

      {/* Quick action pills - shown on welcome screen */}
      {/* -> if the prop told us to show them, generate all the suggestion buttons */}
      {showQuickActions && (
        <div className="flex flex-wrap justify-center gap-2 mt-4">
          {quickActions.map((action, index) => (
            <button
              key={index}
              className="flex items-center gap-2 px-4 py-2 rounded-lg text-sm transition-colors hover:bg-[#252422]"
              style={{ 
                backgroundColor: "#2b2a27", 
                color: "#ececec",
                border: "1px solid #4a4945"
              }}
            >
              <action.icon size={16} strokeWidth={1.5} />
              <span>{action.label}</span>
            </button>
          ))}
        </div>
      )}


    </div>
  )
}
