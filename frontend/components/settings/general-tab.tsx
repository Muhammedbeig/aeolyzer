import { Theme } from "./types"
import { ProfileSection } from "./profile-section"
import { AppearanceSection } from "./appearance-section"

interface GeneralTabProps {
  theme: Theme;
  onThemeChange: (theme: Theme) => void;
}

export function GeneralTab({ theme, onThemeChange }: GeneralTabProps) {
  return (
    <div className="space-y-8 max-w-2xl">
      <ProfileSection />
      
      <AppearanceSection theme={theme} onThemeChange={onThemeChange} />

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
  )
}
