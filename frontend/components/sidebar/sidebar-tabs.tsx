import { LayoutDashboard, Bot, FileText } from "lucide-react"
import { cn } from "@/lib/utils"

interface SidebarTabsProps {
  isOpen: boolean
  activeTab: string
  onTabChange: (tab: string) => void
}

export function SidebarTabs({ isOpen, activeTab, onTabChange }: SidebarTabsProps) {
  if (isOpen) {
    return (
      <div className="px-2 pt-2 pb-1">
        <div role="tablist" aria-label="Sidebar mode" className="relative flex items-center gap-0.5 p-0.5 rounded-lg bg-black/5 dark:bg-white/5">
          {/* Sliding Background */}
          <div 
            aria-hidden="true" 
            className="absolute top-0.5 bottom-0.5 rounded-md bg-white dark:bg-zinc-800 shadow-sm pointer-events-none transition-all duration-200 ease-out" 
            style={{
              width: 'calc(33.333% - 1.333px)',
              left: activeTab === 'Home' ? '2px' : activeTab === 'Agent' ? 'calc(33.333% + 0.667px)' : 'calc(66.667% - 0.667px)'
            }}
          />
          <button 
            type="button" 
            role="tab" 
            aria-selected={activeTab === 'Home'} 
            onClick={() => onTabChange('Home')}
            className={cn(
              "relative z-10 flex-1 h-7 rounded-md text-[12.5px] tracking-[-0.005em] cursor-pointer transition-colors duration-150 ease-out",
              activeTab === 'Home' 
                ? "text-zinc-900 dark:text-zinc-100 font-semibold" 
                : "text-zinc-600 dark:text-zinc-400 font-medium hover:text-zinc-900 dark:hover:text-zinc-200"
            )}
          >
            Home
          </button>
          <button 
            type="button" 
            role="tab" 
            aria-selected={activeTab === 'Agent'} 
            onClick={() => onTabChange('Agent')}
            className={cn(
              "relative z-10 flex-1 h-7 rounded-md text-[12.5px] tracking-[-0.005em] cursor-pointer transition-colors duration-150 ease-out",
              activeTab === 'Agent' 
                ? "text-zinc-900 dark:text-zinc-100 font-semibold" 
                : "text-zinc-600 dark:text-zinc-400 font-medium hover:text-zinc-900 dark:hover:text-zinc-200"
            )}
          >
            Agent
          </button>
          <button 
            type="button" 
            role="tab" 
            aria-selected={activeTab === 'Content'} 
            onClick={() => onTabChange('Content')}
            className={cn(
              "relative z-10 flex-1 h-7 rounded-md text-[12.5px] tracking-[-0.005em] cursor-pointer transition-colors duration-150 ease-out",
              activeTab === 'Content' 
                ? "text-zinc-900 dark:text-zinc-100 font-semibold" 
                : "text-zinc-600 dark:text-zinc-400 font-medium hover:text-zinc-900 dark:hover:text-zinc-200"
            )}
          >
            Content
          </button>
        </div>
      </div>
    )
  }

  return (
    <div className="px-3 py-2 space-y-0.5 flex flex-col items-center">
      <button
        onClick={() => onTabChange('Home')}
        className={cn(
          "p-2 w-full flex justify-center rounded-lg transition-colors group",
          activeTab === 'Home' 
            ? "bg-sidebar-hover text-sidebar-text font-medium" 
            : "text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text"
        )}
        title="Home"
        aria-label="Home"
      >
        <LayoutDashboard size={20} strokeWidth={1.5} className={cn("flex-shrink-0", activeTab === 'Home' ? "text-sidebar-text" : "text-sidebar-muted group-hover:text-sidebar-text")} />
      </button>

      <button
        onClick={() => onTabChange('Agent')}
        className={cn(
          "p-2 w-full flex justify-center rounded-lg transition-colors group",
          activeTab === 'Agent' 
            ? "bg-sidebar-hover text-sidebar-text font-medium" 
            : "text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text"
        )}
        title="Agent"
        aria-label="Agent"
      >
        <Bot size={20} strokeWidth={1.5} className={cn("flex-shrink-0", activeTab === 'Agent' ? "text-sidebar-text" : "text-sidebar-muted group-hover:text-sidebar-text")} />
      </button>

      <button
        onClick={() => onTabChange('Content')}
        className={cn(
          "p-2 w-full flex justify-center rounded-lg transition-colors group",
          activeTab === 'Content' 
            ? "bg-sidebar-hover text-sidebar-text font-medium" 
            : "text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text"
        )}
        title="Content"
        aria-label="Content"
      >
        <FileText size={20} strokeWidth={1.5} className={cn("flex-shrink-0", activeTab === 'Content' ? "text-sidebar-text" : "text-sidebar-muted group-hover:text-sidebar-text")} />
      </button>
    </div>
  )
}
