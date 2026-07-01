"use client"

import { useCallback, useState } from "react"
import { useTheme } from "next-themes"
import { Menu, Plus } from "lucide-react"
import { AeolyzerSidebar } from "@/components/sidebar/aeolyzer-sidebar"
import { AeolyzerChatArea } from "@/components/chat/chat-area"
import { AeolyzerWelcome } from "@/components/chat/welcome-screen"
import { AeolyzerChatInput } from "@/components/chat/chat-input"
import { AeolyzerSettings } from "@/components/settings/aeolyzer-settings"
import { AeolyzerKnowledgeBase } from "@/components/knowledge/knowledge-base"
import type { ConversationSummary } from "@/components/chat/types"
import { useConversations } from "@/hooks/use-conversations"

type Theme = "light" | "system" | "dark"

export default function AeolyzerChatbot() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false)
  const [activeSidebarTab, setActiveSidebarTab] = useState("Home")
  const [activeKnowledgeSection, setActiveKnowledgeSection] = useState("profile")
  const [activeView, setActiveView] = useState("Agent")
  const [settingsOpen, setSettingsOpen] = useState(false)
  const { theme, setTheme } = useTheme()
  const conversations = useConversations()

  const handleTabChange = useCallback(
    (tab: string) => {
      setActiveSidebarTab(tab)
      if (tab === "Agent") {
        conversations.setAgent("audit")
        setActiveView(tab)
      } else if (tab === "Content") {
        conversations.setAgent("content")
        setActiveView(tab)
      } else if (tab === "Home") {
        setActiveView("Home")
      }
    },
    [conversations],
  )

  const handleKnowledgeSectionSelect = useCallback((section: string) => {
    setActiveKnowledgeSection(section)
    setActiveView("Knowledge")
  }, [])

  const handleNewChat = useCallback(() => {
    conversations.beginNewConversation()
    setActiveView(conversations.agent === "content" ? "Content" : "Agent")
  }, [conversations])

  const handleSelectConversation = useCallback(
    (conversation: ConversationSummary) => {
      const tab = conversation.agent === "content" ? "Content" : "Agent"
      setActiveSidebarTab(tab)
      setActiveView(tab)
      void conversations.selectConversation(conversation)
    },
    [conversations],
  )

  const handleSendMessage = useCallback(
    (content: string, files: File[] = []) => {
      void conversations.submitMessage(content, files)
    },
    [conversations],
  )

  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <AeolyzerSidebar
        isOpen={sidebarOpen}
        isMobileOpen={mobileSidebarOpen}
        onToggle={() => setSidebarOpen(!sidebarOpen)}
        onMobileClose={() => setMobileSidebarOpen(false)}
        onNewChat={handleNewChat}
        conversations={conversations.currentConversations}
        allConversations={conversations.allConversations}
        activeConversationID={conversations.activeConversationID}
        onSelectConversation={handleSelectConversation}
        onToggleStar={(conversation) => {
          void conversations.toggleStar(conversation)
        }}
        activeTab={activeSidebarTab}
        onTabChange={handleTabChange}
        activeKnowledgeSection={activeKnowledgeSection}
        onKnowledgeSectionChange={handleKnowledgeSectionSelect}
        onOpenSettings={() => setSettingsOpen(true)}
      />

      <main className="flex-1 flex flex-col min-w-0 min-h-0 relative">
        <div className="md:hidden flex items-center justify-between p-3 border-b-[0.5px] border-black/10 dark:border-white/10 bg-background flex-shrink-0 z-10">
          <button
            onClick={() => setMobileSidebarOpen(true)}
            className="p-1.5 rounded-md text-muted-foreground hover:bg-muted"
            aria-label="Open sidebar"
          >
            <Menu size={20} strokeWidth={1.5} />
          </button>
          <span className="text-sm font-medium text-foreground">AEOlyzer</span>
          <button
            onClick={handleNewChat}
            className="p-1.5 rounded-md text-muted-foreground hover:bg-muted"
            aria-label="New chat"
          >
            <Plus size={20} strokeWidth={1.5} />
          </button>
        </div>

        {activeView === "Knowledge" ? (
          <div className="flex-1 flex flex-col bg-background overflow-y-auto">
            <AeolyzerKnowledgeBase activeSection={activeKnowledgeSection} />
          </div>
        ) : activeView === "Content" && conversations.messages.length === 0 ? (
          <div className="flex-1 flex flex-col bg-background overflow-y-auto">
            <AeolyzerWelcome
              title="What can I help you create?"
              placeholder="Describe what you want to write..."
              showContentOptions
              onSend={handleSendMessage}
              isGenerating={conversations.isGenerating}
            />
          </div>
        ) : conversations.messages.length === 0 ? (
          <div className="flex-1 flex flex-col bg-background overflow-y-auto">
            <AeolyzerWelcome
              onSend={handleSendMessage}
              isGenerating={conversations.isGenerating}
            />
          </div>
        ) : (
          <>
            <AeolyzerChatArea
              messages={conversations.messages}
              isGenerating={conversations.isGenerating}
              chatTitle={conversations.currentConversation?.title}
              onSend={handleSendMessage}
            />

            <div className="flex-shrink-0 px-4 pb-4 pt-2 bg-background">
              <AeolyzerChatInput
                onSend={handleSendMessage}
                isGenerating={conversations.isGenerating}
                placeholder="Reply..."
              />
            </div>
          </>
        )}
      </main>

      <AeolyzerSettings
        isOpen={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        theme={(theme as Theme) || "system"}
        onThemeChange={(nextTheme) => setTheme(nextTheme)}
      />
    </div>
  )
}
