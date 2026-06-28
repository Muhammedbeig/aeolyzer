"use client"

import {
  ArrowLeft,
  ArrowRight,
  Building2,
  Check,
  CheckCircle2,
  ChevronDown,
  CircleGauge,
  Globe2,
  LoaderCircle,
  MapPin,
  Plus,
  Search,
  Sparkles,
  Target,
  Users,
  WandSparkles,
  X,
} from "lucide-react"
import { useEffect, useMemo, useState } from "react"
import { ThemeToggle, type ThemeMode } from "@/components/aeolyzer-app"
import { AeolyzerLogo } from "@/components/aeolyzer-logo"
import {
  type AccountType,
  type CountryOption,
  type ProjectProfile,
  type Reach,
  completeOnboarding,
  countries,
  fallbackBrandName,
  inspectSite,
  languages,
  reachOptions,
} from "@/lib/aeolyzer"

type Stage =
  | "account"
  | "domain"
  | "brand"
  | "market"
  | "scan"
  | "competitors"
  | "prompts"
  | "ready"

const panel =
  "rounded-2xl border border-[#e8e3dc] bg-white shadow-[0_18px_60px_rgba(31,27,22,0.07)] dark:border-[#4a4945] dark:bg-[#393836] dark:shadow-[0_20px_70px_rgba(0,0,0,0.28)]"
const field =
  "flex h-12 w-full items-center gap-3 rounded-xl border border-[#ddd7cf] bg-white px-4 text-sm text-[#211f1c] outline-none transition focus-within:border-[#d97757] focus-within:ring-4 focus-within:ring-[#d97757]/10 dark:border-[#4a4945] dark:bg-[#393836] dark:text-[#ececec]"
const primary =
  "inline-flex h-12 items-center justify-center gap-2 rounded-xl bg-[#d97757] px-5 text-sm font-semibold text-white transition hover:bg-[#c96849] disabled:cursor-not-allowed disabled:opacity-40 dark:bg-[#e07b53] dark:text-[#241d19] dark:hover:bg-[#ef8a62]"
const eyebrow =
  "mb-3 text-[11px] font-bold tracking-[0.16em] text-[#c96849] dark:text-[#e9916d]"
const heading =
  "font-display text-[42px] font-medium leading-[1.04] tracking-[-0.035em] text-[#171512] dark:text-[#f1eee8]"
const muted = "text-[#77716a] dark:text-[#aaa7a1]"

const initialForm = {
  accountType: "" as AccountType | "",
  domain: "",
  canonicalDomain: "",
  brandName: "",
  suggestedBrandName: "",
  reach: "" as Reach | "",
  countryCode: "",
  countryName: "",
  language: "",
  competitors: [] as string[],
  iconURL: "",
  description: "",
  category: "",
  competitorCandidates: [] as string[],
}

export function OnboardingFlow({
  theme,
  onThemeChange,
  onComplete,
}: {
  theme: ThemeMode
  onThemeChange: (theme: ThemeMode) => void
  onComplete: (project: ProjectProfile) => void
}) {
  // Memoize static dictionaries to prevent repeated execution during state-driven UI reconciliations.
  const allCountries = useMemo(() => countries(), [])
  const allLanguages = useMemo(() => languages(), [])
  const [stage, setStage] = useState<Stage>("account")
  const [form, setForm] = useState(initialForm)
  const [sessionID] = useState(() => crypto.randomUUID())
  const [busy, setBusy] = useState(false)
  const [notice, setNotice] = useState("")
  const [competitorInput, setCompetitorInput] = useState("")
  const [prompts, setPrompts] = useState<string[]>([])

  // Orchestrate sequential UI transitions.
  // Cleans up timeouts on unmount or dependency change to avoid memory leaks and ghost state updates.
  useEffect(() => {
    if (stage !== "scan") return
    const timer = window.setTimeout(() => setStage("competitors"), 1700)
    return () => window.clearTimeout(timer)
  }, [stage])

  async function readWebsite() {
    setBusy(true)
    setNotice("")
    try {
      const inspection = await inspectSite(sessionID, form.domain)
      const suggestion =
        inspection.suggested_brand_name || fallbackBrandName(form.domain)
      setForm((current) => ({
        ...current,
        canonicalDomain: inspection.canonical_url,
        suggestedBrandName: suggestion,
        brandName: current.brandName || suggestion,
        iconURL: inspection.icon_url,
        description: inspection.description,
        category: inspection.category,
        competitorCandidates: inspection.competitor_candidates,
      }))
    } catch (error) {
      const suggestion = fallbackBrandName(form.domain)
      setForm((current) => ({
        ...current,
        canonicalDomain: normalizeURL(current.domain),
        suggestedBrandName: suggestion,
        brandName: current.brandName || suggestion,
      }))
      setNotice(
        error instanceof Error
          ? error.message
          : "We could not read the site. Continue with the details manually.",
      )
    } finally {
      setBusy(false)
      setStage("brand")
    }
  }

  function selectCountry(country: CountryOption) {
    const language =
      allLanguages.find((item) => item.code === country.languageCode) ??
      allLanguages.find((item) => item.code === "en")
    setForm((current) => ({
      ...current,
      countryCode: country.code,
      countryName: country.name,
      language: language?.name ?? "",
    }))
  }

  function addCompetitor(raw: string) {
    const value = raw.trim().replace(/^https?:\/\//, "").replace(/\/.*$/, "")
    if (!value || form.competitors.includes(value) || form.competitors.length >= 5) {
      return
    }
    setForm((current) => ({
      ...current,
      competitors: [...current.competitors, value],
    }))
    setCompetitorInput("")
  }

  async function createPrompts() {
    if (!form.accountType || !form.reach) return
    setBusy(true)
    setNotice("")
    try {
      const frame = await completeOnboarding({
        session_id: sessionID,
        account_type: form.accountType,
        domain: form.canonicalDomain || normalizeURL(form.domain),
        brand_name: form.brandName,
        reach: form.reach,
        country_code: form.countryCode,
        country_name: form.countryName,
        language: form.language,
        competitors: form.competitors,
        icon_url: form.iconURL,
        description: form.description,
        category: form.category,
      })
      setPrompts(frame.prompts)
      setStage("prompts")
    } catch (error) {
      setNotice(
        error instanceof Error
          ? error.message
          : "AEOlyzer could not create the prompts.",
      )
    } finally {
      setBusy(false)
    }
  }

  function finish() {
    if (!form.accountType || !form.reach) return
    const project: ProjectProfile = {
      session_id: sessionID,
      account_type: form.accountType,
      domain: form.canonicalDomain || normalizeURL(form.domain),
      brand_name: form.brandName,
      reach: form.reach,
      country_code: form.countryCode,
      country_name: form.countryName,
      language: form.language,
      competitors: form.competitors,
      icon_url: form.iconURL,
      description: form.description,
      category: form.category,
      prompts,
    }
    sessionStorage.setItem("aeolyzer.pending-project", JSON.stringify(project))
    setStage("ready")
  }

  if (stage === "account") {
    return (
      <main className="relative grid min-h-screen place-items-center overflow-hidden bg-[#fbfaf8] px-6 text-[#1d1b18] dark:bg-[#2b2a27] dark:text-[#ececec]">
        <Header theme={theme} onThemeChange={onThemeChange} />
        <div className="pointer-events-none absolute size-[680px] rounded-full border border-[#ece7e0] dark:border-white/[0.035]" />
        <div className="pointer-events-none absolute size-[850px] rounded-full border border-[#f1ede8] dark:border-white/[0.025]" />

        <section className="relative z-10 flex w-full max-w-[720px] flex-col items-center py-28 text-center">
          <p className={eyebrow}>LET&apos;S PERSONALIZE YOUR WORKSPACE</p>
          <h1 className={`${heading} text-[54px]`}>Are you a brand or an agency?</h1>
          <p className={`mt-4 text-sm ${muted}`}>
            We&apos;ll tailor your visibility insights and recommendations.
          </p>
          <div className="mt-9 grid w-full grid-cols-2 gap-4 max-sm:grid-cols-1">
            <AccountCard
              active={form.accountType === "brand"}
              icon={<Building2 size={22} />}
              title="I'm a brand"
              description="Track and improve one brand's visibility."
              onClick={() =>
                setForm((current) => ({ ...current, accountType: "brand" }))
              }
            />
            <AccountCard
              active={form.accountType === "agency"}
              icon={<Users size={22} />}
              title="I'm an agency"
              description="Manage visibility for client brands."
              onClick={() =>
                setForm((current) => ({ ...current, accountType: "agency" }))
              }
            />
          </div>
          {form.accountType && (
            <button className={`${primary} mt-6 min-w-40`} onClick={() => setStage("domain")}>
              Continue <ArrowRight size={17} />
            </button>
          )}
        </section>
        <p className={`absolute bottom-7 text-xs ${muted}`}>
          No account needed · Guest session only
        </p>
      </main>
    )
  }

  if (stage === "scan") {
    return (
      <main className="relative grid min-h-screen place-items-center bg-[#fbfaf8] text-[#1d1b18] dark:bg-[#252421] dark:text-[#ececec]">
        <Header theme={theme} onThemeChange={onThemeChange} />
        <section className="flex flex-col items-center text-center">
          <div className="relative grid size-40 place-items-center rounded-full border border-dashed border-[#d97757]/50">
            <div className="absolute inset-6 animate-spin rounded-full border border-dashed border-[#d97757]/35 [animation-duration:4s]" />
            <span className="absolute left-0 top-1/2 size-2 rounded-full bg-[#d97757]" />
            <span className="absolute right-5 top-5 size-2 animate-pulse rounded-full bg-[#d97757]" />
            <AeolyzerLogo size={50} animate />
          </div>
          <h1 className={`${heading} mt-8`}>Learning about {form.brandName}</h1>
          <p className={`mt-3 text-sm ${muted}`}>
            Scanning public brand signals and preparing your context.
          </p>
          <div className="mt-7 flex gap-2 text-xs">
            <span className="flex items-center gap-1.5 rounded-full border border-[#dfe9e1] bg-[#f4faf5] px-3 py-2 text-[#4b8a5f] dark:border-[#45604d] dark:bg-[#34433a] dark:text-[#8ec49d]">
              <Check size={14} /> Brand identity
            </span>
            <span className={`flex items-center gap-1.5 rounded-full border border-[#e7e2db] px-3 py-2 dark:border-[#4a4945] ${muted}`}>
              <LoaderCircle size={14} className="animate-spin" /> Market context
            </span>
            <span className={`rounded-full border border-[#e7e2db] px-3 py-2 dark:border-[#4a4945] ${muted}`}>
              Competitor signals
            </span>
          </div>
        </section>
      </main>
    )
  }

  if (stage === "ready") {
    return (
      <main className="relative grid min-h-screen place-items-center bg-[#fbfaf8] text-[#1d1b18] dark:bg-[#252421] dark:text-[#ececec]">
        <Header theme={theme} onThemeChange={onThemeChange} />
        <section className="flex max-w-xl flex-col items-center px-6 text-center">
          <div className="grid size-20 place-items-center rounded-full border border-[#b9dac2] bg-[#f2faf4] text-[#4b8a5f] dark:border-[#45604d] dark:bg-[#34433a] dark:text-[#8ec49d]">
            <CheckCircle2 size={36} />
          </div>
          <p className={`${eyebrow} mt-7`}>SETUP COMPLETE</p>
          <h1 className={heading}>Your first AEO report is ready</h1>
          <p className={`mt-4 text-sm leading-6 ${muted}`}>
            AEOlyzer prepared your brand workspace and {prompts.length} AI
            visibility prompts.
          </p>
          <button
            className={`${primary} mt-7`}
            onClick={() => {
              const stored = sessionStorage.getItem("aeolyzer.pending-project")
              if (!stored) return
              sessionStorage.removeItem("aeolyzer.pending-project")
              onComplete(JSON.parse(stored) as ProjectProfile)
            }}
          >
            Open dashboard <ArrowRight size={17} />
          </button>
        </section>
      </main>
    )
  }

  return (
    <main className="grid min-h-screen grid-cols-[46%_54%] bg-white text-[#1d1b18] max-lg:grid-cols-1 dark:bg-[#2b2a27] dark:text-[#ececec]">
      <section className="relative min-h-screen border-r border-[#ece7e0] bg-[#fff] dark:border-[#3b3a36] dark:bg-[#2b2a27]">
        <Header theme={theme} onThemeChange={onThemeChange} compact />
        <div className="absolute left-8 top-20 flex items-center gap-3 text-xs text-[#8b857e] dark:text-[#8f8d87]">
          <button
            aria-label="Go back"
            className="grid size-9 place-items-center rounded-lg border border-[#e3ded7] bg-white transition hover:bg-[#f7f4f0] dark:border-[#4a4945] dark:bg-[#393836] dark:hover:bg-[#44433f]"
            onClick={() => setStage(previousStage(stage))}
          >
            <ArrowLeft size={17} />
          </button>
          {stageLabel(stage)}
        </div>

        <div className="mx-auto max-w-[530px] px-10 pb-14 pt-36 max-sm:px-6">
          {stage === "domain" && (
            <>
              <p className={eyebrow}>YOUR WEBSITE</p>
              <h1 className={heading}>Where can we find your brand?</h1>
              <p className={`mb-8 mt-4 text-sm leading-6 ${muted}`}>
                We&apos;ll use the public site to understand your positioning.
              </p>
              <FieldLabel htmlFor="website">Website URL</FieldLabel>
              <div className={field}>
                <Globe2 size={18} className="shrink-0 text-[#8a847d]" />
                <input
                  id="website"
                  autoFocus
                  className="min-w-0 flex-1 bg-transparent outline-none placeholder:text-[#aaa49d]"
                  value={form.domain}
                  placeholder="yourbrand.com"
                  onChange={(event) =>
                    setForm((current) => ({ ...current, domain: event.target.value }))
                  }
                  onKeyDown={(event) => {
                    if (event.key === "Enter" && plausibleDomain(form.domain)) {
                      void readWebsite()
                    }
                  }}
                />
                {plausibleDomain(form.domain) && (
                  <CheckCircle2 size={17} className="text-[#5f9d70]" />
                )}
              </div>
              <button
                className={`${primary} mt-6 w-full`}
                disabled={!plausibleDomain(form.domain) || busy}
                onClick={() => void readWebsite()}
              >
                {busy && <LoaderCircle size={16} className="animate-spin" />}
                {busy ? "Reading website…" : "Continue"}
              </button>
            </>
          )}

          {stage === "brand" && (
            <>
              <p className={eyebrow}>BRAND IDENTITY</p>
              <h1 className={heading}>What should we call your brand?</h1>
              <p className={`mb-8 mt-4 text-sm leading-6 ${muted}`}>
                This name will appear across your reports and recommendations.
              </p>
              {notice && <Notice>{notice}</Notice>}
              <FieldLabel htmlFor="brand">Brand name</FieldLabel>
              <div className={field}>
                <SiteImage
                  src={form.iconURL}
                  className="size-5 rounded object-contain"
                  fallback={<Sparkles size={18} className="text-[#d97757]" />}
                />
                <input
                  id="brand"
                  autoFocus
                  className="min-w-0 flex-1 bg-transparent outline-none placeholder:text-[#aaa49d]"
                  value={form.brandName}
                  placeholder="Your brand name"
                  onChange={(event) =>
                    setForm((current) => ({
                      ...current,
                      brandName: event.target.value,
                    }))
                  }
                />
              </div>
              {form.suggestedBrandName &&
                form.brandName !== form.suggestedBrandName && (
                  <button
                    className="mt-3 flex items-center gap-2 text-xs font-medium text-[#c96849] dark:text-[#e9916d]"
                    onClick={() =>
                      setForm((current) => ({
                        ...current,
                        brandName: current.suggestedBrandName,
                      }))
                    }
                  >
                    <WandSparkles size={14} />
                    Suggested: {form.suggestedBrandName}
                  </button>
                )}
              <button
                className={`${primary} mt-6 w-full`}
                disabled={!form.brandName.trim()}
                onClick={() => setStage("market")}
              >
                Continue
              </button>
            </>
          )}

          {stage === "market" && (
            <>
              <p className={eyebrow}>YOUR MARKET</p>
              <h1 className={heading}>Where does your brand reach?</h1>
              <p className={`mb-6 mt-4 text-sm leading-6 ${muted}`}>
                Location and language make every insight more relevant.
              </p>
              <FieldLabel>Brand reach</FieldLabel>
              <div className="grid grid-cols-2 gap-2">
                {reachOptions.map((option, index) => (
                  <button
                    key={option.id}
                    className={[
                      "relative rounded-xl border p-3 text-left transition",
                      index === 0 ? "col-span-2" : "",
                      form.reach === option.id
                        ? "border-[#d97757] bg-[#fff7f3] ring-2 ring-[#d97757]/10 dark:bg-[#493a34]"
                        : "border-[#e4dfd8] bg-white hover:border-[#cfc7bd] dark:border-[#4a4945] dark:bg-[#393836] dark:hover:border-[#66635d]",
                    ].join(" ")}
                    onClick={() =>
                      setForm((current) => ({ ...current, reach: option.id }))
                    }
                  >
                    <strong className="block text-xs font-semibold">
                      {option.title}
                    </strong>
                    <small className={`mt-1 block text-[10px] leading-4 ${muted}`}>
                      {option.description}
                    </small>
                    {form.reach === option.id && (
                      <Check
                        size={14}
                        className="absolute right-3 top-3 text-[#d97757]"
                      />
                    )}
                  </button>
                ))}
              </div>
              <div className="mt-4 grid grid-cols-2 gap-3">
                <SearchSelect
                  label="Primary location"
                  placeholder="Search countries"
                  icon={<MapPin size={16} />}
                  value={form.countryName}
                  options={allCountries}
                  optionKey={(option) => option.code}
                  optionLabel={(option) => option.name}
                  onSelect={selectCountry}
                />
                <SearchSelect
                  label="Workspace language"
                  placeholder="Search languages"
                  icon={<Globe2 size={16} />}
                  value={form.language}
                  options={allLanguages}
                  optionKey={(option) => option.code}
                  optionLabel={(option) => option.name}
                  onSelect={(language) =>
                    setForm((current) => ({
                      ...current,
                      language: language.name,
                    }))
                  }
                />
              </div>
              <button
                className={`${primary} mt-6 w-full`}
                disabled={!form.reach || !form.countryCode || !form.language}
                onClick={() => setStage("scan")}
              >
                Scan my website <Sparkles size={16} />
              </button>
            </>
          )}

          {stage === "competitors" && (
            <>
              <p className={eyebrow}>COMPETITIVE LANDSCAPE</p>
              <h1 className={heading}>Who do customers compare you with?</h1>
              <p className={`mb-7 mt-4 text-sm leading-6 ${muted}`}>
                Add up to five relevant competitors. You can change them later.
              </p>
              {form.competitorCandidates.length > 0 ? (
                <>
                  <FieldLabel>Suggested from your site</FieldLabel>
                  <div className="flex flex-wrap gap-2">
                    {form.competitorCandidates.map((candidate) => (
                      <button
                        key={candidate}
                        className="flex items-center gap-1.5 rounded-full border border-[#e2ddd6] bg-white px-3 py-2 text-xs hover:border-[#d97757] dark:border-[#4a4945] dark:bg-[#393836]"
                        disabled={form.competitors.includes(candidate)}
                        onClick={() => addCompetitor(candidate)}
                      >
                        <Plus size={13} /> {candidate}
                      </button>
                    ))}
                  </div>
                </>
              ) : (
                <Notice>
                  No reliable competitors were found automatically. Add the brands
                  you know below.
                </Notice>
              )}
              <FieldLabel htmlFor="competitor">Add a competitor</FieldLabel>
              <div className="flex gap-2">
                <div className={field}>
                  <Target size={17} className="text-[#8a847d]" />
                  <input
                    id="competitor"
                    className="min-w-0 flex-1 bg-transparent outline-none placeholder:text-[#aaa49d]"
                    value={competitorInput}
                    placeholder="competitor.com"
                    onChange={(event) => setCompetitorInput(event.target.value)}
                    onKeyDown={(event) => {
                      if (event.key === "Enter") addCompetitor(competitorInput)
                    }}
                  />
                </div>
                <button
                  aria-label="Add competitor"
                  className="grid size-12 shrink-0 place-items-center rounded-xl border border-[#ddd7cf] bg-white hover:bg-[#f7f4f0] disabled:opacity-40 dark:border-[#4a4945] dark:bg-[#393836]"
                  disabled={!competitorInput.trim() || form.competitors.length >= 5}
                  onClick={() => addCompetitor(competitorInput)}
                >
                  <Plus size={17} />
                </button>
              </div>
              <div className="mt-3 flex min-h-8 flex-wrap gap-2">
                {form.competitors.map((competitor) => (
                  <span
                    key={competitor}
                    className="flex items-center gap-2 rounded-lg border border-[#efcbbd] bg-[#fff5f0] px-2.5 py-1.5 text-xs text-[#a95034] dark:border-[#67483d] dark:bg-[#493a34] dark:text-[#efad92]"
                  >
                    {competitor}
                    <button
                      aria-label={`Remove ${competitor}`}
                      onClick={() =>
                        setForm((current) => ({
                          ...current,
                          competitors: current.competitors.filter(
                            (item) => item !== competitor,
                          ),
                        }))
                      }
                    >
                      <X size={12} />
                    </button>
                  </span>
                ))}
              </div>
              {notice && <Notice>{notice}</Notice>}
              <button
                className={`${primary} mt-5 w-full`}
                disabled={busy}
                onClick={() => void createPrompts()}
              >
                {busy ? (
                  <LoaderCircle size={16} className="animate-spin" />
                ) : (
                  <Sparkles size={16} />
                )}
                {busy ? "Building prompts…" : "Create my AEO prompts"}
              </button>
            </>
          )}

          {stage === "prompts" && (
            <>
              <p className={eyebrow}>AI VISIBILITY PROMPTS</p>
              <h1 className={heading}>Review the questions that matter</h1>
              <p className={`mb-6 mt-4 text-sm leading-6 ${muted}`}>
                AEOlyzer will track how answer engines respond to these questions.
              </p>
              <div className="max-h-[430px] space-y-2 overflow-y-auto pr-2">
                {prompts.map((prompt, index) => (
                  <label
                    key={`${index}-${prompt}`}
                    className="grid grid-cols-[30px_1fr] items-center gap-2 rounded-xl border border-[#e8e3dc] bg-[#fcfbf9] px-3 py-2 dark:border-[#4a4945] dark:bg-[#393836]"
                  >
                    <span className="font-mono text-[10px] text-[#d97757]">
                      {String(index + 1).padStart(2, "0")}
                    </span>
                    <textarea
                      aria-label={`Prompt ${index + 1}`}
                      className="resize-none bg-transparent text-xs leading-5 outline-none"
                      rows={2}
                      value={prompt}
                      onChange={(event) =>
                        setPrompts((current) =>
                          current.map((item, promptIndex) =>
                            promptIndex === index ? event.target.value : item,
                          ),
                        )
                      }
                    />
                  </label>
                ))}
              </div>
              <button className={`${primary} mt-5 w-full`} onClick={finish}>
                Finish setup <ArrowRight size={16} />
              </button>
            </>
          )}
        </div>
      </section>

      <LivePreview
        stage={stage}
        brandName={
          form.brandName ||
          form.suggestedBrandName ||
          fallbackBrandName(form.domain)
        }
        domain={displayHost(form.canonicalDomain || form.domain)}
        iconURL={form.iconURL}
        country={form.countryName}
        language={form.language}
      />
    </main>
  )
}

function Header({
  theme,
  onThemeChange,
  compact = false,
}: {
  theme: ThemeMode
  onThemeChange: (theme: ThemeMode) => void
  compact?: boolean
}) {
  return (
    <>
      <div className="absolute left-8 top-7 z-30 flex items-center gap-2">
        <AeolyzerLogo size={compact ? 25 : 29} />
        <span className="font-display text-xl font-semibold tracking-[-0.03em]">
          AEOlyzer
        </span>
      </div>
      <div className="absolute right-8 top-6 z-30 max-lg:right-5">
        <ThemeToggle theme={theme} onChange={onThemeChange} />
      </div>
    </>
  )
}

function AccountCard({
  active,
  icon,
  title,
  description,
  onClick,
}: {
  active: boolean
  icon: React.ReactNode
  title: string
  description: string
  onClick: () => void
}) {
  return (
    <button
      className={[
        panel,
        "grid grid-cols-[46px_1fr_22px] items-center gap-3 p-5 text-left transition hover:-translate-y-0.5",
        active
          ? "border-[#d97757] bg-[#fff8f4] ring-4 ring-[#d97757]/10 dark:bg-[#493a34]"
          : "",
      ].join(" ")}
      onClick={onClick}
    >
      <span className="grid size-11 place-items-center rounded-xl bg-[#f6eee9] text-[#d97757] dark:bg-[#4a4945]">
        {icon}
      </span>
      <span>
        <strong className="block text-sm">{title}</strong>
        <small className={`mt-1 block text-xs leading-5 ${muted}`}>
          {description}
        </small>
      </span>
      <span
        className={[
          "grid size-[22px] place-items-center rounded-full border",
          active
            ? "border-[#d97757] bg-[#d97757] text-white"
            : "border-[#cec7bf] dark:border-[#66635d]",
        ].join(" ")}
      >
        {active && <Check size={13} />}
      </span>
    </button>
  )
}

function FieldLabel({
  htmlFor,
  children,
}: {
  htmlFor?: string
  children: React.ReactNode
}) {
  return (
    <label
      htmlFor={htmlFor}
      className="mb-2 mt-4 block text-xs font-semibold text-[#4c4741] dark:text-[#d1cdc6]"
    >
      {children}
    </label>
  )
}

function SearchSelect<T>({
  label,
  placeholder,
  icon,
  value,
  options,
  optionKey,
  optionLabel,
  onSelect,
}: {
  label: string
  placeholder: string
  icon: React.ReactNode
  value: string
  options: T[]
  optionKey: (option: T) => string
  optionLabel: (option: T) => string
  onSelect: (option: T) => void
}) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState("")
  const filtered = options
    .filter((option) =>
      optionLabel(option).toLowerCase().includes(query.toLowerCase()),
    )
    .slice(0, 80)

  return (
    <div className="relative">
      <FieldLabel>{label}</FieldLabel>
      <button
        className={`${field} text-left`}
        onClick={() => setOpen((current) => !current)}
      >
        <span className="text-[#8a847d]">{icon}</span>
        <span
          className={[
            "min-w-0 flex-1 truncate text-xs",
            value ? "" : "text-[#aaa49d]",
          ].join(" ")}
        >
          {value || placeholder}
        </span>
        <ChevronDown size={15} />
      </button>
      {open && (
        <div className="absolute left-0 top-[calc(100%+6px)] z-50 w-full overflow-hidden rounded-xl border border-[#ded8d0] bg-white shadow-xl dark:border-[#55524c] dark:bg-[#393836]">
          <div className="flex items-center gap-2 border-b border-[#ebe6df] px-3 py-2 dark:border-[#4a4945]">
            <Search size={14} className="text-[#8b857e]" />
            <input
              autoFocus
              className="min-w-0 flex-1 bg-transparent text-xs outline-none"
              placeholder={placeholder}
              value={query}
              onChange={(event) => setQuery(event.target.value)}
            />
          </div>
          <div className="max-h-52 overflow-y-auto p-1.5">
            {filtered.map((option) => (
              <button
                key={optionKey(option)}
                className="flex w-full items-center justify-between rounded-lg px-2.5 py-2 text-left text-xs hover:bg-[#f5f1ec] dark:hover:bg-[#4a4945]"
                onClick={() => {
                  onSelect(option)
                  setOpen(false)
                  setQuery("")
                }}
              >
                {optionLabel(option)}
                {optionLabel(option) === value && (
                  <Check size={13} className="text-[#d97757]" />
                )}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

function Notice({ children }: { children: React.ReactNode }) {
  return (
    <div className="my-4 flex items-start gap-2 rounded-xl border border-[#efd2c6] bg-[#fff7f3] p-3 text-xs leading-5 text-[#8f4b35] dark:border-[#67483d] dark:bg-[#493a34] dark:text-[#efb29a]">
      <CircleGauge size={16} className="mt-0.5 shrink-0 text-[#d97757]" />
      {children}
    </div>
  )
}

function LivePreview({
  stage,
  brandName,
  domain,
  iconURL,
  country,
  language,
}: {
  stage: Stage
  brandName: string
  domain: string
  iconURL: string
  country: string
  language: string
}) {
  return (
    <aside className="relative flex min-h-screen items-center justify-center overflow-hidden bg-[#f5f3f0] p-12 max-lg:hidden dark:bg-[#24231f]">
      <div className="absolute size-[720px] rounded-full bg-[#d97757]/[0.06] blur-3xl dark:bg-[#e07b53]/10" />
      <div className="relative w-full max-w-[680px]">
        <div className="grid min-h-[450px] grid-cols-[58px_1fr] overflow-hidden rounded-2xl border border-[#ded8d0] bg-white shadow-[0_35px_90px_rgba(35,30,25,0.14)] dark:border-[#4a4945] dark:bg-[#2b2a27] dark:shadow-[0_40px_100px_rgba(0,0,0,0.4)]">
          <div className="flex flex-col items-center gap-5 border-r border-[#eee9e3] bg-[#faf9f7] pt-5 dark:border-[#3d3c38] dark:bg-[#252422]">
            <AeolyzerLogo size={21} />
            {[0, 1, 2, 3].map((item) => (
              <span
                key={item}
                className="h-2 w-5 rounded bg-[#dad4cc] dark:bg-[#4a4945]"
              />
            ))}
          </div>
          <div className="p-6">
            <div className="mb-6 flex justify-between">
              <span className="h-2 w-20 rounded bg-[#e2ddd6] dark:bg-[#4a4945]" />
              <span className="h-2 w-8 rounded bg-[#e2ddd6] dark:bg-[#4a4945]" />
            </div>
            <div className="flex items-center gap-3 rounded-xl border border-[#e5dfd8] bg-[#fbfaf8] p-3 dark:border-[#4a4945] dark:bg-[#393836]">
              <span className="grid size-9 place-items-center overflow-hidden rounded-lg bg-[#f6e7df] font-semibold text-[#c96849] dark:bg-[#4a4945]">
                <SiteImage
                  src={iconURL}
                  className="size-full object-contain"
                  fallback={(brandName || "A").charAt(0).toUpperCase()}
                />
              </span>
              <span className="min-w-0 flex-1">
                <strong className="block truncate text-xs">
                  {brandName || "Your brand"}
                </strong>
                <small className={`mt-1 block truncate text-[9px] ${muted}`}>
                  {domain || "yourbrand.com"}
                </small>
              </span>
              <CheckCircle2 size={17} className="text-[#5f9d70]" />
            </div>
            <p className={`${eyebrow} mb-1 mt-8 text-[8px]`}>AI VISIBILITY</p>
            <h2 className="font-display text-2xl font-medium tracking-[-0.03em]">
              Good morning, {brandName || "your brand"}
            </h2>
            <div className="mt-5 grid grid-cols-3 gap-3">
              {["Visibility score", "Brand mentions", "Source coverage"].map(
                (label) => (
                  <div
                    key={label}
                    className="rounded-xl border border-[#e7e2db] bg-[#fbfaf8] p-3 dark:border-[#4a4945] dark:bg-[#393836]"
                  >
                    <small className={`block text-[8px] ${muted}`}>{label}</small>
                    <b className="my-2 block text-lg">—</b>
                    <span className="block h-1 w-2/3 rounded bg-[#ded8d0] dark:bg-[#4a4945]" />
                  </div>
                ),
              )}
            </div>
            <div className="relative mt-3 h-28 overflow-hidden rounded-xl border border-[#e7e2db] bg-[linear-gradient(#eee9e3_1px,transparent_1px),linear-gradient(90deg,#eee9e3_1px,transparent_1px)] bg-[size:30px_30px] dark:border-[#4a4945] dark:bg-[linear-gradient(#3f3e3a_1px,transparent_1px),linear-gradient(90deg,#3f3e3a_1px,transparent_1px)]">
              <span className="absolute left-[12%] top-[62%] h-0.5 w-3/4 -skew-y-6 bg-[#d97757]" />
            </div>
            {(country || language) && (
              <div className={`mt-3 flex justify-end gap-3 text-[9px] ${muted}`}>
                <span>{country}</span>
                <span className="border-l border-[#ddd7cf] pl-3 dark:border-[#4a4945]">
                  {language}
                </span>
              </div>
            )}
          </div>
        </div>
        <p className={`mt-7 text-center text-xs ${muted}`}>
          {previewCaption(stage)}
        </p>
      </div>
    </aside>
  )
}

function SiteImage({
  src,
  className,
  fallback,
}: {
  src: string
  className?: string
  fallback: React.ReactNode
}) {
  const [failed, setFailed] = useState(false)
  if (!src || failed) return <>{fallback}</>
  return (
    // eslint-disable-next-line @next/next/no-img-element
    <img
      className={className}
      src={src}
      alt=""
      onError={() => setFailed(true)}
    />
  )
}

function previewCaption(stage: Stage) {
  const captions: Partial<Record<Stage, string>> = {
    domain: "Your dashboard updates as you enter your website.",
    brand: "Your brand identity, reflected instantly.",
    market: "Local context makes your insights more relevant.",
    competitors: "Competitor context sharpens every recommendation.",
    prompts: "Your visibility questions are ready to track.",
  }
  return captions[stage] ?? ""
}

function previousStage(stage: Stage): Stage {
  const stages: Partial<Record<Stage, Stage>> = {
    domain: "account",
    brand: "domain",
    market: "brand",
    competitors: "market",
    prompts: "competitors",
  }
  return stages[stage] ?? "account"
}

function stageLabel(stage: Stage) {
  const labels: Partial<Record<Stage, string>> = {
    domain: "Step 2 of 6",
    brand: "Step 3 of 6",
    market: "Step 4 of 6",
    competitors: "Step 5 of 6",
    prompts: "Step 6 of 6",
  }
  return labels[stage] ?? ""
}

function plausibleDomain(value: string) {
  try {
    const url = new URL(value.includes("://") ? value : `https://${value}`)
    return url.hostname.includes(".") && !url.hostname.endsWith(".")
  } catch {
    return false
  }
}

function normalizeURL(value: string) {
  try {
    return new URL(value.includes("://") ? value : `https://${value}`).toString()
  } catch {
    return `https://${value}`
  }
}

function displayHost(value: string) {
  if (!value) return ""
  try {
    return new URL(value.includes("://") ? value : `https://${value}`).hostname.replace(
      /^www\./,
      "",
    )
  } catch {
    return value
  }
}
