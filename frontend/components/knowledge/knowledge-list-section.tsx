"use client"

import { useState, type FormEvent } from "react"
import { X, type LucideIcon } from "lucide-react"

interface KnowledgeListSectionProps {
  title: string
  description: string
  emptyTitle: string
  emptyDescription: string
  inputLabel: string
  placeholder: string
  actionLabel: string
  icon: LucideIcon
  items: string[]
  saving: boolean
  onSave: (items: string[]) => Promise<boolean>
}

export function KnowledgeListSection({
  title,
  description,
  emptyTitle,
  emptyDescription,
  inputLabel,
  placeholder,
  actionLabel,
  icon: Icon,
  items,
  saving,
  onSave,
}: KnowledgeListSectionProps) {
  const [draft, setDraft] = useState("")

  const handleAdd = async (event: FormEvent) => {
    event.preventDefault()
    const value = draft.trim()
    if (!value || saving) {
      return
    }
    if (await onSave([...items, value])) {
      setDraft("")
    }
  }

  const handleRemove = async (index: number) => {
    if (saving) {
      return
    }
    await onSave(items.filter((_, itemIndex) => itemIndex !== index))
  }

  return (
    <div className="max-w-3xl space-y-6" data-testid="knowledge-list-section">
      <div>
        <h1 className="mb-1 text-xl font-semibold">{title}</h1>
        <p className="text-sm text-muted-foreground">{description}</p>
      </div>

      <form onSubmit={handleAdd} className="flex gap-2">
        <label htmlFor="knowledge-list-value" className="sr-only">
          {inputLabel}
        </label>
        <input
          id="knowledge-list-value"
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          placeholder={placeholder}
          disabled={saving}
          className="min-w-0 flex-1 rounded-md border border-black/10 bg-transparent px-3 py-2 transition-colors focus:border-accent focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10"
        />
        <button
          type="submit"
          disabled={!draft.trim() || saving}
          className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-white transition-opacity hover:opacity-90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent/40 active:scale-95 disabled:pointer-events-none disabled:opacity-50"
        >
          {saving ? "Saving..." : actionLabel}
        </button>
      </form>

      {items.length === 0 ? (
        <div className="flex flex-col items-center justify-center space-y-3 rounded-xl border border-black/10 bg-black/[0.02] p-8 text-center dark:border-white/10 dark:bg-white/[0.02]">
          <Icon className="h-8 w-8 text-muted-foreground opacity-50" aria-hidden="true" />
          <div>
            <h3 className="font-medium">{emptyTitle}</h3>
            <p className="mt-1 max-w-sm text-sm text-muted-foreground">
              {emptyDescription}
            </p>
          </div>
        </div>
      ) : (
        <ul className="space-y-2" aria-label={title}>
          {items.map((item, index) => (
            <li
              key={`${item}-${index}`}
              className="flex items-start justify-between gap-3 rounded-lg border border-black/10 bg-black/[0.02] px-4 py-3 dark:border-white/10 dark:bg-white/[0.02]"
            >
              <span className="min-w-0 break-words text-sm">{item}</span>
              <button
                type="button"
                onClick={() => void handleRemove(index)}
                disabled={saving}
                className="shrink-0 rounded-md p-1 text-muted-foreground transition-colors hover:bg-black/5 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent active:scale-95 disabled:pointer-events-none disabled:opacity-50 dark:hover:bg-white/5"
                aria-label={`Remove ${item}`}
              >
                <X className="h-4 w-4" aria-hidden="true" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
