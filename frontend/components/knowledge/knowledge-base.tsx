"use client"

import { Award, Brain, List, Users } from "lucide-react"
import { useKnowledgeBase } from "@/hooks/use-knowledge-base"
import { KnowledgeError, KnowledgeLoading } from "./knowledge-state"
import { KnowledgeListSection } from "./knowledge-list-section"
import { ProfileSection } from "./profile-section"
import { ToneSection } from "./tone-section"
import {
  isKnowledgeSection,
  type KnowledgeDocument,
  type KnowledgeSection,
} from "./types"

interface KnowledgeBaseProps {
  activeSection: string
}

export function AeolyzerKnowledgeBase({
  activeSection,
}: KnowledgeBaseProps) {
  if (!isKnowledgeSection(activeSection)) {
    return <div>Select a section from the sidebar.</div>
  }
  return (
    <KnowledgeSectionController
      key={activeSection}
      section={activeSection}
    />
  )
}

function KnowledgeSectionController({
  section,
}: {
  section: KnowledgeSection
}) {
  const knowledge = useKnowledgeBase(section)

  return (
    <div
      className="flex-1 w-full overflow-y-auto p-6 font-outfit md:p-8 custom-scrollbar"
      data-testid="knowledge-base"
    >
      <div className="mx-auto max-w-5xl">
        {knowledge.loading ? (
          <KnowledgeLoading section={section} />
        ) : (
          <>
            {knowledge.error && (
              <KnowledgeError
                message={knowledge.error}
                onRetry={knowledge.reload}
              />
            )}
            <KnowledgeSectionContent
              key={`${section}-${knowledge.document.version}`}
              document={knowledge.document}
              saving={knowledge.saving}
              onSave={knowledge.save}
            />
          </>
        )}
      </div>
    </div>
  )
}

interface KnowledgeSectionContentProps {
  document: KnowledgeDocument
  saving: boolean
  onSave: (document: KnowledgeDocument) => Promise<boolean>
}

function KnowledgeSectionContent({
  document,
  saving,
  onSave,
}: KnowledgeSectionContentProps) {
  switch (document.section) {
    case "profile":
      return (
        <ProfileSection
          document={document}
          saving={saving}
          onSave={onSave}
        />
      )
    case "tone":
      return (
        <ToneSection
          document={document}
          saving={saving}
          onSave={onSave}
        />
      )
    case "eeat":
      return (
        <KnowledgeListSection
          title="E-E-A-T Guidelines"
          description="Configure Experience, Expertise, Authoritativeness, and Trustworthiness signals."
          emptyTitle="No E-E-A-T Rules Defined"
          emptyDescription="Add your specific authority markers and trust signals to ensure content meets Google's quality rater guidelines."
          inputLabel="E-E-A-T guideline"
          placeholder="Add an authority marker or trust guideline..."
          actionLabel="Add Guideline"
          icon={Award}
          items={document.eeat?.guidelines ?? []}
          saving={saving}
          onSave={(guidelines) =>
            onSave({ ...document, eeat: { guidelines } })
          }
        />
      )
    case "competitors":
      return (
        <KnowledgeListSection
          title="Competitor Analysis"
          description="Track competitor websites and analyze their content strategies."
          emptyTitle="No competitors added yet"
          emptyDescription="Add a competitor URL to make it available to both AEOlyzer agents."
          inputLabel="Competitor URL"
          placeholder="https://competitor.com"
          actionLabel="Add"
          icon={Users}
          items={document.competitors?.urls ?? []}
          saving={saving}
          onSave={(urls) =>
            onSave({ ...document, competitors: { urls } })
          }
        />
      )
    case "topics":
      return (
        <KnowledgeListSection
          title="Topic Clusters"
          description="Manage your core content pillars and keyword clusters."
          emptyTitle="Map your Content"
          emptyDescription="Define the main topics your AI should focus on for semantic relevance."
          inputLabel="Topic cluster"
          placeholder="Add a topic or content pillar..."
          actionLabel="Create Cluster"
          icon={List}
          items={document.topics?.topics ?? []}
          saving={saving}
          onSave={(topics) =>
            onSave({ ...document, topics: { topics } })
          }
        />
      )
    case "memory":
      return (
        <KnowledgeListSection
          title="Agent Memory"
          description="Manage approved facts and contextual knowledge available to both agents."
          emptyTitle="Memory Vault is Empty"
          emptyDescription="Add an approved brand fact here. AEOlyzer never saves conversational facts silently."
          inputLabel="Approved memory fact"
          placeholder="Add a durable brand or audience fact..."
          actionLabel="Add Fact"
          icon={Brain}
          items={document.memory?.facts ?? []}
          saving={saving}
          onSave={(facts) =>
            onSave({ ...document, memory: { facts } })
          }
        />
      )
  }
}
