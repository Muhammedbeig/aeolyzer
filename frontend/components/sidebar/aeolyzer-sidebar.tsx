"use client"

import { useState } from "react"
import { cn } from "@/lib/utils"
import { SidebarProps } from "./types"
import { SidebarHeader } from "./sidebar-header"
import { SidebarTabs } from "./sidebar-tabs"
import { SidebarAgentContent } from "./sidebar-agent-content"
import { SidebarHomeContent } from "./sidebar-home-content"
import { SidebarUserProfile } from "./sidebar-user-profile"
import { SearchDialog } from "./search-dialog"

interface AeolyzerSidebarProps {
  isOpen: boolean
  isMobileOpen?: boolean
  onToggle: () => void
  onMobileClose?: () => void
  activeTab: string
  onTabChange: (tab: string) => void
  onNewChat?: () => void
  currentChatTitle?: string
  onOpenSettings?: () => void
}

export function AeolyzerSidebar({ 
  isOpen, 
  isMobileOpen,
  onToggle, 
  onMobileClose,
  activeTab, 
  onTabChange,
  onNewChat,
  currentChatTitle,
  onOpenSettings
}: AeolyzerSidebarProps) {
  const [searchOpen, setSearchOpen] = useState(false)

  const displayOpen = isOpen || isMobileOpen

  const handleToggle = () => {
    // If the mobile sidebar is currently open, clicking the toggle should close it
    if (isMobileOpen && onMobileClose) {
      onMobileClose()
    } else {
      // Otherwise toggle the desktop sidebar
      onToggle()
    }
  }

  return (
    <>
      {/* Mobile Overlay */}
      {isMobileOpen && (
        <div 
          className="fixed inset-0 bg-black/50 z-40 md:hidden transition-opacity"
          onClick={onMobileClose}
        />
      )}
      <aside 
        className={cn(
          "flex flex-col bg-sidebar-bg border-r-[0.5px] border-black/10 dark:border-white/10 transition-all duration-300 ease-in-out h-full relative z-50 font-outfit",
          "max-md:fixed max-md:inset-y-0 max-md:left-0 max-md:w-[260px]",
          isOpen ? "md:w-[260px]" : "w-[0px] md:w-[60px] md:border-r-[0.5px]",
          !isOpen && "max-md:border-r-[0.5px]",
          isMobileOpen ? "max-md:translate-x-0" : "max-md:-translate-x-full"
        )}
      >
        <SidebarHeader 
          isOpen={displayOpen} 
          onToggle={handleToggle} 
          onSearchOpen={() => setSearchOpen(true)} 
        />

        <div className={cn(
          "flex-1 overflow-y-auto overflow-x-hidden custom-scrollbar",
          !displayOpen && "hidden md:block"
        )}>
          <SidebarTabs 
            isOpen={displayOpen} 
            activeTab={activeTab} 
            onTabChange={onTabChange} 
          />

          {displayOpen && (activeTab === 'Agent' || activeTab === 'Content') && (
            <SidebarAgentContent 
              onNewChat={onNewChat} 
              currentChatTitle={currentChatTitle} 
            />
          )}

          {displayOpen && activeTab === 'Home' && (
            <SidebarHomeContent />
          )}
        </div>

        <SidebarUserProfile 
          isOpen={displayOpen} 
          onOpenSettings={onOpenSettings} 
        />
      </aside>

      <SearchDialog 
        open={searchOpen} 
        onOpenChange={setSearchOpen} 
      />
    </>
  )
}
