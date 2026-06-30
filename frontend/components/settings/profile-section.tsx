import { ChevronDown } from "lucide-react"

export function ProfileSection() {
  return (
    <section>
      <h3 className="text-base font-semibold mb-6 text-foreground">Profile</h3>
      <div className="space-y-6">
        {/* Avatar */}
        <div className="flex items-center justify-between gap-4">
          <label className="text-sm font-medium text-foreground flex-1">Avatar</label>
          <div className="flex-1 flex justify-end">
            <div className="w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium flex-shrink-0 bg-muted text-foreground">
              M
            </div>
          </div>
        </div>
        
        {/* Full name */}
        <div className="flex items-center justify-between gap-4">
          <label className="text-sm font-medium text-foreground flex-1 break-words pr-2">Full name</label>
          <div className="flex-1 flex justify-end">
            <input
              type="text"
              defaultValue="Muhammad"
              className="w-full max-w-[220px] md:max-w-[280px] px-3 py-2 rounded-lg text-sm outline-none transition-colors focus:ring-2 focus:ring-[var(--accent)] bg-black/5 dark:bg-white/5 text-foreground border border-transparent focus:border-border"
            />
          </div>
        </div>

        {/* What should Aeolyzer call you? */}
        <div className="flex items-center justify-between gap-4">
          <label className="text-sm font-medium text-foreground flex-1 break-words pr-2">What should Aeolyzer call you?</label>
          <div className="flex-1 flex justify-end">
            <input
              type="text"
              defaultValue="Muhammad"
              className="w-full max-w-[220px] md:max-w-[280px] px-3 py-2 rounded-lg text-sm outline-none transition-colors focus:ring-2 focus:ring-[var(--accent)] bg-black/5 dark:bg-white/5 text-foreground border border-transparent focus:border-border"
            />
          </div>
        </div>

        {/* What best describes your work? */}
        <div className="flex items-center justify-between gap-4">
          <label className="text-sm font-medium text-foreground flex-1 break-words pr-2">What best describes your work?</label>
          <div className="flex-1 flex justify-end">
            <button className="flex items-center justify-between gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors w-full max-w-[220px] md:max-w-[280px] px-3 py-2 bg-transparent">
              <span>Select</span>
              <ChevronDown size={14} />
            </button>
          </div>
        </div>
      </div>

      <div className="mt-8">
        <label className="block text-sm font-medium mb-2 text-foreground">Instructions for Aeolyzer</label>
        <p className="text-[13px] mb-3 text-muted-foreground leading-relaxed">
          Aeolyzer will keep these in mind across chats and Cowork within <span className="underline cursor-pointer">Anthropic&apos;s guidelines</span>. <span className="underline cursor-pointer">Learn more</span>
        </p>
        <textarea
          placeholder="e.g. when learning new concepts, I find analogies particularly helpful"
          className="w-full px-4 py-3 rounded-xl text-sm outline-none resize-none h-24 transition-colors focus:ring-2 focus:ring-[var(--accent)] bg-muted/30 text-foreground border border-border"
        />
      </div>
    </section>
  )
}
