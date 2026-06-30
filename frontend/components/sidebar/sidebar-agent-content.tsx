import { useState } from "react"
import { Plus, MoreHorizontal, Star } from "lucide-react"
import { cn } from "@/lib/utils"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { RECENT_CHATS } from "./constants"

interface SidebarAgentContentProps {
  onNewChat: () => void
  currentChatTitle?: string
}

export function SidebarAgentContent({ onNewChat, currentChatTitle }: SidebarAgentContentProps) {
  const [recents, setRecents] = useState<string[]>(RECENT_CHATS)
  const [starred, setStarred] = useState<string[]>(["Core software engineering princ..."])

  const handleStarToggle = (chat: string) => {
    if (starred.includes(chat)) {
      setStarred(starred.filter(c => c !== chat))
      setRecents([chat, ...recents])
    } else {
      setStarred([...starred, chat])
      setRecents(recents.filter(c => c !== chat))
    }
  }

  return (
    <>
      {/* Primary Actions */}
      <div className="px-3 py-2 space-y-0.5">
        <button onClick={onNewChat} className="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group transition-colors" aria-label="New chat">
          <div className="w-5 h-5 rounded-full border-[0.5px] border-black/10 dark:border-white/10 flex items-center justify-center group-hover:border-black/20 dark:group-hover:border-white/20 transition-colors">
            <Plus size={14} />
          </div>
          <span className="truncate text-[15px]">New chat</span>
        </button>
      </div>

      {/* Starred Section */}
      {starred.length > 0 && (
        <div className="px-3 py-2">
          <div className="px-3 pb-2 pt-1">
            <h3 className="text-xs font-medium text-sidebar-muted">Starred</h3>
          </div>
          <div className="space-y-0.5">
            {starred.map((chat, index) => (
              <div key={`starred-${index}`} className="flex items-center justify-between w-full px-3 py-2 rounded-lg group transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text cursor-pointer">
                <span className="truncate text-sm pr-2 text-left">{chat}</span>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button className="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 rounded hover:bg-black/10 dark:hover:bg-white/10" aria-label="More options">
                      <MoreHorizontal size={14} />
                    </button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="start" className="w-40 bg-popover border-[0.5px] border-black/10 dark:border-white/10">
                    <DropdownMenuItem onClick={(e) => { e.stopPropagation(); handleStarToggle(chat); }} className="gap-2 cursor-pointer">
                      <Star size={14} className="fill-muted-foreground text-muted-foreground" />
                      Unstar
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Recents Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1 flex items-center justify-between">
          <h3 className="text-xs font-medium text-sidebar-muted">Recents</h3>
          <button className="p-1 rounded opacity-0 hover:bg-sidebar-hover text-sidebar-muted hover:text-sidebar-text transition-all" aria-hidden="true" disabled>
            <Plus size={14} strokeWidth={2} className="invisible" />
          </button>
        </div>
        <div className="space-y-0.5">
          {recents.slice(0, 8).map((chat, index) => (
            <div key={`recent-${index}`} className={cn("flex items-center justify-between w-full px-3 py-2 rounded-lg group transition-colors cursor-pointer", chat === currentChatTitle ? "bg-sidebar-hover text-sidebar-text" : "text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text")}>
              <span className="truncate text-sm pr-2 text-left">{chat}</span>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <button className="opacity-0 group-hover:opacity-100 transition-opacity p-0.5 rounded hover:bg-black/10 dark:hover:bg-white/10" aria-label="More options">
                    <MoreHorizontal size={14} />
                  </button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="start" className="w-40 bg-popover border-[0.5px] border-black/10 dark:border-white/10">
                  <DropdownMenuItem onClick={(e) => { e.stopPropagation(); handleStarToggle(chat); }} className="gap-2 cursor-pointer">
                    <Star size={14} className="text-muted-foreground" />
                    Star
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          ))}
        </div>
      </div>
    </>
  )
}
