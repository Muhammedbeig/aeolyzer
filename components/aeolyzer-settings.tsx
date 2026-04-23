"use client"

import { useState } from "react"
import { X, ChevronDown } from "lucide-react"
import { cn } from "@/lib/utils"

type Theme = "light" | "auto" | "dark"

// -> defines the shape of the data (props) the AeolyzerSettings block expects
// - `isOpen` -> true/false switch to decide if the settings menu should be shown or hidden
// - `onClose` -> function to call when the user wants to close the menu
// - `theme` -> remembers the current visual theme ("light", "auto", or "dark")
// - `onThemeChange` -> function to call to actually change the theme when the user selects a new one
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

// -> This draws the main popup window containing all user settings
export function AeolyzerSettings({ isOpen, onClose, theme, onThemeChange }: SettingsProps) {
  // - `const [activeTab, setActiveTab]` -> memory that tracks which menu tab on the left is currently clicked (defaults to "General")
  const [activeTab, setActiveTab] = useState("General")
  
  // - `const [backgroundAnimation, setBackgroundAnimation]` -> memory tracking if background movement is enabled
  const [backgroundAnimation, setBackgroundAnimation] = useState<"enabled" | "auto" | "disabled">("auto")

  // -> If `isOpen` is false, this completely stops the component from drawing anything at all
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
        className="relative w-full max-w-4xl max-h-[85vh] rounded-2xl overflow-hidden shadow-2xl flex"
        style={{ backgroundColor: "#2b2a27" }}
      >
        {/* Sidebar */}
        <div 
          className="w-48 flex-shrink-0 p-4 border-r"
          style={{ borderColor: "#4a4945" }}
        >
          <h2 className="text-xl font-semibold mb-6" style={{ color: "#ececec" }}>Settings</h2>
          <nav className="space-y-1">
            {settingsTabs.map((tab) => (
              <button
                key={tab}
                onClick={() => setActiveTab(tab)}
                className={cn(
                  "w-full text-left px-3 py-2 rounded-lg text-sm transition-colors",
                  activeTab === tab 
                    ? "bg-[#252422]" 
                    : "hover:bg-[#252422]"
                )}
                style={{ color: activeTab === tab ? "#ececec" : "#a3a29e" }}
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
            className="absolute top-4 right-4 p-2 rounded-lg transition-colors hover:bg-[#252422]"
            style={{ color: "#a3a29e" }}
          >
            <X size={20} />
          </button>

          {activeTab === "General" && (
            <div className="space-y-8">
              {/* Profile Section */}
              <section>
                <h3 className="text-base font-medium mb-4" style={{ color: "#ececec" }}>Profile</h3>
                <div className="grid grid-cols-2 gap-6">
                  <div>
                    <label className="block text-sm mb-2" style={{ color: "#a3a29e" }}>Full name</label>
                    <div className="flex items-center gap-3">
                      <div 
                        className="w-10 h-10 rounded-full flex items-center justify-center text-sm font-medium flex-shrink-0"
                        style={{ backgroundColor: "#4a4945", color: "#ececec" }}
                      >
                        M
                      </div>
                      <input
                        type="text"
                        defaultValue="Muhammad"
                        className="flex-1 px-4 py-2.5 rounded-lg text-sm outline-none transition-colors focus:ring-2 focus:ring-[#e07b53]"
                        style={{ backgroundColor: "#393836", color: "#ececec", border: "1px solid #4a4945" }}
                      />
                    </div>
                  </div>
                  <div>
                    <label className="block text-sm mb-2" style={{ color: "#a3a29e" }}>What should Aeolyzer call you?</label>
                    <input
                      type="text"
                      defaultValue="Muhammad"
                      className="w-full px-4 py-2.5 rounded-lg text-sm outline-none transition-colors focus:ring-2 focus:ring-[#e07b53]"
                      style={{ backgroundColor: "#393836", color: "#ececec", border: "1px solid #4a4945" }}
                    />
                  </div>
                </div>

                <div className="mt-4">
                  <label className="block text-sm mb-2" style={{ color: "#a3a29e" }}>What best describes your work?</label>
                  <button
                    className="w-full flex items-center justify-between px-4 py-2.5 rounded-lg text-sm"
                    style={{ backgroundColor: "#393836", color: "#a3a29e", border: "1px solid #4a4945" }}
                  >
                    <span>Select your work function</span>
                    <ChevronDown size={16} />
                  </button>
                </div>

                <div className="mt-4">
                  <label className="block text-sm mb-2" style={{ color: "#a3a29e" }}>
                    What <span className="underline">personal preferences</span> should Aeolyzer consider in responses?
                  </label>
                  <p className="text-xs mb-2" style={{ color: "#a3a29e" }}>
                    Your preferences will apply to all conversations, within <span className="underline">Anthropic&apos;s guidelines</span>.
                  </p>
                  <textarea
                    placeholder="e.g. I primarily code in Python (not a coding beginner)"
                    className="w-full px-4 py-3 rounded-lg text-sm outline-none resize-none h-24 transition-colors focus:ring-2 focus:ring-[#e07b53]"
                    style={{ backgroundColor: "#393836", color: "#ececec", border: "1px solid #4a4945" }}
                  />
                </div>
              </section>

              {/* Appearance Section */}
              <section className="border-t pt-8" style={{ borderColor: "#4a4945" }}>
                <h3 className="text-base font-medium mb-4" style={{ color: "#ececec" }}>Appearance</h3>
                
                {/* Color mode */}
                <div className="mb-6">
                  <label className="block text-sm mb-3" style={{ color: "#a3a29e" }}>Color mode</label>
                  <div className="flex gap-4">
                    {(["light", "auto", "dark"] as const).map((mode) => (
                      <button
                        key={mode}
                        onClick={() => onThemeChange(mode)}
                        className={cn(
                          "flex flex-col items-center gap-2 p-3 rounded-xl transition-all",
                          theme === mode 
                            ? "ring-2 ring-[#e07b53]" 
                            : "hover:bg-[#252422]"
                        )}
                        style={{ border: "1px solid #4a4945" }}
                      >
                        {/* Theme preview */}
                        <div 
                          className="w-28 h-20 rounded-lg overflow-hidden"
                          style={{ 
                            backgroundColor: mode === "light" ? "#ffffff" : mode === "dark" ? "#2b2a27" : "#393836",
                            border: "1px solid #4a4945"
                          }}
                        >
                          <div className="p-2 space-y-1.5">
                            <div 
                              className="h-1.5 rounded-full w-12"
                              style={{ backgroundColor: mode === "light" ? "#d1d5db" : "#4a4945" }}
                            />
                            <div 
                              className="h-1.5 rounded-full w-16"
                              style={{ backgroundColor: mode === "light" ? "#d1d5db" : "#4a4945" }}
                            />
                            <div 
                              className="h-1.5 rounded-full w-10"
                              style={{ backgroundColor: mode === "light" ? "#d1d5db" : "#4a4945" }}
                            />
                          </div>
                          <div 
                            className="mx-2 mt-2 h-4 rounded flex items-center justify-end pr-1"
                            style={{ backgroundColor: mode === "light" ? "#f3f4f6" : "#393836" }}
                          >
                            <div 
                              className="w-3 h-2 rounded"
                              style={{ backgroundColor: "#e07b53" }}
                            />
                          </div>
                        </div>
                        <span 
                          className="text-sm capitalize"
                          style={{ color: theme === mode ? "#ececec" : "#a3a29e" }}
                        >
                          {mode}
                        </span>
                      </button>
                    ))}
                  </div>
                </div>

                {/* Background animation */}
                <div>
                  <label className="block text-sm mb-3" style={{ color: "#a3a29e" }}>Background animation</label>
                  <div className="flex gap-4">
                    {(["enabled", "auto", "disabled"] as const).map((mode) => (
                      <button
                        key={mode}
                        onClick={() => setBackgroundAnimation(mode)}
                        className={cn(
                          "flex flex-col items-center gap-2 p-3 rounded-xl transition-all",
                          backgroundAnimation === mode 
                            ? "ring-2 ring-[#e07b53]" 
                            : "hover:bg-[#252422]"
                        )}
                        style={{ border: "1px solid #4a4945" }}
                      >
                        {/* Animation preview */}
                        <div 
                          className="w-28 h-20 rounded-lg flex items-center justify-center"
                          style={{ backgroundColor: "#393836", border: "1px solid #4a4945" }}
                        >
                          <div className="flex gap-1.5">
                            {mode === "enabled" && (
                              <>
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e] animate-pulse" />
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e] animate-pulse delay-75" />
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e] animate-pulse delay-150" />
                              </>
                            )}
                            {mode === "auto" && (
                              <>
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e]" />
                                <div className="w-px h-4 bg-[#4a4945]" />
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e]" />
                              </>
                            )}
                            {mode === "disabled" && (
                              <>
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e]" />
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e]" />
                                <div className="w-1.5 h-1.5 rounded-full bg-[#a3a29e]" />
                              </>
                            )}
                          </div>
                        </div>
                        <span 
                          className="text-sm capitalize"
                          style={{ color: backgroundAnimation === mode ? "#ececec" : "#a3a29e" }}
                        >
                          {mode}
                        </span>
                      </button>
                    ))}
                  </div>
                </div>
              </section>

              {/* Notifications Section */}
              <section className="border-t pt-8" style={{ borderColor: "#4a4945" }}>
                <h3 className="text-base font-medium mb-4" style={{ color: "#ececec" }}>Notifications</h3>
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium" style={{ color: "#e07b53" }}>Response completions</p>
                    <p className="text-sm mt-1" style={{ color: "#a3a29e" }}>
                      Get notified when Aeolyzer has finished a response. Most useful for long-running tasks like tool calls and Research.
                    </p>
                  </div>
                  <button
                    className="w-11 h-6 rounded-full transition-colors relative bg-[#e07b53]"
                  >
                    <span className="absolute top-1 right-1 w-4 h-4 rounded-full bg-white transition-transform" />
                  </button>
                </div>
              </section>

              {/* Chat font Section */}
              <section className="border-t pt-8" style={{ borderColor: "#4a4945" }}>
                <h3 className="text-base font-medium mb-4" style={{ color: "#ececec" }}>Chat font</h3>
                <div className="flex gap-3">
                  <button 
                    className="px-4 py-2 rounded-lg text-sm ring-2 ring-[#e07b53]"
                    style={{ backgroundColor: "#252422", color: "#ececec" }}
                  >
                    Default
                  </button>
                  <button 
                    className="px-4 py-2 rounded-lg text-sm hover:bg-[#252422]"
                    style={{ color: "#a3a29e", border: "1px solid #4a4945" }}
                  >
                    Monospace
                  </button>
                </div>
              </section>
            </div>
          )}

          {activeTab !== "General" && (
            <div className="flex items-center justify-center h-64">
              <p style={{ color: "#a3a29e" }}>{activeTab} settings coming soon...</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
