"use client"

import { useState, type FormEvent } from "react"
import type {
  KnowledgeDocument,
  PrimaryTone,
} from "./types"

interface ToneSectionProps {
  document: KnowledgeDocument
  saving: boolean
  onSave: (document: KnowledgeDocument) => Promise<boolean>
}

export function ToneSection({
  document,
  saving,
  onSave,
}: ToneSectionProps) {
  const [primaryTone, setPrimaryTone] = useState<PrimaryTone>(
    document.tone?.primary_tone ?? "professional_authoritative",
  )
  const [customInstructions, setCustomInstructions] = useState(
    document.tone?.custom_instructions ?? "",
  )

  const handleSubmit = async (event: FormEvent) => {
    event.preventDefault()
    await onSave({
      ...document,
      tone: {
        primary_tone: primaryTone,
        custom_instructions: customInstructions,
      },
    })
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="max-w-3xl space-y-6"
      data-testid="knowledge-tone"
    >
      <div>
        <h1 className="mb-1 text-xl font-semibold">Brand Tone & Voice</h1>
        <p className="text-sm text-muted-foreground">
          Set the personality, reading level, and style for generated content.
        </p>
      </div>
      <div className="grid gap-6">
        <div className="space-y-2">
          <label htmlFor="knowledge-primary-tone" className="text-sm font-medium">
            Primary Tone
          </label>
          <select
            id="knowledge-primary-tone"
            value={primaryTone}
            onChange={(event) => setPrimaryTone(event.target.value as PrimaryTone)}
            disabled={saving}
            className="w-full cursor-pointer appearance-none rounded-md border border-black/10 bg-transparent px-3 py-2 transition-colors focus:border-accent focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-white/10"
          >
            <option value="professional_authoritative">
              Professional &amp; Authoritative
            </option>
            <option value="conversational_friendly">
              Conversational &amp; Friendly
            </option>
            <option value="academic_technical">Academic &amp; Technical</option>
            <option value="persuasive_direct">Persuasive &amp; Direct</option>
          </select>
        </div>
        <div className="space-y-2">
          <label
            htmlFor="knowledge-tone-instructions"
            className="text-sm font-medium"
          >
            Custom Voice Instructions
          </label>
          <textarea
            id="knowledge-tone-instructions"
            value={customInstructions}
            onChange={(event) => setCustomInstructions(event.target.value)}
            placeholder="e.g. Always use active voice, avoid jargon, use short paragraphs..."
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
            {saving ? "Saving..." : "Save Tone Settings"}
          </button>
        </div>
      </div>
    </form>
  )
}
