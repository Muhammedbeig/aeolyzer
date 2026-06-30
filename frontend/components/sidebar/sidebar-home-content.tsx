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
          <button className="w-full flex items-center justify-between px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="LLM Analytics">
            <div className="flex items-center gap-3">
              <BarChart2 size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
              <span className="truncate text-[15px]">LLM Analytics</span>
            </div>
            <span className="text-[10px] px-1.5 py-0.5 rounded-full bg-accent/10 text-[#a53b15] dark:text-accent font-medium border-[0.5px] border-black/10 dark:border-white/10">Beta</span>
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
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Prompt Research">
            <Search size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">Prompt Research</span>
          </button>
        </div>
      </div>

      {/* Actions Section */}
      <div className="px-3 py-2">
        <div className="px-3 pb-2 pt-1">
          <h3 className="text-xs font-medium text-sidebar-muted">Actions</h3>
        </div>
        <div className="space-y-0.5">
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="Opportunities">
            <Zap size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">Opportunities</span>
          </button>
          <button className="w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors text-sidebar-muted hover:bg-sidebar-hover hover:text-sidebar-text group" aria-label="On Page">
            <Layout size={20} strokeWidth={1.5} className="flex-shrink-0 group-hover:text-sidebar-text" />
            <span className="truncate text-[15px]">On Page</span>
          </button>
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
