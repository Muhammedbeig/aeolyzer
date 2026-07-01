import { AlertCircle, LoaderCircle } from "lucide-react"

interface KnowledgeLoadingProps {
  section: string
}

export function KnowledgeLoading({ section }: KnowledgeLoadingProps) {
  return (
    <div
      className="flex min-h-64 items-center justify-center text-muted-foreground"
      data-testid={`knowledge-${section}-loading`}
    >
      <LoaderCircle className="h-5 w-5 animate-spin" aria-hidden="true" />
      <span className="ml-2 text-sm">Loading settings...</span>
    </div>
  )
}

interface KnowledgeErrorProps {
  message: string
  onRetry: () => void
}

export function KnowledgeError({ message, onRetry }: KnowledgeErrorProps) {
  return (
    <div
      className="mb-5 flex items-center justify-between gap-4 rounded-lg bg-destructive/10 px-4 py-3 text-sm text-destructive"
      role="alert"
      data-testid="knowledge-error"
    >
      <span className="flex items-center gap-2">
        <AlertCircle className="h-4 w-4 shrink-0" aria-hidden="true" />
        {message}
      </span>
      <button
        type="button"
        onClick={onRetry}
        className="rounded-md px-2 py-1 font-medium transition-colors hover:bg-destructive/10 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent active:scale-95"
      >
        Retry
      </button>
    </div>
  )
}
