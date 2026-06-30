import { useState } from "react"
import { cn } from "@/lib/utils"
import { Theme } from "./types"

interface AppearanceSectionProps {
  theme: Theme;
  onThemeChange: (theme: Theme) => void;
}

export function AppearanceSection({ theme, onThemeChange }: AppearanceSectionProps) {
  const [backgroundAnimation, setBackgroundAnimation] = useState<"enabled" | "auto" | "disabled">("auto")

  return (
    <section className="border-t pt-8 border-border">
      <h3 className="text-base font-medium mb-4 text-foreground">Appearance</h3>
      
      {/* Color mode */}
      <div className="mb-6">
        <label className="block text-sm mb-3 text-muted-foreground">Color mode</label>
        <div className="flex gap-4">
          {(["light", "system", "dark"] as const).map((mode) => (
            <button
              key={mode}
              onClick={() => onThemeChange(mode)}
              className={cn(
                "flex flex-col items-center gap-2 p-3 rounded-xl transition-all",
                theme === mode 
                  ? "ring-2 ring-[var(--accent)]" 
                  : "hover:bg-muted"
              )}
              style={{ border: "1px solid var(--border)" }}
            >
              {/* Theme preview */}
              <div 
                className="w-28 h-20 rounded-lg overflow-hidden border border-border"
                style={{ backgroundColor: mode === "light" ? "#ffffff" : mode === "dark" ? "var(--background)" : "var(--card)" }}
              >
                <div className="p-2 space-y-1.5">
                  <div 
                    className="h-1.5 rounded-full w-12"
                    style={{ backgroundColor: mode === "light" ? "#d1d5db" : "var(--border)" }}
                  />
                  <div 
                    className="h-1.5 rounded-full w-16"
                    style={{ backgroundColor: mode === "light" ? "#d1d5db" : "var(--border)" }}
                  />
                  <div 
                    className="h-1.5 rounded-full w-10"
                    style={{ backgroundColor: mode === "light" ? "#d1d5db" : "var(--border)" }}
                  />
                </div>
                <div 
                  className="mx-2 mt-2 h-4 rounded flex items-center justify-end pr-1"
                  style={{ backgroundColor: mode === "light" ? "#f3f4f6" : "var(--card)" }}
                >
                  <div 
                    className="w-3 h-2 rounded bg-accent"
                  />
                </div>
              </div>
              <span 
                className="text-sm capitalize"
                style={{ color: theme === mode ? "var(--foreground)" : "var(--muted-foreground)" }}
              >
                {mode}
              </span>
            </button>
          ))}
        </div>
      </div>

      {/* Background animation */}
      <div>
        <label className="block text-sm mb-3 text-muted-foreground">Background animation</label>
        <div className="flex gap-4">
          {(["enabled", "auto", "disabled"] as const).map((mode) => (
            <button
              key={mode}
              onClick={() => setBackgroundAnimation(mode)}
              className={cn(
                "flex flex-col items-center gap-2 p-3 rounded-xl transition-all",
                backgroundAnimation === mode 
                  ? "ring-2 ring-[var(--accent)]" 
                  : "hover:bg-muted"
              )}
              style={{ border: "1px solid var(--border)" }}
            >
              {/* Animation preview */}
              <div 
                className="w-28 h-20 rounded-lg flex items-center justify-center"
                style={{ backgroundColor: "var(--card)", border: "1px solid var(--border)" }}
              >
                <div className="flex gap-1.5">
                  {mode === "enabled" && (
                    <>
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)] animate-pulse" />
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)] animate-pulse delay-75" />
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)] animate-pulse delay-150" />
                    </>
                  )}
                  {mode === "auto" && (
                    <>
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)]" />
                      <div className="w-px h-4 bg-[var(--border)]" />
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)]" />
                    </>
                  )}
                  {mode === "disabled" && (
                    <>
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)]" />
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)]" />
                      <div className="w-1.5 h-1.5 rounded-full bg-[var(--muted-foreground)]" />
                    </>
                  )}
                </div>
              </div>
              <span 
                className="text-sm capitalize"
                style={{ color: backgroundAnimation === mode ? "var(--foreground)" : "var(--muted-foreground)" }}
              >
                {mode}
              </span>
            </button>
          ))}
        </div>
      </div>
    </section>
  )
}
