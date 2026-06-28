"use client"

import { MonitorCog, Sun } from "lucide-react"
import { useEffect, useState } from "react"
import { AeolyzerLogo } from "@/components/aeolyzer-logo"
import { OnboardingFlow } from "@/components/onboarding-flow"
import { AeolyzerWorkspace } from "@/components/aeolyzer-workspace"
import type { ProjectProfile } from "@/lib/aeolyzer"

export type ThemeMode = "light" | "system"

const projectStorageKey = "aeolyzer.guest.project"
const themeStorageKey = "aeolyzer.theme"

export function AeolyzerApp() {
  const [hydrated, setHydrated] = useState(false)
  const [project, setProject] = useState<ProjectProfile | null>(null)
  const [theme, setTheme] = useState<ThemeMode>("light")

  useEffect(() => {
    const storedTheme = localStorage.getItem(themeStorageKey)
    if (storedTheme === "system") setTheme("system")

    const storedProject = sessionStorage.getItem(projectStorageKey)
    if (storedProject) {
      try {
        setProject(JSON.parse(storedProject) as ProjectProfile)
      } catch {
        sessionStorage.removeItem(projectStorageKey)
      }
    }
    setHydrated(true)
  }, [])

  useEffect(() => {
    document.documentElement.classList.toggle("dark", theme === "system")
    document.documentElement.dataset.theme = theme
    localStorage.setItem(themeStorageKey, theme)
  }, [theme])

  if (!hydrated) {
    return (
      <main className="grid min-h-screen place-items-center bg-[#fbfaf8] text-[#1d1b18] dark:bg-[#2b2a27] dark:text-[#ececec]">
        <AeolyzerLogo size={38} animate />
      </main>
    )
  }

  if (!project) {
    return (
      <OnboardingFlow
        theme={theme}
        onThemeChange={setTheme}
        onComplete={(nextProject) => {
          sessionStorage.setItem(projectStorageKey, JSON.stringify(nextProject))
          setProject(nextProject)
        }}
      />
    )
  }

  return (
    <AeolyzerWorkspace
      project={project}
      theme={theme}
      onThemeChange={setTheme}
      onReset={() => {
        sessionStorage.removeItem(projectStorageKey)
        sessionStorage.removeItem("aeolyzer.guest.audit-history")
        sessionStorage.removeItem("aeolyzer.guest.content-history")
        setProject(null)
      }}
    />
  )
}

export function ThemeToggle({
  theme,
  onChange,
  compact = false,
}: {
  theme: ThemeMode
  onChange: (theme: ThemeMode) => void
  compact?: boolean
}) {
  return (
    <div
      aria-label="Color mode"
      className={[
        "inline-flex items-center rounded-xl border border-[#e6e1da] bg-white p-1 shadow-sm",
        "dark:border-[#4a4945] dark:bg-[#393836]",
        compact ? "w-full" : "",
      ].join(" ")}
    >
      <button
        className={[
          "flex items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-xs font-medium transition",
          compact ? "flex-1" : "",
          theme === "light"
            ? "bg-[#f2ede7] text-[#1d1b18] dark:bg-[#4a4945] dark:text-white"
            : "text-[#7a746d] hover:text-[#1d1b18] dark:text-[#aaa7a1] dark:hover:text-white",
        ].join(" ")}
        onClick={() => onChange("light")}
      >
        <Sun size={14} />
        Light
      </button>
      <button
        className={[
          "flex items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-xs font-medium transition",
          compact ? "flex-1" : "",
          theme === "system"
            ? "bg-[#2b2a27] text-white"
            : "text-[#7a746d] hover:text-[#1d1b18] dark:text-[#aaa7a1] dark:hover:text-white",
        ].join(" ")}
        onClick={() => onChange("system")}
      >
        <MonitorCog size={14} />
        System
      </button>
    </div>
  )
}
