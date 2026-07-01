"use client"

import { useCallback, useEffect, useState } from "react"
import {
  createEmptyKnowledgeDocument,
  type KnowledgeDocument,
  type KnowledgeSection,
} from "@/components/knowledge/types"
import { AeolyzerAPIError } from "@/lib/aeolyzer-api"
import { getKnowledge, updateKnowledge } from "@/lib/knowledge-api"

export function useKnowledgeBase(section: KnowledgeSection) {
  const [document, setDocument] = useState<KnowledgeDocument>(() =>
    createEmptyKnowledgeDocument(section),
  )
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string>()
  const [reloadKey, setReloadKey] = useState(0)

  useEffect(() => {
    const controller = new AbortController()
    getKnowledge(section, controller.signal)
      .then(setDocument)
      .catch((cause: unknown) => {
        if (!(cause instanceof DOMException && cause.name === "AbortError")) {
          setDocument(createEmptyKnowledgeDocument(section))
          setError(errorMessage(cause))
        }
      })
      .finally(() => {
        if (!controller.signal.aborted) {
          setLoading(false)
        }
      })
    return () => controller.abort()
  }, [reloadKey, section])

  const save = useCallback(
    async (next: KnowledgeDocument) => {
      setSaving(true)
      setError(undefined)
      try {
        const updated = await updateKnowledge({
          ...next,
          section,
          version: document.version,
        })
        setDocument(updated)
        return true
      } catch (cause) {
        setError(errorMessage(cause))
        return false
      } finally {
        setSaving(false)
      }
    },
    [document.version, section],
  )

  return {
    document,
    loading,
    saving,
    error,
    save,
    reload: () => {
      setLoading(true)
      setError(undefined)
      setReloadKey((value) => value + 1)
    },
  }
}

function errorMessage(error: unknown) {
  return error instanceof AeolyzerAPIError
    ? error.message
    : "AEOlyzer could not load these settings. Please try again."
}
