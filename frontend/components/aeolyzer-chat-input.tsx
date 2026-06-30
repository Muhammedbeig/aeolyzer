"use client"

import { useState, useRef, useEffect } from "react"
import { Plus, ChevronDown, PenLine, GraduationCap, Code2, Armchair, HardDrive } from "lucide-react"
import { cn } from "@/lib/utils"

interface ChatInputProps {
  onSend: (message: string) => void
  isGenerating: boolean
  placeholder?: string
  showQuickActions?: boolean
}

const quickActions = [
  { icon: PenLine, label: "Write" },
  { icon: GraduationCap, label: "Learn" },
  { icon: Code2, label: "Code" },
  { icon: Armchair, label: "Life stuff" },
  { icon: HardDrive, label: "From Drive", hasIcon: true },
]

export function AeolyzerChatInput({ onSend, isGenerating, placeholder = "How can I help you today?", showQuickActions = false }: ChatInputProps) {
  const [message, setMessage] = useState("")
  const [showModelDropdown, setShowModelDropdown] = useState(false)
  const [selectedModel, setSelectedModel] = useState("Sonnet 4.6")
  const [extendedThinking, setExtendedThinking] = useState(true)
  const [isFocused, setIsFocused] = useState(false)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "24px"
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 200)}px`
    }
  }, [message])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (message.trim() && !isGenerating) {
      onSend(message.trim())
      setMessage("")
      if (textareaRef.current) {
        textareaRef.current.style.height = "24px"
      }
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSubmit(e)
    }
  }

  const models = [
    { name: "Opus 4.6", description: "Most capable for ambitious work", upgrade: true },
    { name: "Sonnet 4.6", description: "Most efficient for everyday tasks", selected: true },
    { name: "Haiku 4.5", description: "Fastest for quick answers" },
  ]

  return (
    <div className="w-full max-w-2xl mx-auto">
      <form onSubmit={handleSubmit}>
        <div 
          className={cn(
            "relative rounded-xl transition-all bg-card border",
            isFocused ? "border-accent shadow-sm" : "border-muted-foreground/20 shadow-sm"
          )}
        >
          {/* Input area - taller */}
          <div className="px-4 pt-5 pb-4">
            <textarea
              ref={textareaRef}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              onKeyDown={handleKeyDown}
              onFocus={() => setIsFocused(true)}
              onBlur={() => setIsFocused(false)}
              placeholder={placeholder}
              disabled={isGenerating}
              className={cn(
                "w-full bg-transparent resize-none outline-none text-[15px] leading-7 text-foreground caret-foreground",
                "placeholder:font-light placeholder:tracking-wide min-h-[32px] max-h-[200px]"
              )}
              rows={1}
            />
          </div>

          {/* Bottom toolbar */}
          <div className="flex items-center justify-between px-3 pb-3">
            {/* Left side - Add button */}
            <button
              type="button"
              className="p-2 rounded-md transition-colors hover:bg-muted text-muted-foreground"
            >
              <Plus size={20} strokeWidth={1.5} />
            </button>

            {/* Right side - Model selector and voice */}
            <div className="flex items-center gap-2">
              {/* Model dropdown */}
              <div className="relative">
                <button
                  type="button"
                  onClick={() => setShowModelDropdown(!showModelDropdown)}
                  className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-sm transition-colors hover:bg-muted text-foreground"
                >
                  <span className="font-medium">{selectedModel}</span>
                  {extendedThinking && (
                    <span className="text-muted-foreground">Extended</span>
                  )}
                  <ChevronDown size={14} className="text-muted-foreground" />
                </button>

                {/* Dropdown menu */}
                {showModelDropdown && (
                  <div 
                    className="absolute bottom-full left-0 mb-2 w-[320px] rounded-xl shadow-xl overflow-hidden z-50 bg-card border border-border"
                  >
                    <div className="p-2">
                      {models.map((model) => (
                        <button
                          key={model.name}
                          type="button"
                          onClick={() => {
                            setSelectedModel(model.name)
                            setShowModelDropdown(false)
                          }}
                          className={cn(
                            "w-full flex items-center justify-between px-3 py-2.5 rounded-md text-left transition-colors",
                            selectedModel === model.name ? "bg-muted" : "hover:bg-muted"
                          )}
                        >
                          <div>
                            <p className="text-sm font-medium text-foreground">{model.name}</p>
                            <p className="text-xs mt-0.5 text-muted-foreground">{model.description}</p>
                          </div>
                          {model.upgrade && (
                            <span className="text-xs px-2 py-0.5 rounded bg-accent text-white">
                              Upgrade
                            </span>
                          )}
                          {selectedModel === model.name && !model.upgrade && (
                            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                              <polyline points="20,6 9,17 4,12" />
                            </svg>
                          )}
                        </button>
                      ))}
                    </div>

                    {/* Extended thinking toggle */}
                    <div className="px-3 py-3 border-t border-border">
                      <div className="flex items-center justify-between">
                        <div>
                          <p className="text-sm font-medium text-foreground">Extended thinking</p>
                          <p className="text-xs text-muted-foreground">Think longer for complex tasks</p>
                        </div>
                        <button
                          type="button"
                          onClick={() => setExtendedThinking(!extendedThinking)}
                          className={cn(
                            "w-11 h-6 rounded-full transition-colors relative",
                            extendedThinking ? "bg-primary" : "bg-muted"
                          )}
                        >
                          <span
                            className={cn(
                              "absolute top-1 w-4 h-4 rounded-full bg-background transition-transform",
                              extendedThinking ? "right-1" : "left-1"
                            )}
                          />
                        </button>
                      </div>
                    </div>

                    {/* More models */}
                    <button
                      type="button"
                      className="w-full flex items-center justify-between px-5 py-3 border-t transition-colors hover:bg-muted border-border text-foreground"
                    >
                      <span className="text-sm">More models</span>
                      <ChevronDown size={14} className="-rotate-90 text-muted-foreground" />
                    </button>
                  </div>
                )}
              </div>

              {/* Voice/Send button */}
              <button
                type={message.trim() ? "submit" : "button"}
                className="p-2 rounded-md transition-colors hover:bg-muted text-muted-foreground"
              >
                {message.trim() ? (
                  <svg width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                  </svg>
                ) : (
                  <div className="flex items-center gap-0.5">
                    <div className="w-0.5 h-3 rounded-full bg-[var(--muted-foreground)]" />
                    <div className="w-0.5 h-4 rounded-full bg-[var(--muted-foreground)]" />
                    <div className="w-0.5 h-2 rounded-full bg-[var(--muted-foreground)]" />
                    <div className="w-0.5 h-5 rounded-full bg-[var(--muted-foreground)]" />
                    <div className="w-0.5 h-3 rounded-full bg-[var(--muted-foreground)]" />
                  </div>
                )}
              </button>
            </div>
          </div>
        </div>
      </form>

      {/* Quick action pills - shown on welcome screen */}
      {showQuickActions && (
        <div className="flex flex-wrap justify-center gap-2 mt-4">
          {quickActions.map((action, index) => (
            <button
              key={index}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-lg text-sm transition-colors hover:bg-muted border border-border bg-card text-foreground"
              )}
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
