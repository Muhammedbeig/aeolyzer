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

export function AeolyzerSidebar({ 
  isOpen, 
  onToggle, 
  onNewChat, 
  currentChatTitle, 
  onOpenSettings, 
  activeTab, 
  onTabChange 
}: SidebarProps) {
  const [searchOpen, setSearchOpen] = useState(false)

  return (
    <>
      {/* Mobile backdrop */}
      {isOpen && (
        <div 
          className="fixed inset-0 bg-black/50 z-40 md:hidden" 
          onClick={onToggle}
          aria-hidden="true"
        />
      )}
      <aside 
        className={cn(
          "flex flex-col bg-sidebar-bg border-r-[0.5px] border-black/10 dark:border-white/10 transition-all duration-300 ease-in-out h-full relative z-50",
          "max-md:fixed max-md:inset-y-0 max-md:left-0",
          isOpen ? "w-[260px] max-md:translate-x-0" : "w-[0px] border-r-0 md:w-[60px] md:border-r-[0.5px] max-md:-translate-x-full"
        )}
      >
        <SidebarHeader 
          isOpen={isOpen} 
          onToggle={onToggle} 
          onSearchOpen={() => setSearchOpen(true)} 
        />

        <div className={cn(
          "flex-1 overflow-y-auto overflow-x-hidden custom-scrollbar",
          !isOpen && "hidden md:block"
        )}>
          <SidebarTabs 
            isOpen={isOpen} 
            activeTab={activeTab} 
            onTabChange={onTabChange} 
          />

          {isOpen && (activeTab === 'Agent' || activeTab === 'Content') && (
            <SidebarAgentContent 
              onNewChat={onNewChat} 
              currentChatTitle={currentChatTitle} 
            />
          )}

          {isOpen && activeTab === 'Home' && (
            <SidebarHomeContent />
          )}
        </div>

        <SidebarUserProfile 
          isOpen={isOpen} 
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
