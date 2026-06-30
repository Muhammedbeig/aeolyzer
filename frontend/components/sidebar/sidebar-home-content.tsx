import { LineChart, BarChart2, Activity, Terminal, Search, Zap, Layout, ActivitySquare, AlertCircle } from "lucide-react"

export function SidebarHomeContent() {
  return (
    <>
      {/* Analytics Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted">Analytics</h3>
        </div>
        <div className="space-y-0.5">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="AEO Insights">
            <LineChart size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">AEO Insights</span>
          </button>

          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Traffic">
            <Activity size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">Traffic</span>
          </button>
        </div>
      </div>

      {/* Prompts Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted">Prompts</h3>
        </div>
        <div className="space-y-0.5">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Your Prompts">
            <Terminal size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">Your Prompts</span>
          </button>

        </div>
      </div>

      {/* Actions Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted">Actions</h3>
        </div>
        <div className="space-y-0.5">

          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Site Health">
            <ActivitySquare size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">Site Health</span>
          </button>
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Issues">
            <AlertCircle size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">Issues</span>
          </button>
        </div>
      </div>
    </>
  )
}
