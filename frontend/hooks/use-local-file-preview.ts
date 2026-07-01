"use client"

import { useEffect, useState } from "react"

interface LoadedPreview {
  file: File
  url: string
}

export function useLocalFilePreview(file: File | undefined) {
  const [loadedPreview, setLoadedPreview] = useState<LoadedPreview>()

  useEffect(() => {
    if (!file) {
      return
    }

    let active = true
    const reader = new FileReader()

    reader.addEventListener("load", () => {
      if (active && typeof reader.result === "string") {
        setLoadedPreview({ file, url: reader.result })
      }
    })
    reader.readAsDataURL(file)

    return () => {
      active = false
      if (reader.readyState === FileReader.LOADING) {
        reader.abort()
      }
    }
  }, [file])

  return loadedPreview && loadedPreview.file === file
    ? loadedPreview.url
    : undefined
}
