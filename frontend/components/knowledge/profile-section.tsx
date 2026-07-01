"use client"

import { useState, type FormEvent } from "react"
import type { KnowledgeDocument } from "./types"

interface ProfileSectionProps {
  document: KnowledgeDocument
  saving: boolean
  onSave: (document: KnowledgeDocument) => Promise<boolean>
}

export function ProfileSection({
  document,
  saving,
  onSave,
}: ProfileSectionProps) {
  const [name, setName] = useState(document.profile?.name ?? "")
  const [description, setDescription] = useState(
    document.profile?.description ?? "",
  )

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    await onSave({
      ...document,
      profile: { name, description },
    })
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="max-w-3xl space-y-6"
      data-testid="knowledge-profile"
    >
      <div>
        <h1 className="mb-1 text-xl font-semibold">Knowledge Base Profile</h1>
        <p className="text-sm text-muted-foreground">
          Manage your main profile details, company information, and primary directives.
        </p>
      </div>

      <div className="grid gap-6">
        <div className="space-y-2">
          <label htmlFor="knowledge-profile-name" className="text-sm font-medium">
            Company/Agent Name
          </label>
          <input
            id="knowledge-profile-name"
            type="text"
            value={name}
            onChange={(event) => setName(event.target.value)}
            placeholder="e.g. AEOlyzer SEO Expert"
            disabled={saving}
            className="w-full rounded-md border border-black/10 bg-transparent px-3 py-2 transition-colors focus:border-accent focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10"
          />
        </div>
        <div className="space-y-2">
          <label
            htmlFor="knowledge-profile-description"
            className="text-sm font-medium"
          >
            Core Description
          </label>
          <textarea
            id="knowledge-profile-description"
            value={description}
            onChange={(event) => setDescription(event.target.value)}
            placeholder="Describe the primary purpose and scope of this knowledge base..."
            disabled={saving}
            className="min-h-32 w-full resize-y rounded-md border border-black/10 bg-transparent px-3 py-2 transition-colors focus:border-accent focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10"
          />
        </div>
        <div className="pt-4">
          <button
            type="submit"
            disabled={saving}
            className="rounded-md bg-accent px-4 py-2 text-sm font-medium text-plum-900 transition-opacity hover:opacity-90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent/40 active:scale-95 disabled:pointer-events-none disabled:opacity-50"
          >
            {saving ? "Saving..." : "Save Profile"}
          </button>
        </div>
      </div>
    </form>
  )
}
