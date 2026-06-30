import { X, Settings, User, Shield, CreditCard, Layers, Plug, Code, Search } from "lucide-react"
import { cn } from "@/lib/utils"
import { SETTINGS_TABS, SettingsTab } from "./constants"

interface SettingsNavProps {
  activeTab: SettingsTab | string;
  onTabChange: (tab: string) => void;
  onClose: () => void;
}

const TAB_ICONS: Record<string, React.ElementType> = {
  "General": Settings,
  "Account": User,
  "Privacy": Shield,
  "Billing": CreditCard,
  "Capabilities": Layers,
  "Connectors": Plug,
  "Aeolyzer Code": Code,
}

export function SettingsNav({ activeTab, onTabChange, onClose }: SettingsNavProps) {
  return (
    <div className="md:w-[260px] flex-shrink-0 flex flex-col border-b md:border-b-0 md:border-r border-border">
      <div className="flex flex-col p-4 md:pb-2 md:pt-4">
        {/* Mobile Header */}
        <div className="flex items-center justify-between md:hidden">
          <h2 className="text-xl font-semibold text-foreground">Settings</h2>
          <button 
            onClick={onClose}
            className="p-1 rounded-lg transition-colors hover:bg-muted text-muted-foreground"
          >
            <X size={20} />
          </button>
        </div>

        {/* Desktop Search & Header */}
        <div className="hidden md:flex flex-col">
          <div className="flex items-center gap-2 px-3 py-1.5 mb-6 rounded-lg bg-white/5 dark:bg-white/5 border border-transparent focus-within:border-border transition-colors">
            <Search size={16} className="text-muted-foreground flex-shrink-0" />
            <input 
              type="text" 
              placeholder="Search" 
              className="w-full bg-transparent border-none outline-none text-sm placeholder:text-muted-foreground text-foreground h-6"
            />
          </div>
          <h2 className="text-[13px] font-medium text-muted-foreground px-2">Settings</h2>
        </div>
      </div>
      <nav className="flex md:flex-col overflow-x-auto hide-scrollbar px-4 pb-2 md:px-4 md:pb-4 space-x-2 md:space-x-0 md:space-y-1 snap-x">
        {SETTINGS_TABS.map((tab) => {
          const Icon = TAB_ICONS[tab]
          return (
            <button
              key={tab}
              onClick={() => onTabChange(tab)}
              className={cn(
                "flex items-center gap-3 flex-shrink-0 px-3 py-1.5 md:px-3 md:py-2 rounded-full md:rounded-lg text-sm font-medium transition-colors snap-start md:w-full md:justify-start",
                activeTab === tab 
                  ? "bg-muted text-foreground" 
                  : "text-muted-foreground hover:bg-muted hover:text-foreground"
              )}
            >
              {Icon && <Icon size={16} strokeWidth={2} className="hidden md:block flex-shrink-0" />}
              <span>{tab}</span>
            </button>
          )
        })}
      </nav>
    </div>
  )
}
