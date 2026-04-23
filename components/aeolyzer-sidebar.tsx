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

// -> defines the shape of the data (props) the AeolyzerSidebar block expects to receive
// - `isOpen` -> true/false switch telling us if the sidebar is expanded (wide) or collapsed (narrow)
// - `onToggle` -> function to call when the user clicks the open/close button
// - `onNewChat` -> function to call when the user clicks the "+" New Chat button
// - `currentChatTitle` -> optional text to highlight which chat is currently active in the list
// - `onOpenSettings` -> function to call when the user clicks Settings in their user menu
interface SidebarProps {
  isOpen: boolean
  onToggle: () => void
  onNewChat: () => void
  currentChatTitle?: string
  onOpenSettings: () => void
}

// -> a fake list of past conversations (history) hardcoded for demonstration purposes
const recentChats = [
  "Your greatest strengths",
  "Downloaded page color mismatch",
  "Answer engine optimization key...",
  "Futbol fibre audience analysis",
  "Football websites that went viral...",
  "Growing soccer live score websi...",
  "SEO-optimized external links",
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
// -> A mini-component that draws the exact little SVG icon used to open/close the sidebar
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

// -> This is the main Sidebar panel that holds the navigation links, chat history, and user settings
export function AeolyzerSidebar({ isOpen, onToggle, onNewChat, currentChatTitle, onOpenSettings }: SidebarProps) {
  // - `const [hoveredChat, setHoveredChat]` -> memory tracking which past chat your mouse is currently floating over
  const [hoveredChat, setHoveredChat] = useState<number | null>(null)
  
  // - `const [showUserMenu, setShowUserMenu]` -> memory toggling the bottom popup menu (Settings, Log out, etc.)
  const [showUserMenu, setShowUserMenu] = useState(false)

  // Sidebar text color - more muted grey like Claude
  // -> Variables holding the exact grey colors we want to use for text
  const sidebarTextColor = "#a3a29e"
  const sidebarTextColorHover = "#d4d4d4"

  return (
    <aside
      className={cn(
        "flex flex-col h-full transition-all duration-300 ease-in-out relative",
        isOpen ? "w-[260px]" : "w-[60px]"
      )}
      style={{ backgroundColor: "#2b2a27" }}
    >
      {/* Top section */}
      <div className="flex items-center justify-between p-3 flex-shrink-0">
        {isOpen ? (
          <>
            <span className="text-lg font-medium" style={{ color: "#d4d4d0" }}>AEOlyzer</span>
            <button
              onClick={onToggle}
              className="p-1.5 rounded-md transition-colors"
              style={{ color: "#6b6b66" }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = "#252422"
                e.currentTarget.style.color = "#a3a29e"
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = "transparent"
                e.currentTarget.style.color = "#6b6b66"
              }}
            >
              <SidebarToggleIcon isOpen={true} />
            </button>
          </>
        ) : (
          <button
            onClick={onToggle}
            className="p-1.5 rounded-md transition-colors mx-auto"
            style={{ color: "#6b6b66" }}
            onMouseEnter={(e) => {
              e.currentTarget.style.backgroundColor = "#252422"
              e.currentTarget.style.color = "#a3a29e"
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.backgroundColor = "transparent"
              e.currentTarget.style.color = "#6b6b66"
            }}
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
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <Plus size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">New chat</span>}
        </button>

        {/* Search */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <Search size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Search</span>}
        </button>

        {/* Customize */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <Settings2 size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Customize</span>}
        </button>
      </nav>

      {/* Scrollable section */}
      <div className="flex-1 overflow-y-auto px-2 mt-1" style={{ scrollbarWidth: "thin", scrollbarColor: "#4a4945 transparent" }}>
        {/* Divider */}
        <div className="my-1.5 border-t border-[#3a3936]" />

        {/* Chats */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <MessageSquare size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Chats</span>}
        </button>

        {/* Projects */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <FolderOpen size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Projects</span>}
        </button>

        {/* Artifacts */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <Grid2X2 size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Artifacts</span>}
        </button>

        {/* Code */}
        <button
          className={cn(
            "flex items-center w-full rounded-lg transition-colors",
            isOpen ? "gap-3 px-3 py-1.5" : "justify-center p-2"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
            e.currentTarget.style.color = sidebarTextColorHover
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
            e.currentTarget.style.color = sidebarTextColor
          }}
        >
          <Code size={18} strokeWidth={1.5} />
          {isOpen && <span className="text-sm">Code</span>}
        </button>

        {/* Recents Section - Scrollable */}
        {isOpen && (
          <>
            <div className="mt-3 mb-1.5 px-3">
              <span className="text-xs" style={{ color: "#6b6b66" }}>Recents</span>
            </div>
            <div className="pb-4">
              {recentChats.map((chat, index) => (
                <button
                  key={index}
                  onMouseEnter={() => setHoveredChat(index)}
                  onMouseLeave={() => setHoveredChat(null)}
                  className={cn(
                    "flex items-center justify-between w-full px-3 py-1 rounded-lg transition-colors text-left group",
                    currentChatTitle === chat ? "bg-[#252422]" : ""
                  )}
                  style={{ 
                    color: sidebarTextColor,
                    backgroundColor: hoveredChat === index && currentChatTitle !== chat ? "#252422" : undefined
                  }}
                >
                  <span className="text-sm truncate flex-1">{chat}</span>
                  {hoveredChat === index && (
                    <MoreHorizontal size={14} style={{ color: "#6b6b66" }} className="flex-shrink-0 ml-2" />
                  )}
                </button>
              ))}
            </div>
          </>
        )}
      </div>

      {/* User section - fixed at bottom */}
      <div className="flex-shrink-0 p-2 border-t border-[#3a3936] relative">
        <button
          onClick={() => setShowUserMenu(!showUserMenu)}
          className={cn(
            "flex items-center w-full rounded-lg transition-colors p-2",
            isOpen ? "gap-3" : "justify-center"
          )}
          style={{ color: sidebarTextColor }}
          onMouseEnter={(e) => {
            e.currentTarget.style.backgroundColor = "#252422"
          }}
          onMouseLeave={(e) => {
            e.currentTarget.style.backgroundColor = "transparent"
          }}
        >
          <div 
            className="w-7 h-7 rounded-full flex items-center justify-center text-sm font-medium flex-shrink-0"
            style={{ backgroundColor: "#4a4945", color: "#ececec" }}
          >
            M
          </div>
          {isOpen && (
            <>
              <div className="flex-1 min-w-0 text-left">
                <p className="text-sm font-medium truncate" style={{ color: "#ececec" }}>Muhammad</p>
                <p className="text-xs" style={{ color: "#6b6b66" }}>Free plan</p>
              </div>
              <div className="flex items-center gap-1 flex-shrink-0">
                <Download size={16} style={{ color: "#6b6b66" }} />
                <ChevronDown size={14} style={{ color: "#6b6b66" }} />
              </div>
            </>
          )}
        </button>

        {/* User dropdown menu */}
        {showUserMenu && isOpen && (
          <div 
            className="absolute bottom-full left-2 right-2 mb-2 rounded-xl shadow-xl overflow-hidden z-50"
            style={{ backgroundColor: "#2b2a27", border: "1px solid #3a3936" }}
          >
            <div className="px-4 py-3 border-b" style={{ borderColor: "#3a3936" }}>
              <p className="text-xs" style={{ color: "#6b6b66" }}>muhammed.beig@gmail.com</p>
            </div>
            
            <div className="py-1">
              <button
                onClick={() => {
                  setShowUserMenu(false)
                  onOpenSettings()
                }}
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <Settings size={16} strokeWidth={1.5} />
                <span>Settings</span>
                <span className="ml-auto text-xs" style={{ color: "#6b6b66" }}>Shift+Ctrl+,</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <Globe size={16} strokeWidth={1.5} />
                <span>Language</span>
                <ChevronDown size={14} className="-rotate-90 ml-auto" style={{ color: "#6b6b66" }} />
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <HelpCircle size={16} strokeWidth={1.5} />
                <span>Get help</span>
              </button>
            </div>

            <div className="py-1 border-t" style={{ borderColor: "#3a3936" }}>
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <Sparkles size={16} strokeWidth={1.5} />
                <span>Upgrade plan</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <Download size={16} strokeWidth={1.5} />
                <span>Get apps and extensions</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <Gift size={16} strokeWidth={1.5} />
                <span>Gift AEOlyzer</span>
              </button>
              
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <BookOpen size={16} strokeWidth={1.5} />
                <span>Learn more</span>
                <ChevronDown size={14} className="-rotate-90 ml-auto" style={{ color: "#6b6b66" }} />
              </button>
            </div>

            <div className="py-1 border-t" style={{ borderColor: "#3a3936" }}>
              <button
                className="flex items-center gap-3 w-full px-4 py-2 text-sm transition-colors"
                style={{ color: "#a3a29e" }}
                onMouseEnter={(e) => e.currentTarget.style.backgroundColor = "#252422"}
                onMouseLeave={(e) => e.currentTarget.style.backgroundColor = "transparent"}
              >
                <LogOut size={16} strokeWidth={1.5} />
                <span>Log out</span>
              </button>
            </div>
          </div>
        )}
      </div>
    </aside>
  )
}
