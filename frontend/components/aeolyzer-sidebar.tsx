"use client"

import { useState } from "react"
import { 
  Plus, 
  Search, 
  Settings2, 
  MessageSquare, 
  FolderOpen, 
  Grid2X2, 
  Code,
  ChevronDown,
  Download,
  Settings,
  Globe,
  HelpCircle,
  Sparkles,
  Gift,
  BookOpen,
  LogOut,
  MoreHorizontal
} from "lucide-react"
import { cn } from "@/lib/utils"

interface SidebarProps {
  isOpen: boolean
  onToggle: () => void
  onNewChat: () => void
  currentChatTitle?: string
  onOpenSettings: () => void
}

const recentChats = [
  "Your greatest strengths",
  "Downloaded page color mismatch",
  "Answer engine optimization key...",
  "Futbol fibre audience analysis",
  "Football websites that went viral...",
  "Growing soccer live score websi...",
  "SEO-optimized domain names fo...",
  "Fixing mismatched sports icons",
  "README file rewrite",
  "Social media style ad post with e...",
  "Unified AI generation API for ima...",
  "AI-powered ad generation with t...",
  "Fixing modal styling without bre...",
  "Email verification requirement f...",
  "README File Correction",
  "DuckDuckGo Search Views Rev...",
]

// Custom sidebar toggle icon matching Claude's design
function SidebarToggleIcon({ isOpen }: { isOpen: boolean }) {
  return (
    <svg 
      width="18" 
      height="18" 
      viewBox="0 0 24 24" 
      fill="none" 
      stroke="currentColor" 
      strokeWidth="2" 
      strokeLinecap="round" 
      strokeLinejoin="round"
    >
      {isOpen ? (
        <>
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <line x1="9" y1="3" x2="9" y2="21" />
        </>
      ) : (
        <>
          <rect x="3" y="3" width="18" height="18" rx="2" />
          <line x1="15" y1="3" x2="15" y2="21" />
        </>
      )}
    </svg>
  )
}

export function AeolyzerSidebar({ isOpen, onToggle, onNewChat, currentChatTitle, onOpenSettings }: SidebarProps) {
  const [hoveredChat, setHoveredChat] = useState<number | null>(null)
  const [showUserMenu, setShowUserMenu] = useState(false)

  // Sidebar text color - more muted grey like Claude
  const sidebarTextColor = "var(--muted-foreground)"
  const sidebarTextColorHover = "var(--sidebar-text)"

  return (
    <>
      {/* Mobile backdrop */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black/50 z-40 md:hidden" 
          onClick={onToggle}
        />
      )}
      <aside
      className={cn(
        "flex flex-col h-full transition-all duration-300 ease-in-out bg-sidebar-bg shrink-0",
        "max-md:fixed max-md:inset-y-0 max-md:left-0 max-md:z-50 md:relative",
        isOpen ? "w-[260px] max-md:translate-x-0" : "w-[60px] max-md:-translate-x-full"
      )}
    >
      {/* Top section */}
      <div className="flex items-center justify-between p-3 flex-shrink-0">
        {isOpen ? (
          <>
            <span className="text-lg font-medium text-sidebar-text">AEOlyzer</span>
            <button
              onClick={onToggle}
              className="p-1.5 rounded-md transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-muted-foreground"
>
              <SidebarToggleIcon isOpen={true} />
            </button>
          </>
        ) : (
          <button
            onClick={onToggle}
            className="p-1.5 rounded-md transition-colors mx-auto text-sidebar-muted hover:bg-sidebar-hover hover:text-muted-foreground"
>
            <SidebarToggleIcon isOpen={false} />
          </button>
        )}
      </div>

      {/* Navigation items - fixed */}
      <nav className="px-2 flex-shrink-0">
        {/* New Chat */}
        <button
          onClick={onNewChat}
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <Plus size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">New chat</span>}
        </button>

        {/* Search */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <Search size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Search</span>}
        </button>

        {/* Customize */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <Settings2 size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Customize</span>}
        </button>
      </nav>

      {/* Scrollable section */}
      <div className="flex-1 overflow-y-auto px-2 mt-1" style={{ scrollbarWidth: "thin", scrollbarColor: "var(--border) transparent" }}>
        {/* Divider */}
        <div className="my-1.5 border-t border-[var(--border)]" />

        {/* Chats */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <MessageSquare size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Chats</span>}
        </button>

        {/* Projects */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <FolderOpen size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Projects</span>}
        </button>

        {/* Artifacts */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <Grid2X2 size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Artifacts</span>}
        </button>

        {/* Code */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors text-sidebar-muted hover:text-sidebar-text hover:bg-sidebar-hover",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
>
          <Code size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Code</span>}
        </button>

        {/* Recents Section - Scrollable */}
        {isOpen && (
          <>
            <div className="mt-3 mb-1.5 px-3">
              <span className="text-xs text-sidebar-muted">Recents</span>
            </div>
            <div className="pb-4">
              {recentChats.map((chat, index) => (
                <button
                  key={index}
                  className={cn(
                    "flex items-center justify-between w-full px-3 py-1 rounded-lg transition-colors text-left group text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text",
                    currentChatTitle === chat ? "bg-muted text-foreground" : ""
                  )}
                >
                  <span className="text-sm truncate flex-1">{chat}</span>
                  <MoreHorizontal size={14} className="text-sidebar-muted flex-shrink-0 ml-2 opacity-0 group-hover:opacity-100 transition-opacity" />
                </button>
              ))}
            </div>
          </>
        )}
      </div>

      {/* User section - fixed at bottom */}
      <div className="flex-shrink-0 p-2 border-t border-[var(--border)] relative">
        <button
          onClick={() => setShowUserMenu(!showUserMenu)}
          className={cn(
            "flex items-center w-full rounded-lg transition-colors p-2 text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text",
            isOpen ? "gap-3" : "justify-center"
          )}
>
          <div 
            className="w-7 h-7 rounded-full flex items-center justify-center text-sm font-medium flex-shrink-0 bg-border text-foreground"
          >
            M
          </div>
          {isOpen && (
            <>
              <div className="flex-1 min-w-0 text-left">
                <p className="text-sm font-medium truncate text-foreground">Muhammad</p>
                <p className="text-xs text-sidebar-muted">Free plan</p>
              </div>
              <div className="flex items-center gap-1 flex-shrink-0">
                <Download size={16} className="text-sidebar-muted" />
                <ChevronDown size={14} className="text-sidebar-muted" />
              </div>
            </>
          )}
        </button>

        {/* User dropdown menu */}
        {showUserMenu && isOpen && (
          <div 
            className="absolute bottom-full left-2 right-2 mb-2 rounded-xl shadow-xl overflow-hidden z-50"
            style={{ backgroundColor: "var(--background)", border: "1px solid var(--border)" }}
          >
            <div className="px-4 py-3 border-b border-border">
              <p className="text-xs text-sidebar-muted">muhammed.beig@gmail.com</p>
            </div>
            
            <div className="py-1">
              <button
                onClick={() => {
                  setShowUserMenu(false)
                  onOpenSettings()
                }}
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <Settings size={16} strokeWidth={1.5} />
                <span>Settings</span>
                <span className="ml-auto text-xs text-sidebar-muted">Shift+Ctrl+,</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <Globe size={16} strokeWidth={1.5} />
                <span>Language</span>
                <ChevronDown size={14} className="-rotate-90 ml-auto text-sidebar-muted" />
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <HelpCircle size={16} strokeWidth={1.5} />
                <span>Get help</span>
              </button>
            </div>

            <div className="py-1 border-t border-border">
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <Sparkles size={16} strokeWidth={1.5} />
                <span>Upgrade plan</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <Download size={16} strokeWidth={1.5} />
                <span>Get apps and extensions</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <Gift size={16} strokeWidth={1.5} />
                <span>Gift AEOlyzer</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <BookOpen size={16} strokeWidth={1.5} />
                <span>Learn more</span>
                <ChevronDown size={14} className="-rotate-90 ml-auto text-sidebar-muted" />
              </button>
            </div>

            <div className="py-1 border-t border-border">
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors text-muted-foreground hover:bg-sidebar-hover"
>
                <LogOut size={16} strokeWidth={1.5} />
                <span>Log out</span>
              </button>
            </div>
          </div>
        )}
      </div>
    </aside>
    </>
  )
}
