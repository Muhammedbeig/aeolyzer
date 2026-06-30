"use client"

import { useState } from "react"
import { X } from "lucide-react"
import { SettingsProps } from "./types"
import { SettingsNav } from "./settings-nav"
import { GeneralTab } from "./general-tab"

export function AeolyzerSettings({ isOpen, onClose, theme, onThemeChange }: SettingsProps) {
  const [activeTab, setActiveTab] = useState("General")

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
        className="relative w-full max-w-4xl h-full md:h-auto max-h-[100dvh] md:max-h-[85vh] md:rounded-2xl overflow-hidden shadow-2xl flex flex-col md:flex-row bg-background"
      >
        <SettingsNav activeTab={activeTab} onTabChange={setActiveTab} onClose={onClose} />

        {/* Content */}
        <div className="flex-1 overflow-y-auto p-4 md:p-6 relative custom-scrollbar">
          {/* Close button (Desktop) */}
          <button
            onClick={onClose}
            className="hidden md:block absolute top-4 right-4 p-2 rounded-lg transition-colors hover:bg-muted text-muted-foreground"
          >
            <X size={20} />
          </button>

          {activeTab === "General" ? (
            <GeneralTab theme={theme} onThemeChange={onThemeChange} />
          ) : (
            <div className="flex items-center justify-center h-64">
              <p className="text-muted-foreground">{activeTab} settings coming soon...</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
