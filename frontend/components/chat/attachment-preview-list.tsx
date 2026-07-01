"use client"

import { X } from "lucide-react"
import { useLocalFilePreview } from "@/hooks/use-local-file-preview"

interface AttachmentPreviewListProps {
  files: File[]
  disabled: boolean
  onRemove: (index: number) => void
}

export function AttachmentPreviewList({
  files,
  disabled,
  onRemove,
}: AttachmentPreviewListProps) {
  if (files.length === 0) {
    return null
  }
  return (
    <div
      className="flex max-w-full gap-2 overflow-x-auto px-3 pt-3 custom-scrollbar"
      data-testid="attachment-preview-list"
      aria-label="Selected attachments"
    >
      {files.map((file, index) => (
        <AttachmentPreviewCard
          key={`${file.name}-${file.size}-${file.lastModified}-${index}`}
          file={file}
          disabled={disabled}
          onRemove={() => onRemove(index)}
        />
      ))}
    </div>
  )
}

interface AttachmentPreviewCardProps {
  file: File
  disabled: boolean
  onRemove: () => void
}

function AttachmentPreviewCard({
  file,
  disabled,
  onRemove,
}: AttachmentPreviewCardProps) {
  const isImage = file.type.startsWith("image/")
  const previewURL = useLocalFilePreview(isImage ? file : undefined)

  return (
    <div
      className="group relative h-28 w-28 shrink-0 overflow-hidden rounded-xl bg-muted text-foreground"
      data-testid="attachment-preview-card"
    >
      {previewURL ? (
        // Local attachment bytes must stay in-browser and bypass the remote image optimizer.
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={previewURL}
          alt={`Preview of ${file.name}`}
          width={112}
          height={112}
          className="h-full w-full object-cover"
        />
      ) : (
        <div className="flex h-full flex-col justify-between p-3">
          <span className="text-sm font-medium text-muted-foreground">
            {fileExtension(file)}
          </span>
          <div className="min-w-0">
            <p className="truncate text-sm font-medium" title={file.name}>
              {file.name}
            </p>
            <p className="text-xs text-muted-foreground">{formatBytes(file.size)}</p>
          </div>
        </div>
      )}
      {previewURL && (
        <div className="absolute inset-x-0 bottom-0 bg-black/60 px-2 py-1.5 text-white">
          <p className="truncate text-xs font-medium" title={file.name}>
            {file.name}
          </p>
        </div>
      )}
      <button
        type="button"
        onClick={onRemove}
        disabled={disabled}
        className="absolute right-1.5 top-1.5 flex h-6 w-6 items-center justify-center rounded-full bg-background/95 text-foreground transition-colors hover:bg-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent active:scale-95 disabled:pointer-events-none disabled:opacity-50"
        aria-label={`Remove ${file.name}`}
      >
        <X className="h-3.5 w-3.5" aria-hidden="true" />
      </button>
    </div>
  )
}

function fileExtension(file: File) {
  const extension = file.name.split(".").pop()
  if (extension && extension !== file.name && extension.length <= 8) {
    return extension.toUpperCase()
  }
  return file.type.split("/").pop()?.toUpperCase() ?? "FILE"
}

function formatBytes(bytes: number) {
  if (bytes < 1_024) {
    return `${bytes} B`
  }
  if (bytes < 1_048_576) {
    return `${Math.ceil(bytes / 1_024)} KB`
  }
  return `${(bytes / 1_048_576).toFixed(1)} MB`
}
