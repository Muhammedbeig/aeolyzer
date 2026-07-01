"use client"

import { useState, useRef, useEffect } from "react"
import { PenLine, Search, FileText, Activity, ListTodo, LineChart, Sparkles, Eye, LayoutTemplate, Users, Target, Plus } from "lucide-react"
import { cn } from "@/lib/utils"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { AttachmentPreviewList } from "./attachment-preview-list"
import { ContentTypeSelector } from "./content-type-selector"
import type { ContentType } from "./types"

interface ChatInputProps {
  onSend: (message: string, files?: File[], contentType?: ContentType) => void
  isGenerating: boolean
  placeholder?: string
  showQuickActions?: boolean
  showContentOptions?: boolean
  contentType?: ContentType
  onContentTypeChange?: (contentType: ContentType) => void
}

const MAX_ATTACHMENTS = 5
const MAX_ATTACHMENT_BYTES = 10 << 20
const MAX_TOTAL_ATTACHMENT_BYTES = 20 << 20

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

export function AeolyzerChatInput({
  onSend,
  isGenerating,
  placeholder = "How can I help you today?",
  showQuickActions = false,
  showContentOptions = false,
  contentType = "article",
  onContentTypeChange,
}: ChatInputProps) {
  const [message, setMessage] = useState("")
  const [selectedFiles, setSelectedFiles] = useState<File[]>([])
  const [attachmentError, setAttachmentError] = useState<string>()
  const suggestions = PROMPT_SUGGESTIONS.slice(0, 4)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "24px"
      textareaRef.current.style.height = `${Math.min(textareaRef.current.scrollHeight, 200)}px`
    }
  }, [message])

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if ((message.trim() || selectedFiles.length > 0) && !isGenerating) {
      onSend(message.trim(), selectedFiles, showContentOptions ? contentType : undefined)
      setMessage("")
      setSelectedFiles([])
      setAttachmentError(undefined)
      if (fileInputRef.current) {
        fileInputRef.current.value = ""
      }
      if (textareaRef.current) {
        textareaRef.current.style.height = "24px"
      }
    }
  }

  const handleFilesSelected = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(event.target.files ?? [])
    const available = MAX_ATTACHMENTS - selectedFiles.length
    const accepted = files.slice(0, Math.max(available, 0))
    const candidate = [...selectedFiles, ...accepted]
    const oversized = accepted.some((file) => file.size > MAX_ATTACHMENT_BYTES)
    const totalBytes = candidate.reduce((total, file) => total + file.size, 0)
    if (oversized) {
      setAttachmentError("Each attachment must be 10 MB or smaller.")
      event.currentTarget.value = ""
      return
    }
    if (totalBytes > MAX_TOTAL_ATTACHMENT_BYTES) {
      setAttachmentError("Attachments can total no more than 20 MB.")
      event.currentTarget.value = ""
      return
    }
    setSelectedFiles(candidate)
    setAttachmentError(
      files.length > available ? "You can attach up to 5 files." : undefined,
    )
    event.currentTarget.value = ""
  }

  const removeFile = (index: number) => {
    setSelectedFiles((current) =>
      current.filter((_, fileIndex) => fileIndex !== index),
    )
    setAttachmentError(undefined)
  }

  const canSend = Boolean(message.trim() || selectedFiles.length > 0)

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      handleSubmit(e)
    }
  }

  return (
    <div
      className="w-full max-w-3xl mx-auto flex flex-col-reverse sm:flex-col"
      data-testid="chat-input"
    >
      <form onSubmit={handleSubmit} className="mt-2 sm:mt-0">
        <div className="relative rounded-[26px] bg-white shadow-sm hover:shadow-md transition-all duration-300 focus-within:outline focus-within:outline-1 focus-within:outline-accent/40 dark:bg-card dark:shadow-black/20">
          <AttachmentPreviewList
            files={selectedFiles}
            disabled={isGenerating}
            onRemove={removeFile}
          />

          <div className="relative w-full">
            <textarea
              ref={textareaRef}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              onKeyDown={handleKeyDown}
              placeholder={placeholder}
              disabled={isGenerating}
              className={cn(
                "w-full px-5 py-4 bg-transparent rounded-t-[26px] outline-none",
                "text-[15px] leading-relaxed text-plum-800 dark:text-foreground caret-plum-900 dark:caret-foreground",
                "placeholder:text-[#9CA3AF] dark:placeholder:text-muted-foreground placeholder:font-normal",
                "min-h-[60px] max-h-[200px] resize-none"
              )}
              rows={1}
            />
          </div>

          <div className="flex items-center justify-between px-2.5 pb-2.5 pt-0">
            <div className="flex items-center gap-1">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button
                    type="button"
                    disabled={isGenerating}
                    className="w-8 h-8 rounded-lg bg-transparent text-plum-500 dark:text-muted-foreground hover:text-plum-700 dark:hover:text-foreground hover:bg-sand-50 dark:hover:bg-muted flex items-center justify-center transition-colors cursor-pointer disabled:pointer-events-none disabled:opacity-50"
                  >
                    <Plus className="w-5 h-5" strokeWidth={2} />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start" className="w-48 bg-white dark:bg-card border border-black/10 dark:border-white/10 rounded-xl shadow-lg p-1">
                  <DropdownMenuItem 
                    onClick={() => fileInputRef.current?.click()}
                    className="flex items-center gap-2 px-3 py-2 cursor-pointer rounded-lg hover:bg-sand-50 dark:hover:bg-muted focus:bg-sand-50 dark:focus:bg-muted focus:text-plum-900 dark:focus:text-foreground outline-none text-[13px] font-medium transition-colors group"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4 text-muted-foreground group-focus:text-plum-900 dark:group-focus:text-foreground transition-colors">
                      <path d="m21.44 11.05-9.19 9.19a6 6 0 0 1-8.49-8.49l9.19-9.19a4 4 0 0 1 5.66 5.66l-9.2 9.19a2 2 0 0 1-2.83-2.83l8.49-8.48"></path>
                    </svg>
                    Add Attachment
                  </DropdownMenuItem>
                  <DropdownMenuItem 
                    className="flex items-center gap-2 px-3 py-2 cursor-pointer rounded-lg hover:bg-sand-50 dark:hover:bg-muted focus:bg-sand-50 dark:focus:bg-muted focus:text-plum-900 dark:focus:text-foreground outline-none text-[13px] font-medium transition-colors group"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4 text-muted-foreground group-focus:text-plum-900 dark:group-focus:text-foreground transition-colors">
                      <path d="M15 6v12a3 3 0 1 0 3-3H6a3 3 0 1 0 3 3V6a3 3 0 1 0-3 3h12a3 3 0 1 0-3-3"></path>
                    </svg>
                    Shortcuts
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
              <input
                ref={fileInputRef}
                type="file"
                multiple
                accept=".pdf,.png,.jpg,.jpeg,.gif,.txt,.md,.markdown,.csv,.json,.html,.htm,.css,.js,.ts,.tsx,.jsx,.go,.py,.java,.rs,.xml,.yaml,.yml"
                onChange={handleFilesSelected}
                className="sr-only"
                aria-label="Choose attachments"
              />
            </div>

            <div className="flex items-center gap-1.5 sm:gap-2">
              <button
                type="submit"
                disabled={!canSend || isGenerating}
                className={cn(
                  "w-8 h-8 sm:w-9 sm:h-9 rounded-full flex items-center justify-center transition-all",
                  canSend
                    ? "bg-accent dark:bg-accent text-white hover:bg-accent/90 cursor-pointer shadow-sm hover:shadow-md"
                    : "bg-sand-200 dark:bg-muted text-plum-400 dark:text-muted-foreground/50 cursor-default"
                )}
                aria-label="Send message"
              >
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="w-4 h-4 sm:w-5 sm:h-5 transition-transform duration-200 ease-out">
                  <path d="m5 12 7-7 7 7"></path>
                  <path d="M12 19V5"></path>
                </svg>
              </button>
            </div>
          </div>
          {attachmentError && (
            <p className="px-4 pb-3 text-xs text-destructive" role="alert">
              {attachmentError}
            </p>
          )}
        </div>
      </form>

      {showContentOptions && (
        <ContentTypeSelector
          value={contentType}
          disabled={isGenerating}
          onChange={(value) => onContentTypeChange?.(value)}
        />
      )}

      {showQuickActions && (
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3 mb-4 sm:mb-0 sm:mt-4 px-1 sm:px-0">
          {suggestions.map((action, index) => (
            <button
              key={index}
              onClick={() => onSend(action.prompt)}
              className="flex flex-col items-start justify-between p-3 sm:p-4 h-24 sm:h-28 rounded-xl bg-white dark:bg-card shadow-sm hover:shadow-md transition-[box-shadow] duration-200 text-left group cursor-pointer"
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
