import { LineChart, BarChart2, Activity, Terminal, Search, Zap, Layout, ActivitySquare, AlertCircle, FileText, Award, Users, List, Mic, Brain } from "lucide-react"
import { cn } from "@/lib/utils"

interface SidebarHomeContentProps {
  isOpen: boolean
  activeSection?: string
  onSectionChange?: (section: string) => void
}

export function SidebarHomeContent({ isOpen, activeSection, onSectionChange }: SidebarHomeContentProps) {
  const knowledgeSections = [
    { id: "profile", label: "Profile", icon: FileText },
    { id: "eeat", label: "EEAT", icon: Award },
    { id: "competitors", label: "Competitors", icon: Users },
    { id: "topics", label: "Topics", icon: List },
    { id: "tone", label: "Tone", icon: Mic },
    { id: "memory", label: "Memory", icon: Brain },
  ]

  if (!isOpen) {
    return (
      <div className="px-3 py-2 space-y-2 flex flex-col items-center">
        {/* Analytics Icons */}
        <button className="p-2 w-full flex justify-center rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" title="AEO Insights">
          <LineChart size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
        </button>
        <button className="p-2 w-full flex justify-center rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" title="Traffic">
          <Activity size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
        </button>
        
        {/* Divider */}
        <div className="w-8 h-[1px] bg-black/10 dark:bg-white/10 my-1" />

        {/* Prompts Icons */}
        <button className="p-2 w-full flex justify-center rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" title="Your Prompts">
          <Terminal size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
        </button>

        {/* Divider */}
        <div className="w-8 h-[1px] bg-black/10 dark:bg-white/10 my-1" />

        {/* Actions Icons */}
        <button className="p-2 w-full flex justify-center rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" title="Site Health">
          <ActivitySquare size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
        </button>
        <button className="p-2 w-full flex justify-center rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" title="Issues">
          <AlertCircle size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
        </button>

        {/* Divider */}
        <div className="w-8 h-[1px] bg-black/10 dark:bg-white/10 my-1" />

        {/* Knowledge Base Icons */}
        {knowledgeSections.map((section) => (
          <button
            key={section.id}
            onClick={() => onSectionChange?.(section.id)}
            className={cn(
              "p-2 w-full flex justify-center rounded-lg transition-colors group",
              activeSection === section.id
                ? "bg-sidebar-hover text-sidebar-text font-medium"
                : "text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text"
            )}
            title={section.label}
          >
            <section.icon size={20} strokeWidth={1.5} className="flex-shrink-0" />
          </button>
        ))}
      </div>
    )
  }

  return (
    <>
      {/* Analytics Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted uppercase tracking-wider">Analytics</h3>
        </div>
        <div className="space-y-0.5">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="AEO Insights">
            <LineChart size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[14px]">AEO Insights</span>
          </button>

          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Traffic">
            <Activity size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[14px]">Traffic</span>
          </button>
        </div>
      </div>

      {/* Prompts Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted uppercase tracking-wider">Prompts</h3>
        </div>
        <div className="space-y-0.5">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Your Prompts">
            <Terminal size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[14px]">Your Prompts</span>
          </button>
        </div>
      </div>

      {/* Actions Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted uppercase tracking-wider">Actions</h3>
        </div>
        <div className="space-y-0.5">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Site Health">
            <ActivitySquare size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[14px]">Site Health</span>
          </button>
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Issues">
            <AlertCircle size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[14px]">Issues</span>
          </button>
        </div>
      </div>

      {/* Knowledge Base Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted uppercase tracking-wider">Knowledge Base</h3>
        </div>
        <div className="space-y-0.5">
          {knowledgeSections.map((section) => (
            <button
              key={section.id}
              onClick={() => onSectionChange?.(section.id)}
              className={cn(
                "w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors group",
                activeSection === section.id
                  ? "bg-sidebar-hover text-sidebar-text font-medium"
                  : "text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text"
              )}
            >
              <section.icon size={18} strokeWidth={1.5} className="flex-shrink-0" />
              <span className="truncate text-[14px]">{section.label}</span>
            </button>
          ))}
        </div>
      </div>
    </>
  )
}
