"use client"

import {
  ArrowRight,
  BarChart3,
  Check,
  CheckCircle2,
  ChevronDown,
  ChevronLeft,
  ChevronRight,
  CircleHelp,
  Clock3,
  Command,
  Database,
  Eye,
  FileText,
  Gift,
  Globe2,
  GraduationCap,
  HeartPulse,
  Lightbulb,
  ListChecks,
  MessageCircle,
  MessageSquareText,
  Paperclip,
  PenLine,
  Plus,
  Search,
  Send,
  Stethoscope,
  Target,
  TrendingUp,
  WandSparkles,
  Zap,
} from "lucide-react"
import { useEffect, useState } from "react"
import { ThemeToggle, type ThemeMode } from "@/components/aeolyzer-app"
import { AeolyzerLogo } from "@/components/aeolyzer-logo"
import {
  type ChatMessage,
  type ProjectProfile,
  contentTypes,
  quickActions,
} from "@/lib/aeolyzer"

type Mode = "dashboard" | "agent" | "content"
type DashboardTab =
  | "aeo-insights"
  | "traffic"
  | "your-prompts"
  | "prompt-research"
  | "site-health"

const auditHistoryKey = "aeolyzer.guest.audit-history"
const contentHistoryKey = "aeolyzer.guest.content-history"
const dashboardActionOrder = [
  "Audit my Homepage",
  "Quick Visibility Check",
  "Analyze competitors",
  "Make me a Task List",
  "Optimize Meta tags",
  "Generate a llm.txt",
  "Review Technical SEO",
  "Summarize recent Performance",
  "Find Keyword opportunities",
  "Draft a Blog Post",
  "Find content gaps",
]
const agentActionOrder = [
  "Optimize Meta tags",
  "Make me a Task List",
  "Quick Visibility Check",
  "Review Technical SEO",
  "Audit my Homepage",
  "Analyze competitors",
  "Generate a llm.txt",
  "Summarize recent Performance",
  "Find Keyword opportunities",
  "Draft a Blog Post",
  "Find content gaps",
]
const sidebarButton =
  "flex h-7 w-full items-center gap-2 rounded-lg px-2 text-left text-[10px] text-[#6f6962] transition hover:bg-[#eeeae5] hover:text-[#1e1b18] dark:text-[#a3a29e] dark:hover:bg-[#393836] dark:hover:text-[#ececec]"
const softCard =
  "rounded-xl border border-[#e6e1da] bg-white dark:border-[#4a4945] dark:bg-[#393836]"
const muted = "text-[#77716a] dark:text-[#a3a29e]"
const heading =
  "font-display font-medium tracking-[-0.035em] text-[#171512] dark:text-[#ececec]"

export function AeolyzerWorkspace({
  project,
  theme,
  onThemeChange,
  onReset,
}: {
  project: ProjectProfile
  theme: ThemeMode
  onThemeChange: (theme: ThemeMode) => void
  onReset: () => void
}) {
  const [mode, setMode] = useState<Mode>("dashboard")
  const [tab, setTab] = useState<DashboardTab>("aeo-insights")
  const [auditMessages, setAuditMessages] = useState<ChatMessage[]>(() =>
    readHistory(auditHistoryKey),
  )
  const [contentMessages, setContentMessages] = useState<ChatMessage[]>(() =>
    readHistory(contentHistoryKey),
  )
  const [auditInput, setAuditInput] = useState("")
  const [contentInput, setContentInput] = useState("")
  const [contentType, setContentType] = useState(contentTypes[0])
  const [dashboardActionPage, setDashboardActionPage] = useState(0)
  const [agentActionPage, setAgentActionPage] = useState(0)

  useEffect(() => {
    sessionStorage.setItem(auditHistoryKey, JSON.stringify(auditMessages))
  }, [auditMessages])

  useEffect(() => {
    sessionStorage.setItem(contentHistoryKey, JSON.stringify(contentMessages))
  }, [contentMessages])

  const dashboardActions = actionPage(dashboardActionOrder, dashboardActionPage)
  const agentActions = actionPage(agentActionOrder, agentActionPage)

  function sendAudit(value: string) {
    const prompt = value.trim()
    if (!prompt) return
    const timestamp = Date.now()
    setAuditMessages((current) => [
      ...current,
      { id: `audit-user-${timestamp}`, role: "user", content: prompt },
      {
        id: `audit-assistant-${timestamp}`,
        role: "assistant",
        content: `I’ve prepared this request for ${project.brand_name}. Live findings will appear here when the audit workflow completes.`,
      },
    ])
    setAuditInput("")
    setMode("agent")
  }

  function sendContent() {
    const prompt = contentInput.trim()
    if (!prompt) return
    const timestamp = Date.now()
    setContentMessages((current) => [
      ...current,
      { id: `content-user-${timestamp}`, role: "user", content: prompt },
      {
        id: `content-assistant-${timestamp}`,
        role: "assistant",
        content: `Your ${contentType.label.toLowerCase()} request is ready. AEOlyzer will use your brand, market, language, and content-type guidance before drafting.`,
      },
    ])
    setContentInput("")
  }

  return (
    <main className="grid h-screen grid-cols-[188px_1fr] overflow-hidden bg-[#fff] font-sans text-[#1d1b18] max-md:grid-cols-[68px_1fr] dark:bg-[#2b2a27] dark:text-[#ececec]">
      <Sidebar
        project={project}
        mode={mode}
        activeTab={tab}
        auditMessages={auditMessages}
        contentMessages={contentMessages}
        theme={theme}
        onModeChange={setMode}
        onTabChange={(nextTab) => {
          setMode("dashboard")
          setTab(nextTab)
        }}
        onNewChat={() => {
          if (mode === "agent") {
            setAuditMessages([])
            setAuditInput("")
          } else {
            setContentMessages([])
            setContentInput("")
          }
        }}
        onThemeChange={onThemeChange}
        onReset={onReset}
      />

      <section className="min-w-0 overflow-y-auto">
        <header className="sticky top-0 z-30 flex h-10 items-center justify-between border-b border-[#ece6e1] bg-white px-3 dark:border-[#3b3a36] dark:bg-[#2b2a27]">
          <div className="flex h-7 w-44 items-center gap-2 rounded-lg border border-[#e7dfd9] px-2 text-[10px] text-[#8b817a] dark:border-[#4a4945] dark:text-[#a3a29e]">
            <Search size={13} />
            <span className="flex-1">Find anything...</span>
            <kbd className="text-[9px]">⌘ K</kbd>
          </div>
          <div className="flex items-center gap-2">
            <button className="flex h-7 items-center gap-1.5 rounded-lg border border-[#e7dfd9] px-3 text-[10px] text-[#6f625b] dark:border-[#4a4945] dark:text-[#d4d1ca]">
              <GraduationCap size={13} /> Learn
            </button>
          </div>
        </header>

        {mode === "dashboard" && (
          <Dashboard
            project={project}
            tab={tab}
            actions={dashboardActions}
            value={auditInput}
            onChange={setAuditInput}
            onSend={() => sendAudit(auditInput)}
            onAction={setAuditInput}
            onPrevious={() =>
              setDashboardActionPage((current) =>
                current === 0 ? Math.ceil(dashboardActionOrder.length / 4) - 1 : current - 1,
              )
            }
            onNext={() =>
              setDashboardActionPage(
                (current) =>
                  (current + 1) % Math.ceil(dashboardActionOrder.length / 4),
              )
            }
          />
        )}

        {mode === "agent" && (
          <AgentScreen
            project={project}
            messages={auditMessages}
            value={auditInput}
            actions={agentActions}
            onChange={setAuditInput}
            onSend={() => sendAudit(auditInput)}
            onAction={setAuditInput}
            onPrevious={() =>
              setAgentActionPage((current) =>
                current === 0 ? Math.ceil(agentActionOrder.length / 4) - 1 : current - 1,
              )
            }
            onNext={() =>
              setAgentActionPage(
                (current) => (current + 1) % Math.ceil(agentActionOrder.length / 4),
              )
            }
          />
        )}

        {mode === "content" && (
          <ContentScreen
            messages={contentMessages}
            value={contentInput}
            selectedType={contentType}
            onChange={setContentInput}
            onTypeChange={setContentType}
            onSend={sendContent}
          />
        )}
      </section>
    </main>
  )
}

function Sidebar({
  project,
  mode,
  activeTab,
  auditMessages,
  contentMessages,
  theme,
  onModeChange,
  onTabChange,
  onNewChat,
  onThemeChange,
  onReset,
}: {
  project: ProjectProfile
  mode: Mode
  activeTab: DashboardTab
  auditMessages: ChatMessage[]
  contentMessages: ChatMessage[]
  theme: ThemeMode
  onModeChange: (mode: Mode) => void
  onTabChange: (tab: DashboardTab) => void
  onNewChat: () => void
  onThemeChange: (theme: ThemeMode) => void
  onReset: () => void
}) {
  const [projectMenuOpen, setProjectMenuOpen] = useState(false)
  const history = mode === "agent" ? auditMessages : contentMessages

  return (
    <aside className="flex min-h-0 flex-col border-r border-[#e8e1dc] bg-[#faf7f5] p-1.5 dark:border-[#3a3936] dark:bg-[#252422]">
      <div className="relative">
        <button
          aria-expanded={projectMenuOpen}
          className={`${softCard} flex h-9 w-full items-center gap-2 px-2 text-left`}
          onClick={() => setProjectMenuOpen((open) => !open)}
        >
          <BrandMark project={project} />
          <span className="min-w-0 flex-1 max-md:hidden">
            <strong className="block truncate text-[10px]">{project.brand_name}</strong>
            <small className={`mt-0.5 block truncate text-[8px] ${muted}`}>
              {hostOnly(project.domain)}
            </small>
          </span>
          <ChevronDown size={12} className="max-md:hidden" />
        </button>
        {projectMenuOpen && (
          <div className="absolute left-0 right-0 top-10 z-50 rounded-xl border border-[#e6e1da] bg-white p-2 shadow-xl dark:border-[#4a4945] dark:bg-[#393836] max-md:hidden">
            <ThemeToggle theme={theme} onChange={onThemeChange} compact />
            <button
              className="mt-1.5 w-full rounded-lg px-2 py-1.5 text-left text-[9px] text-[#8b817a] hover:bg-[#f5efec] dark:text-[#aaa7a1] dark:hover:bg-[#4a4945]"
              onClick={onReset}
            >
              Reset guest workspace
            </button>
          </div>
        )}
      </div>

      <div className="mt-1 grid grid-cols-2 gap-1 max-md:grid-cols-1">
        <ModeButton
          active={mode === "dashboard"}
          label="Dashboard"
          onClick={() => onModeChange("dashboard")}
        />
        <ModeButton
          active={mode === "agent"}
          label="Agent"
          onClick={() => onModeChange("agent")}
        />
      </div>

      {mode === "dashboard" ? (
        <nav className="mt-2 min-h-0 overflow-y-auto [scrollbar-color:#ded8d2_transparent] [scrollbar-width:thin] dark:[scrollbar-color:#4a4945_transparent]">
          <NavSection label="Analytics">
            <NavButton
              active={activeTab === "aeo-insights"}
              icon={<Eye size={13} />}
              label="AEO Insights"
              onClick={() => onTabChange("aeo-insights")}
            />
            <NavButton
              active={false}
              icon={<Zap size={13} />}
              label="LLM Analytics"
              badge="Beta"
              onClick={() => onTabChange("aeo-insights")}
            />
            <NavButton
              active={activeTab === "traffic"}
              icon={<TrendingUp size={13} />}
              label="Traffic"
              onClick={() => onTabChange("traffic")}
            />
          </NavSection>
          <NavSection label="Prompts">
            <NavButton
              active={activeTab === "your-prompts"}
              icon={<MessageSquareText size={13} />}
              label="Your Prompts"
              onClick={() => onTabChange("your-prompts")}
            />
            <NavButton
              active={activeTab === "prompt-research"}
              icon={<GraduationCap size={13} />}
              label="Prompt Research"
              onClick={() => onTabChange("prompt-research")}
            />
          </NavSection>
          <NavSection label="Actions">
            <NavButton
              active={false}
              icon={<PenLine size={13} />}
              label="Content"
              onClick={() => onModeChange("content")}
            />
            <NavButton
              active={false}
              icon={<Lightbulb size={13} />}
              label="Opportunities"
              onClick={() => onTabChange("aeo-insights")}
            />
          </NavSection>
          <NavSection label="On Page">
            <NavButton
              active={activeTab === "site-health"}
              icon={<Globe2 size={13} />}
              label="Site Health"
              onClick={() => onTabChange("site-health")}
            />
          </NavSection>
        </nav>
      ) : (
        <div className="mt-2 min-h-0 overflow-y-auto [scrollbar-color:#ded8d2_transparent] [scrollbar-width:thin] dark:[scrollbar-color:#4a4945_transparent]">
          <button
            className={`${sidebarButton} max-md:justify-center`}
            onClick={onNewChat}
          >
            <Plus size={13} />
            <span className="max-md:hidden">New Chat</span>
          </button>
          <div className="mt-1 rounded-lg bg-[#f2ece8] dark:bg-[#393836]">
            <div className={`${sidebarButton} font-medium`}>
              <Clock3 size={13} />
              <span className="flex-1 max-md:hidden">History</span>
              <ChevronDown size={12} className="max-md:hidden" />
            </div>
          </div>
          <div className="mt-3 max-md:hidden">
            {history
              .filter((message) => message.role === "user")
              .slice(-6)
              .reverse()
              .map((message) => (
                <button
                  key={message.id}
                  className="w-full truncate rounded-lg px-8 py-1.5 text-left text-[9px] text-[#77716a] hover:bg-white dark:text-[#8f8d87] dark:hover:bg-[#393836]"
                >
                  {message.content}
                </button>
              ))}
            {history.length === 0 && (
              <div className="flex flex-col items-center pt-3 text-[#a69f98] dark:text-[#77756f]">
                <MessageSquareText size={14} />
                <small className="mt-2 text-[8px]">No conversations yet</small>
              </div>
            )}
          </div>
        </div>
      )}

      <div className="mt-auto border-t border-[#e5dfd8] pt-2 dark:border-[#3a3936]">
        <button className={`${sidebarButton} max-md:justify-center`}>
          <CircleHelp size={13} />
          <span className="max-md:hidden">Help</span>
        </button>
        <button className={`${sidebarButton} bg-[#eee7e3] max-md:justify-center dark:bg-[#393836]`}>
          <Gift size={13} />
          <span className="max-md:hidden">Invite Friends, Get 20%</span>
        </button>
        <div className="mt-1 flex items-center gap-2 rounded-lg border border-[#e8e1dc] bg-white p-1.5 dark:border-[#4a4945] dark:bg-[#2b2a27]">
          <span className="grid size-6 shrink-0 place-items-center rounded-full bg-[#34202b] text-[9px] font-semibold text-white">
            G
          </span>
          <span className="min-w-0 flex-1 max-md:hidden">
            <strong className="block truncate text-[9px]">Guest workspace</strong>
            <small className="block truncate text-[7px] text-[#9a938d]">
              Session only
            </small>
          </span>
          <ChevronDown size={11} className="max-md:hidden" />
        </div>
      </div>
    </aside>
  )
}

function Dashboard({
  project,
  tab,
  actions,
  value,
  onChange,
  onSend,
  onAction,
  onPrevious,
  onNext,
}: {
  project: ProjectProfile
  tab: DashboardTab
  actions: typeof quickActions
  value: string
  onChange: (value: string) => void
  onSend: () => void
  onAction: (prompt: string) => void
  onPrevious: () => void
  onNext: () => void
}) {
  if (tab !== "aeo-insights") {
    return <DashboardTabPage project={project} tab={tab} />
  }

  return (
    <div className="mx-auto max-w-[760px] px-8 pb-12 pt-11">
      <section className="text-center">
        <h1 className={`${heading} text-[31px]`}>
          {greetingForNow()} {project.brand_name}.
        </h1>
        <p className="mt-1.5 text-[13px] text-[#775247] dark:text-[#c79c89]">
          Want an update or have a question? Just chat below.
        </p>
      </section>

      <ReferenceComposer
        value={value}
        placeholder="What keywords should I target?"
        onChange={onChange}
        onSend={onSend}
      />
      <QuickActionGrid
        actions={actions}
        onAction={onAction}
        onPrevious={onPrevious}
        onNext={onNext}
      />
    </div>
  )
}

function DashboardTabPage({
  project,
  tab,
}: {
  project: ProjectProfile
  tab: Exclude<DashboardTab, "aeo-insights">
}) {
  const details = {
    traffic: {
      eyebrow: "SEARCH PERFORMANCE",
      title: "Connect your traffic data",
      description:
        "Bring Google Search Console data into AEOlyzer to connect search demand with AI visibility.",
      icon: <BarChart3 size={27} />,
    },
    "your-prompts": {
      eyebrow: "AI VISIBILITY",
      title: "Your tracked prompts",
      description:
        "These questions define how AEOlyzer will measure your brand in answer engines.",
      icon: <MessageCircle size={27} />,
    },
    "prompt-research": {
      eyebrow: "DISCOVERY",
      title: "Find the questions your audience asks",
      description:
        "Research new high-intent questions to expand your visibility tracking.",
      icon: <Search size={27} />,
    },
    "site-health": {
      eyebrow: "ON-PAGE",
      title: "Site health",
      description: "Technical checks and crawl findings will be organized here.",
      icon: <HeartPulse size={27} />,
    },
  }[tab]

  return (
    <div className="mx-auto max-w-[940px] px-10 py-16">
      <div className="flex items-start gap-4 border-b border-[#e8e3dc] pb-7 dark:border-[#3b3a36]">
        <span className="grid size-12 shrink-0 place-items-center rounded-xl bg-[#fff1ea] text-[#c96849] dark:bg-[#493a34] dark:text-[#e9916d]">
          {details.icon}
        </span>
        <span>
          <p className="mb-1 text-[10px] font-bold tracking-[0.15em] text-[#c96849] dark:text-[#e9916d]">
            {details.eyebrow}
          </p>
          <h1 className={`${heading} text-[38px]`}>{details.title}</h1>
          <p className={`mt-2 text-xs ${muted}`}>{details.description}</p>
        </span>
      </div>

      {tab === "your-prompts" && (
        <div className="mt-6 space-y-2">
          {project.prompts.map((prompt, index) => (
            <article
              key={prompt}
              className={`${softCard} grid grid-cols-[28px_1fr_auto] items-center gap-3 px-4 py-3`}
            >
              <span className="font-mono text-[10px] text-[#d97757]">
                {String(index + 1).padStart(2, "0")}
              </span>
              <p className="text-xs text-[#4f4943] dark:text-[#d0ccc5]">
                {prompt}
              </p>
              <CheckCircle2 size={16} className="text-[#5f9d70]" />
            </article>
          ))}
        </div>
      )}

      {tab === "traffic" && (
        <div className={`${softCard} mt-6 flex items-center gap-4 p-5`}>
          <div className="grid size-11 place-items-center rounded-lg bg-white text-xl font-black text-[#4285f4] shadow-sm">
            G
          </div>
          <span className="flex-1">
            <strong className="block text-sm">Google Search Console</strong>
            <small className={`mt-1 block text-[10px] ${muted}`}>
              Clicks, impressions, rankings and query performance
            </small>
          </span>
          <button className="rounded-lg bg-[#d97757] px-4 py-2 text-xs font-semibold text-white">
            Connect
          </button>
        </div>
      )}

      {tab === "prompt-research" && (
        <div className={`${softCard} mt-6 flex items-center gap-3 p-4`}>
          <Search size={19} className="text-[#8b857e]" />
          <input
            className="min-w-0 flex-1 bg-transparent text-sm outline-none"
            placeholder={`Research prompts for ${project.brand_name}…`}
          />
          <button className="rounded-lg bg-[#d97757] px-4 py-2 text-xs font-semibold text-white">
            Research
          </button>
        </div>
      )}

      {tab === "site-health" && (
        <div className="mt-6 grid grid-cols-2 gap-3">
          {["Crawlability", "Metadata", "Structured data", "AI discoverability"].map(
            (item) => (
              <article
                key={item}
                className={`${softCard} grid grid-cols-[36px_1fr] items-center gap-3 p-4`}
              >
                <span className="row-span-2 grid size-9 place-items-center rounded-full bg-[#edf7f0] text-[#5f9d70] dark:bg-[#34433a]">
                  <Check size={15} />
                </span>
                <strong className="text-xs">{item}</strong>
                <small className={`text-[9px] ${muted}`}>
                  Ready for first audit
                </small>
              </article>
            ),
          )}
        </div>
      )}
    </div>
  )
}

function AgentScreen({
  project,
  messages,
  value,
  actions,
  onChange,
  onSend,
  onAction,
  onPrevious,
  onNext,
}: {
  project: ProjectProfile
  messages: ChatMessage[]
  value: string
  actions: typeof quickActions
  onChange: (value: string) => void
  onSend: () => void
  onAction: (prompt: string) => void
  onPrevious: () => void
  onNext: () => void
}) {
  if (messages.length > 0) {
    return (
      <ChatLayout
        messages={messages}
        composer={
          <ReferenceComposer
            value={value}
            placeholder="Reply..."
            onChange={onChange}
            onSend={onSend}
            compact
          />
        }
      />
    )
  }

  return (
    <div className="mx-auto flex min-h-[calc(100vh-40px)] max-w-[760px] flex-col items-center px-8 pb-12 pt-11 text-center">
      <h1 className={`${heading} text-[31px]`}>
        {greetingForNow()} {project.brand_name}.
      </h1>
      <p className="mt-1.5 text-[13px] text-[#775247] dark:text-[#c79c89]">
        Want an update or have a question? Just chat below.
      </p>
      <ReferenceComposer
        value={value}
        placeholder="What's my top priority?"
        onChange={onChange}
        onSend={onSend}
      />
      <QuickActionGrid
        actions={actions}
        onAction={onAction}
        onPrevious={onPrevious}
        onNext={onNext}
      />
    </div>
  )
}

function ContentScreen({
  messages,
  value,
  selectedType,
  onChange,
  onTypeChange,
  onSend,
}: {
  messages: ChatMessage[]
  value: string
  selectedType: (typeof contentTypes)[number]
  onChange: (value: string) => void
  onTypeChange: (value: (typeof contentTypes)[number]) => void
  onSend: () => void
}) {
  if (messages.length > 0) {
    return (
      <ChatLayout
        messages={messages}
        composer={
          <Composer
            value={value}
            placeholder={`Describe your ${selectedType.label.toLowerCase()}…`}
            badge={selectedType.label}
            onChange={onChange}
            onSend={onSend}
          />
        }
      />
    )
  }

  return (
    <div className="flex min-h-[calc(100vh-48px)] flex-col items-center px-8 pb-12 pt-20">
      <span className="grid size-11 place-items-center rounded-xl bg-[#fff1ea] text-[#c96849] dark:bg-[#493a34] dark:text-[#e9916d]">
        <WandSparkles size={23} />
      </span>
      <p className="mb-2 mt-5 text-[10px] font-bold tracking-[0.16em] text-[#c96849] dark:text-[#e9916d]">
        CONTENT AGENT
      </p>
      <h1 className={`${heading} text-center text-[43px]`}>
        What can I help you create?
      </h1>
      <p className={`mb-6 mt-2 text-xs ${muted}`}>
        Start with an idea. AEOlyzer will build the brief before drafting.
      </p>

      <div className="w-full max-w-[720px] rounded-2xl border border-[#ddd7cf] bg-white p-4 shadow-[0_15px_45px_rgba(31,27,22,0.07)] dark:border-[#4a4945] dark:bg-[#393836] dark:shadow-[0_15px_45px_rgba(0,0,0,0.2)]">
        <textarea
          aria-label="Describe what you want to write"
          className="min-h-24 w-full resize-none bg-transparent text-sm outline-none placeholder:text-[#aaa49d]"
          placeholder="Describe what you want to write…"
          value={value}
          onChange={(event) => onChange(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter" && !event.shiftKey) {
              event.preventDefault()
              onSend()
            }
          }}
        />
        <div className="flex items-center justify-between">
          <span className={`flex items-center gap-2 text-[10px] ${muted}`}>
            <FileText size={14} /> {selectedType.label}
          </span>
          <button
            aria-label="Create content"
            className="grid size-10 place-items-center rounded-xl bg-[#d97757] text-white disabled:opacity-35"
            disabled={!value.trim()}
            onClick={onSend}
          >
            <ArrowRight size={17} />
          </button>
        </div>
      </div>

      <div className="mt-3 flex max-w-[760px] flex-wrap justify-center gap-2">
        {contentTypes.map((type) => (
          <button
            key={type.id}
            className={[
              "group relative flex h-9 items-center gap-1.5 rounded-lg border px-3 text-[10px] transition",
              selectedType.id === type.id
                ? "border-[#d97757] bg-[#fff1ea] text-[#a95034] dark:bg-[#493a34] dark:text-[#efad92]"
                : "border-[#e2ddd6] bg-white text-[#6f6962] hover:border-[#d9a38e] dark:border-[#4a4945] dark:bg-[#393836] dark:text-[#a3a29e]",
            ].join(" ")}
            onClick={() => onTypeChange(type)}
          >
            <FileText size={14} />
            {type.label}
            <span className="pointer-events-none absolute bottom-[calc(100%+9px)] left-1/2 z-20 hidden w-64 -translate-x-1/2 rounded-xl border border-[#ded8d0] bg-white p-3 text-left shadow-xl group-hover:block group-focus:block dark:border-[#55524c] dark:bg-[#393836]">
              <strong className="block text-[10px] leading-4 text-[#1d1b18] dark:text-[#ececec]">
                {type.description}
              </strong>
              <small className="mt-2 block border-t border-[#ebe6df] pt-2 text-[9px] leading-4 text-[#77716a] dark:border-[#4a4945] dark:text-[#a3a29e]">
                {type.details}
              </small>
            </span>
          </button>
        ))}
      </div>
    </div>
  )
}

function ChatLayout({
  messages,
  composer,
}: {
  messages: ChatMessage[]
  composer: React.ReactNode
}) {
  return (
    <div className="mx-auto flex min-h-[calc(100vh-48px)] max-w-[880px] flex-col px-8 pb-6 pt-8">
      <div className="flex-1 space-y-8 overflow-y-auto pb-10">
        {messages.map((message) =>
          message.role === "user" ? (
            <article
              key={message.id}
              className="ml-auto max-w-[72%] rounded-2xl bg-[#f0ede9] px-5 py-4 text-sm leading-6 text-[#292622] dark:bg-[#393836] dark:text-[#ececec]"
            >
              {message.content}
            </article>
          ) : (
            <article
              key={message.id}
              className="flex max-w-[78%] items-start gap-3 text-sm leading-6 text-[#3e3934] dark:text-[#ececec]"
            >
              <AeolyzerLogo size={27} />
              <p className="m-0 pt-0.5">{message.content}</p>
            </article>
          ),
        )}
      </div>
      <div className="sticky bottom-0 bg-[#fbfaf8] py-3 dark:bg-[#2b2a27]">
        {composer}
        <p className={`mt-2 text-center text-[9px] ${muted}`}>
          AEOlyzer can make mistakes. Verify important findings.
        </p>
      </div>
    </div>
  )
}

function ReferenceComposer({
  value,
  placeholder,
  compact = false,
  onChange,
  onSend,
}: {
  value: string
  placeholder: string
  compact?: boolean
  onChange: (value: string) => void
  onSend: () => void
}) {
  return (
    <div
      className={[
        "mx-auto w-full overflow-hidden rounded-[23px] border border-[#eee8e4] bg-white shadow-[0_8px_24px_rgba(74,52,43,0.04)]",
        "dark:border-[#4a4945] dark:bg-[#393836] dark:shadow-[0_8px_24px_rgba(0,0,0,0.18)]",
        compact ? "max-w-none" : "mt-[26px] max-w-[550px]",
      ].join(" ")}
    >
      {!compact && (
        <div className="flex h-9 items-center justify-between px-3 text-[10px] text-[#755247] dark:text-[#c79c89]">
          <span className="flex items-center gap-2">
            <Database size={12} />
            Get better answers and context by connecting your data
          </span>
          <span className="flex items-center gap-1.5" aria-hidden="true">
            <span className="font-semibold text-[#ed9b37]">▥</span>
            <span className="font-semibold text-[#4c9bea]">◢</span>
            <span className="ml-1 text-[#a49b95]">×</span>
          </span>
        </div>
      )}
      <div className={compact ? "" : "border-t border-[#f0ebe7]"}>
        <textarea
          aria-label={placeholder}
          className={[
            "block w-full resize-none bg-transparent px-4 pt-3 text-[12px] leading-5 outline-none placeholder:text-[#a59d97]",
            "dark:placeholder:text-[#77756f]",
            compact ? "min-h-14" : "min-h-11",
          ].join(" ")}
          rows={compact ? 2 : 1}
          placeholder={placeholder}
          value={value}
          onChange={(event) => onChange(event.target.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter" && !event.shiftKey) {
              event.preventDefault()
              onSend()
            }
          }}
        />
        <div className="flex h-10 items-center gap-2 px-2 pb-2">
          <button
            aria-label="Attach file"
            className="grid size-7 place-items-center rounded-lg border border-[#e8dfd9] text-[#74594e] hover:bg-[#faf6f3] dark:border-[#55514c] dark:text-[#d0b3a7] dark:hover:bg-[#44413d]"
          >
            <Paperclip size={13} />
          </button>
          <button className="flex h-7 items-center gap-1.5 rounded-lg border border-[#e8dfd9] px-2 text-[10px] text-[#74594e] hover:bg-[#faf6f3] dark:border-[#55514c] dark:text-[#d0b3a7] dark:hover:bg-[#44413d]">
            <Command size={12} />
            Shortcuts
          </button>
          <span className="flex items-center gap-1 text-[9px] text-[#74594e] dark:text-[#d0b3a7]">
            Connectors
            <span className="grid size-3.5 place-items-center rounded-full border border-current text-[7px] font-bold">
              G
            </span>
            <span className="grid size-3.5 place-items-center rounded-sm bg-[#171512] text-[7px] font-bold text-white dark:bg-[#ececec] dark:text-[#2b2a27]">
              N
            </span>
          </span>
          <button
            aria-label="Send message"
            className="ml-auto grid size-8 place-items-center rounded-full bg-[#eee5df] text-[#b9a9a0] transition hover:bg-[#e7d8d0] disabled:opacity-70 dark:bg-[#4b4844] dark:text-[#817d77]"
            disabled={!value.trim()}
            onClick={onSend}
          >
            <ArrowRight size={16} />
          </button>
        </div>
      </div>
    </div>
  )
}

function QuickActionGrid({
  actions,
  onAction,
  onPrevious,
  onNext,
}: {
  actions: typeof quickActions
  onAction: (prompt: string) => void
  onPrevious: () => void
  onNext: () => void
}) {
  return (
    <div className="group relative mx-auto mt-4 w-full max-w-[544px]">
      <div className="grid grid-cols-4 gap-2.5">
        {actions.map((action) => (
          <button
            key={action.title}
            className="flex h-20 flex-col items-start justify-between rounded-xl border border-[#eee8e4] bg-white p-2.5 text-left text-[#755247] transition hover:border-[#dccbc2] hover:shadow-sm dark:border-[#4a4945] dark:bg-[#393836] dark:text-[#d0b3a7] dark:hover:border-[#645f59]"
            onClick={() => onAction(action.prompt)}
          >
            {actionIcon(action.title)}
            <span className="line-clamp-2 text-[10px] font-medium leading-4">
              {sentenceCaseAction(action.title)}
            </span>
          </button>
        ))}
      </div>
      <button
        aria-label="Previous suggested actions"
        className="absolute -left-3 top-1/2 grid size-6 -translate-y-1/2 place-items-center rounded-full border border-[#e6dfda] bg-white text-[#8d7c73] opacity-0 shadow-sm transition group-hover:opacity-100 focus:opacity-100 dark:border-[#504d48] dark:bg-[#393836]"
        onClick={onPrevious}
      >
        <ChevronLeft size={12} />
      </button>
      <button
        aria-label="Next suggested actions"
        className="absolute -right-3 top-1/2 grid size-6 -translate-y-1/2 place-items-center rounded-full border border-[#e6dfda] bg-white text-[#8d7c73] opacity-0 shadow-sm transition group-hover:opacity-100 focus:opacity-100 dark:border-[#504d48] dark:bg-[#393836]"
        onClick={onNext}
      >
        <ChevronRight size={12} />
      </button>
    </div>
  )
}

function Composer({
  value,
  placeholder,
  badge,
  large = false,
  onChange,
  onSend,
}: {
  value: string
  placeholder: string
  badge?: string
  large?: boolean
  onChange: (value: string) => void
  onSend: () => void
}) {
  return (
    <div
      className={[
        "flex w-full items-end gap-3 rounded-2xl border border-[#ddd7cf] bg-white p-3 shadow-[0_12px_35px_rgba(31,27,22,0.07)] focus-within:border-[#d97757]",
        "dark:border-[#4a4945] dark:bg-[#393836] dark:shadow-[0_12px_35px_rgba(0,0,0,0.2)]",
        large ? "min-h-[86px]" : "",
      ].join(" ")}
    >
      {badge && (
        <span className="mb-1 rounded-lg bg-[#fff1ea] px-2 py-1.5 text-[9px] text-[#a95034] dark:bg-[#4a4945] dark:text-[#e9916d]">
          {badge}
        </span>
      )}
      <textarea
        aria-label={placeholder}
        className="min-h-10 min-w-0 flex-1 resize-none bg-transparent px-1 py-2 text-sm outline-none placeholder:text-[#aaa49d] dark:placeholder:text-[#77756f]"
        rows={large ? 2 : 1}
        placeholder={placeholder}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        onKeyDown={(event) => {
          if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault()
            onSend()
          }
        }}
      />
      <button
        aria-label="Send message"
        className="grid size-10 shrink-0 place-items-center rounded-xl bg-[#d97757] text-white transition hover:bg-[#c96849] disabled:opacity-35 dark:bg-[#e07b53] dark:text-[#241d19]"
        disabled={!value.trim()}
        onClick={onSend}
      >
        <Send size={16} />
      </button>
    </div>
  )
}

function ModeButton({
  active,
  label,
  onClick,
}: {
  active: boolean
  label: string
  onClick: () => void
}) {
  return (
    <button
      className={[
        "flex h-7 items-center justify-center rounded-lg text-[9px] font-medium transition",
        active
          ? "bg-white text-[#332b27] shadow-sm ring-1 ring-[#ebe5e1] dark:bg-[#393836] dark:text-[#ececec] dark:ring-[#4a4945]"
          : "bg-[#f5efec] text-[#806e65] dark:bg-[#2d2c29] dark:text-[#a3a29e]",
      ].join(" ")}
      onClick={onClick}
    >
      <span className="max-md:hidden">{label}</span>
    </button>
  )
}

function NavSection({
  label,
  children,
}: {
  label: string
  children: React.ReactNode
}) {
  return (
    <section className="mt-1 border-b border-[#eee7e3] pb-2 first:mt-0 last:-mt-1 last:border-b-0 last:pb-0 dark:border-[#393835]">
      <p className="mb-1 px-2 text-[9px] font-bold tracking-[0.14em] text-[#a19b94] max-md:hidden dark:text-[#6f6d68]">
        {label}
      </p>
      {children}
    </section>
  )
}

function NavButton({
  active,
  icon,
  label,
  badge,
  onClick,
}: {
  active: boolean
  icon: React.ReactNode
  label: string
  badge?: string
  onClick: () => void
}) {
  return (
    <button
      className={[
        sidebarButton,
        active
          ? "bg-white text-[#1e1b18] shadow-sm ring-1 ring-[#e8e3dc] dark:bg-[#393836] dark:text-[#ececec] dark:ring-[#4a4945]"
          : "",
        "max-md:justify-center max-md:px-0",
      ].join(" ")}
      onClick={onClick}
    >
      <span className={active ? "text-[#d97757]" : ""}>{icon}</span>
      <span className="flex-1 max-md:hidden">{label}</span>
      {badge && (
        <small className="rounded-full bg-[#ece8e3] px-1.5 py-0.5 text-[9px] max-md:hidden dark:bg-[#4a4945]">
          {badge}
        </small>
      )}
    </button>
  )
}

function BrandMark({ project }: { project: ProjectProfile }) {
  const [failed, setFailed] = useState(false)
  return (
    <span className="grid size-8 shrink-0 place-items-center overflow-hidden rounded-lg bg-[#fff1ea] text-xs font-bold text-[#c96849] dark:bg-[#4a4945] dark:text-[#e9916d]">
      {project.icon_url && !failed ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          className="size-full object-contain"
          src={project.icon_url}
          alt=""
          onError={() => setFailed(true)}
        />
      ) : (
        project.brand_name.charAt(0).toUpperCase()
      )}
    </span>
  )
}

function actionIcon(title: string) {
  const iconClass = "text-[#8d6f63] dark:text-[#d0b3a7]"
  if (title === "Audit my Homepage") {
    return <Stethoscope size={16} className={iconClass} />
  }
  if (title === "Quick Visibility Check") {
    return <Eye size={16} className={iconClass} />
  }
  if (title === "Analyze competitors") {
    return <BarChart3 size={16} className={iconClass} />
  }
  if (title === "Make me a Task List") {
    return <ListChecks size={16} className={iconClass} />
  }
  if (title === "Optimize Meta tags") {
    return <Zap size={16} className={iconClass} />
  }
  if (title === "Review Technical SEO") {
    return <Stethoscope size={16} className={iconClass} />
  }
  if (title === "Find Keyword opportunities") {
    return <Target size={16} className={iconClass} />
  }
  if (title === "Find content gaps") {
    return <Search size={16} className={iconClass} />
  }
  return <FileText size={16} className={iconClass} />
}

function actionPage(order: string[], page: number) {
  const byTitle = new Map(quickActions.map((action) => [action.title, action]))
  const ordered = order
    .map((title) => byTitle.get(title))
    .filter((action): action is (typeof quickActions)[number] => Boolean(action))
  const start = (page * 4) % ordered.length

  return Array.from(
    { length: Math.min(4, ordered.length) },
    (_, index) => ordered[(start + index) % ordered.length],
  )
}

function sentenceCaseAction(title: string) {
  const labels: Record<string, string> = {
    "Audit my Homepage": "Audit my homepage",
    "Quick Visibility Check": "Quick visibility check",
    "Make me a Task List": "Make me a task list",
    "Optimize Meta tags": "Optimize meta tags",
    "Review Technical SEO": "Review technical SEO",
    "Summarize recent Performance": "Summarize recent performance",
    "Find Keyword opportunities": "Find keyword opportunities",
    "Draft a Blog Post": "Draft a blog post",
  }
  return labels[title] ?? title
}

function greetingForNow() {
  const hour = new Date().getHours()
  if (hour < 12) return "Good morning"
  if (hour < 18) return "Good afternoon"
  return "Good evening"
}

function readHistory(key: string): ChatMessage[] {
  if (typeof window === "undefined") return []
  const stored = sessionStorage.getItem(key)
  if (!stored) return []
  try {
    return JSON.parse(stored) as ChatMessage[]
  } catch {
    return []
  }
}

function hostOnly(value: string) {
  try {
    return new URL(value).hostname.replace(/^www\./, "")
  } catch {
    return value
  }
}
