import { Settings } from "lucide-react"
import { cn } from "@/lib/utils"

interface SidebarUserProfileProps {
  isOpen: boolean
  onOpenSettings: () => void
}

export function SidebarUserProfile({ isOpen, onOpenSettings }: SidebarUserProfileProps) {
  return (
    <div className={cn(
      "p-3 border-t border-border/10 mt-auto flex-shrink-0",
      !isOpen && "hidden md:flex justify-center"
    )}>
      <div className="flex items-center gap-3">
        <div className="w-8 h-8 rounded-full bg-accent/20 flex items-center justify-center text-accent font-medium text-sm flex-shrink-0 cursor-pointer">
          M
        </div>
        {isOpen && (
          <div className="flex-1 min-w-0 flex items-center justify-between">
            <div className="flex flex-col min-w-0 pr-2 cursor-pointer">
              <span className="text-[15px] font-medium text-sidebar-text truncate">Muhammad</span>
              <span className="text-xs text-sidebar-muted truncate">Free plan</span>
            </div>
            <button 
              onClick={onOpenSettings}
              className="p-1.5 rounded hover:bg-sidebar-hover text-sidebar-muted transition-colors flex-shrink-0"
              aria-label="Settings"
            >
              <Settings size={18} />
            </button>
          </div>
        )}
      </div>
    </div>
  )
}
