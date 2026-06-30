"use client"

import { useState } from "react"
import { X, ChevronDown } from "lucide-react"
import { cn } from "@/lib/utils"

type Theme = "light" | "system" | "dark"

interface SettingsProps {
  isOpen: boolean
  onClose: () => void
  theme: Theme
  onThemeChange: (theme: Theme) => void
}

const settingsTabs = [
  "General",
  "Account",
  "Privacy",
  "Billing",
  "Capabilities",
  "Connectors",
  "Aeolyzer Code",
]

export function AeolyzerSettings({ isOpen, onClose, theme, onThemeChange }: SettingsProps) {
  const [activeTab, setActiveTab] = useState("General")
  const [backgroundAnimation, setBackgroundAnimation] = useState<"enabled" | "auto" | "disabled">("auto")

  if (!isOpen) return null

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-black/60"
        onClick={onClose}
      />
      
      {/* Modal */}
      <div 
        className="relative w-full max-w-4xl max-h-[85vh] rounded-2xl overflow-hidden shadow-2xl flex bg-background"
      >
        {/* Sidebar */}
        <div 
          className="w-48 flex-shrink-0 p-4 border-r border-border"
        >
          <h2 className="text-xl font-semibold mb-6 text-foreground">Settings</h2>
          <nav className="space-y-1">
            {settingsTabs.map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={cn(
                  "w-full text-left px-3 py-2 rounded-lg text-sm transition-colors",
                  activeTab === tab 
                    ? "bg-muted" 
                    : "hover:bg-muted"
                )}
                style={{ color: activeTab === tab ? "var(--foreground)" : "var(--muted-foreground)" }}
              >
                {tab}
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-6">
          {/* Close button */}
          <button
            onClick={onClose}
            className="absolute top-4 right-4 p-2 rounded-lg transition-colors hover:bg-muted text-muted-foreground"
          >
            <X size={20} />
          </button>

          {activeTab === "General" && (
            <div className="space-y-8">
              {/* Profile Section */}
              <section>
                <h3 className="text-base font-medium mb-4 text-foreground">Profile</h3>
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <label className="block text-sm mb-2 text-muted-foreground">Full name</label>
                    <div className="flex items-center gap-3">
                      <div 
                        className="w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium flex-shrink-0 bg-border text-foreground"
                      >
                        M
                      </div>
                      <input
                        type="text"
                        defaultValue="Muhammad"
                        className="flex-1 px-4 py-2.5 rounded-lg text-sm outline-none transition-colors focus:ring-2 focus:ring-[var(--accent)] bg-card text-foreground border border-border"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm mb-2 text-muted-foreground">What should Aeolyzer call you?</label>
                    <input
                      type="text"
                      defaultValue="Muhammad"
                      className="w-full px-4 py-2.5 rounded-lg text-sm outline-none transition-colors focus:ring-2 focus:ring-[var(--accent)] bg-card text-foreground border border-border"
                    />
                  </div>
                </div>

                <div className="mt-4">
                  <label className="block text-sm mb-2 text-muted-foreground">What best describes your work?</label>
                  <button
                    className="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-sm bg-card text-muted-foreground border border-border"
                  >
                    <span>Select your work function</span>
                    <ChevronDown size={16} />
                  </button>
                </div>

                <div className="mt-4">
                  <label className="block text-sm mb-2 text-muted-foreground">
                    What <span className="underline">personal preferences</span> should Aeolyzer consider in responses?
                  </label>
                  <p className="text-xs mb-2 text-muted-foreground">
                    Your preferences will apply to all conversations, within <span className="underline">Anthropic&apos;s guidelines</span>.
                  </p>
                  <textarea
                    placeholder="e.g. I primarily code in Python (not a coding beginner)"
                    className="w-full px-4 py-3 rounded-lg text-sm outline-none resize-none h-24 transition-colors focus:ring-2 focus:ring-[var(--accent)] bg-card text-foreground border border-border"
                  />
                </div>
              </section>

              {/* Appearance Section */}
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
                          className="w-28 h-20 rounded-lg overflow-hidden"
                          style={{ backgroundColor: mode === "light" ? "#ffffff" : mode === "dark" ? "var(--background)" : "var(--card)" }} className="border border-border"
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

              {/* Notifications Section */}
              <section className="border-t pt-8 border-border">
                <h3 className="text-base font-medium mb-4 text-foreground">Notifications</h3>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-accent">Response completions</p>
                    <p className="text-sm mt-1 text-muted-foreground">
                      Get notified when Aeolyzer has finished a response. Most useful for long-running tasks like tool calls and Research.
                    </p>
                  </div>
                  <button
                    className="w-11 h-6 rounded-full transition-colors relative bg-[var(--accent)]"
                  >
                    <span className="absolute top-1 right-1 w-4 h-4 rounded-full bg-white transition-transform" />
                  </button>
                </div>
              </section>

              {/* Chat font Section */}
              <section className="border-t pt-8 border-border">
                <h3 className="text-base font-medium mb-4 text-foreground">Chat font</h3>
                <div className="flex gap-3">
                  <button 
                    className="px-4 py-2 rounded-lg text-sm ring-2 ring-[var(--accent)] bg-sidebar-hover text-foreground"
                  >
                    Default
                  </button>
                  <button 
                    className="px-4 py-2 rounded-lg text-sm hover:bg-muted text-muted-foreground border border-border"
                  >
                    Monospace
                  </button>
                </div>
              </section>
            </div>
          )}

          {activeTab !== "General" && (
            <div className="flex items-center justify-center h-64">
              <p className="text-muted-foreground">{activeTab} settings coming soon...</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
