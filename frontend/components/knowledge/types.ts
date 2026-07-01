export type KnowledgeSection =
  | "profile"
  | "eeat"
  | "competitors"
  | "topics"
  | "tone"
  | "memory"

export interface KnowledgeProfile {
  name: string
  description: string
}

export interface KnowledgeEEAT {
  guidelines: string[]
}

export interface KnowledgeCompetitors {
  urls: string[]
}

export interface KnowledgeTopics {
  topics: string[]
}

export type PrimaryTone =
  | "professional_authoritative"
  | "conversational_friendly"
  | "academic_technical"
  | "persuasive_direct"

export interface KnowledgeTone {
  primary_tone: PrimaryTone
  custom_instructions: string
}

export interface KnowledgeMemory {
  facts: string[]
}

export interface KnowledgeDocument {
  section: KnowledgeSection
  version: number
  profile?: KnowledgeProfile
  eeat?: KnowledgeEEAT
  competitors?: KnowledgeCompetitors
  topics?: KnowledgeTopics
  tone?: KnowledgeTone
  memory?: KnowledgeMemory
  updated_at?: string
}

export function createEmptyKnowledgeDocument(
  section: KnowledgeSection,
): KnowledgeDocument {
  const document: KnowledgeDocument = { section, version: 0 }
  switch (section) {
    case "profile":
      document.profile = { name: "", description: "" }
      break
    case "eeat":
      document.eeat = { guidelines: [] }
      break
    case "competitors":
      document.competitors = { urls: [] }
      break
    case "topics":
      document.topics = { topics: [] }
      break
    case "tone":
      document.tone = {
        primary_tone: "professional_authoritative",
        custom_instructions: "",
      }
      break
    case "memory":
      document.memory = { facts: [] }
      break
  }
  return document
}

export function isKnowledgeSection(value: string): value is KnowledgeSection {
  return [
    "profile",
    "eeat",
    "competitors",
    "topics",
    "tone",
    "memory",
  ].includes(value)
}
