import type {
  KnowledgeDocument,
  KnowledgeSection,
} from "@/components/knowledge/types"
import { aeolyzerRequest } from "@/lib/aeolyzer-api"

export function getKnowledge(
  section: KnowledgeSection,
  signal?: AbortSignal,
): Promise<KnowledgeDocument> {
  return aeolyzerRequest<KnowledgeDocument>(
    `/v1/knowledge/${encodeURIComponent(section)}`,
    { signal },
  )
}

export function updateKnowledge(
  document: KnowledgeDocument,
): Promise<KnowledgeDocument> {
  return aeolyzerRequest<KnowledgeDocument>(
    `/v1/knowledge/${encodeURIComponent(document.section)}`,
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(updateBody(document)),
    },
  )
}

function updateBody(document: KnowledgeDocument) {
  const payload = sectionPayload(document)
  return {
    version: document.version,
    approved: true,
    [document.section]: payload,
  }
}

function sectionPayload(document: KnowledgeDocument) {
  switch (document.section) {
    case "profile":
      return document.profile
    case "eeat":
      return document.eeat
    case "competitors":
      return document.competitors
    case "topics":
      return document.topics
    case "tone":
      return document.tone
    case "memory":
      return document.memory
  }
}
