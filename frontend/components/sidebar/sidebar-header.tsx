import { Search } from "lucide-react"
import { cn } from "@/lib/utils"

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

interface SidebarHeaderProps {
  isOpen: boolean
  onToggle: () => void
  onSearchOpen: () => void
}

export function SidebarHeader({ isOpen, onToggle, onSearchOpen }: SidebarHeaderProps) {
  return (
    <div className={cn(
      "flex items-center h-14 flex-shrink-0 transition-all duration-300",
      isOpen ? "px-3 justify-between" : "justify-center hidden md:flex"
    )}>
      {isOpen && <span className="text-lg font-medium text-sidebar-text">AEOlyzer</span>}
      <div className="flex items-center gap-1">
        {isOpen && (
          <button 
            onClick={onSearchOpen}
            className="p-1.5 rounded-md text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text transition-colors" 
            aria-label="Search"
          >
            <Search size={18} />
          </button>
        )}
        <button 
          onClick={onToggle}
          className="p-1.5 rounded-md text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text transition-colors"
          aria-label={isOpen ? "Close sidebar" : "Open sidebar"}
        >
          <SidebarToggleIcon isOpen={isOpen} />
        </button>
      </div>
    </div>
  )
}
