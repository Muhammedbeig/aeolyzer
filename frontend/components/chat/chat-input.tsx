"use client"

import { useState, useRef, useEffect } from "react"
import { Plus, PenLine, Search, FileText, Activity, ListTodo, LineChart, Sparkles, Eye, LayoutTemplate, Users, Target } from "lucide-react"
import { cn } from "@/lib/utils"

interface ChatInputProps {
  onSend: (message: string) => void
  isGenerating: boolean
  placeholder?: string
  showQuickActions?: boolean
}

const PROMPT_SUGGESTIONS = [
  { icon: Search, label: "Optimize Meta tags", prompt: "Review the meta titles and descriptions for my key pages. Suggest optimized versions that improve click-through rates (CTR) and include relevant keywords without keyword stuffing." },
  { icon: FileText, label: "Generate a llm.txt", prompt: "Generate an llms.txt file for my website to help AI agents understand my content better. Include key pages, sitemap location, and a concise description of my site's purpose." },
  { icon: Activity, label: "Review Technical SEO", prompt: "Perform a technical SEO review. Check for crawl errors, broken links, duplicate content, canonical tag issues, and mobile-friendliness problems." },
  { icon: ListTodo, label: "Make me a Task List", prompt: "Generate a prioritized SEO task list for my website. Focus on high-impact, low-effort actions I can take this week to improve my search rankings. Categorize them by Technical, Content, and Authority." },
  { icon: LineChart, label: "Summarize recent Performance", prompt: "Summarize my website's search performance over the last 30 days. Highlight significant changes in traffic, rankings, and impressions. What went well and what needs attention?" },
  { icon: Sparkles, label: "Find Keyword opportunities", prompt: "Find underutilized keyword opportunities for my niche. Look for long-tail keywords with decent search volume and low competition that I can target with new content." },
  { icon: Eye, label: "Quick Visibility Check", prompt: "Perform a comprehensive visibility check for my domain. Analyze my ranking for top keywords, identify visibility trends, and summarize my overall presence in search results." },
  { icon: LayoutTemplate, label: "Audit my Homepage", prompt: "Conduct a detailed audit of my homepage. Check for on-page SEO issues, technical errors, user experience friction, and conversion optimization opportunities. Provide actionable recommendations." },
  { icon: Users, label: "Analyze competitors", prompt: "Analyze my top 3 competitors. Compare their search visibility, top-performing content, and backlink profiles to mine. Highlight their strengths and my opportunities to outperform them." },
  { icon: PenLine, label: "Draft a Blog Post", prompt: "Draft a high-quality, SEO-optimized blog post about a trending topic in my industry. Include a catchy title, headers, and a structure that targets user intent." },
  { icon: Target, label: "Find content gaps", prompt: "Identify content gaps on my website compared to my top competitors. What topics are they covering that I am missing? Suggest 5 new article ideas to fill these gaps." },
]

export function AeolyzerChatInput({ onSend, isGenerating, placeholder = "How can I help you today?", showQuickActions = false }: ChatInputProps) {
  const [message, setMessage] = useState("")
  const [suggestions, setSuggestions] = useState<typeof PROMPT_SUGGESTIONS>([])
  
  const [isFocused, setIsFocused] = useState(false)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    // Select 4 random suggestions on mount
    const shuffled = [...PROMPT_SUGGESTIONS].sort(() => 0.5 - Math.random())
    setSuggestions(shuffled.slice(0, 4))
  }, [])

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

  return (
    <div className="w-full max-w-3xl mx-auto flex flex-col-reverse sm:flex-col">
      <form onSubmit={handleSubmit} className="mt-2 sm:mt-0">
        <div 
          className={cn(
            "relative rounded-xl transition-all bg-card border",
            isFocused ? "border-accent shadow-sm" : "border-muted-foreground/20 shadow-sm"
          )}
        >
          {/* Input area - taller */}
          <div className="px-4 pt-4 sm:pt-5 pb-3 sm:pb-4">
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
                "placeholder:font-light placeholder:tracking-wide min-h-[24px] sm:min-h-[32px] max-h-[200px]"
              )}
              rows={1}
            />
          </div>

          {/* Bottom toolbar */}
          <div className="flex items-center justify-between px-2 sm:px-3 pb-2 sm:pb-3">
            {/* Left side - Add button */}
            <button
              type="button"
              className="p-2 rounded-md transition-colors hover:bg-muted text-muted-foreground"
              aria-label="Add attachment"
            >
              <Plus size={20} strokeWidth={1.5} />
            </button>

            {/* Right side - Send button */}
            <div className="flex items-center gap-2">
              <button
                type={message.trim() ? "submit" : "button"}
                className={cn(
                  "p-2 rounded-md transition-colors",
                  message.trim() 
                    ? "bg-foreground text-background hover:bg-foreground/90" 
                    : "hover:bg-muted text-muted-foreground"
                )}
                aria-label={message.trim() ? "Send message" : "Voice input"}
              >
                {message.trim() ? (
                  <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z" />
                  </svg>
                ) : (
                  <div className="flex items-center gap-0.5">
                    <div className="w-0.5 h-3 rounded-full bg-current" />
                    <div className="w-0.5 h-4 rounded-full bg-current" />
                    <div className="w-0.5 h-2 rounded-full bg-current" />
                    <div className="w-0.5 h-5 rounded-full bg-current" />
                    <div className="w-0.5 h-3 rounded-full bg-current" />
                  </div>
                )}
              </button>
            </div>
          </div>
        </div>
      </form>

      {/* Quick action pills */}
      {showQuickActions && (
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3 mb-4 sm:mb-0 sm:mt-4 px-1 sm:px-0">
          {suggestions.map((action, index) => (
            <button
              key={index}
              onClick={() => onSend(action.prompt)}
              className="flex flex-col items-start justify-between p-3 sm:p-4 h-24 sm:h-28 rounded-xl bg-card border border-border/50 hover:border-muted-foreground/30 transition-all text-left group cursor-pointer shadow-sm hover:shadow-md"
              title={action.prompt}
            >
              <action.icon size={20} strokeWidth={1.5} className="text-muted-foreground group-hover:text-foreground transition-colors flex-shrink-0" />
              <span className="text-[12.5px] sm:text-[13.5px] font-medium leading-tight text-foreground/90 group-hover:text-foreground">{action.label}</span>
            </button>
          ))}
        </div>
      )}
    </div>
  )
}
