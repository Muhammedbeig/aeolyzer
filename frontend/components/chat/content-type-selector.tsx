import {
  FileText,
  Linkedin,
  PenLine,
  Rocket,
  Youtube,
  type LucideIcon,
} from "lucide-react"
import type { ContentType } from "./types"
import { cn } from "@/lib/utils"

interface ContentTypeOption {
  value: ContentType
  label: string
  icon: LucideIcon
}

const CONTENT_TYPES: ContentTypeOption[] = [
  { value: "article", label: "Article", icon: FileText },
  { value: "blog_post", label: "Blog Post", icon: PenLine },
  { value: "linkedin_post", label: "LinkedIn", icon: Linkedin },
  { value: "youtube_description", label: "YouTube Desc", icon: Youtube },
  { value: "product_description", label: "Product Desc", icon: Rocket },
]

interface ContentTypeSelectorProps {
  value: ContentType
  disabled: boolean
  onChange: (value: ContentType) => void
}

export function ContentTypeSelector({
  value,
  disabled,
  onChange,
}: ContentTypeSelectorProps) {
  return (
    <div
      className="mb-4 mt-4 flex flex-wrap items-center justify-center gap-1.5"
      data-testid="content-type-selector"
      aria-label="Content type"
    >
      {CONTENT_TYPES.map((option) => {
        const Icon = option.icon
        const selected = option.value === value
        return (
          <button
            key={option.value}
            type="button"
            onClick={() => onChange(option.value)}
            disabled={disabled}
            aria-pressed={selected}
            className={cn(
              "inline-flex items-center gap-1.5 whitespace-nowrap rounded-lg border-[0.5px] px-2.5 py-1.5 text-[12px] font-medium transition-all duration-150 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent active:scale-95 disabled:pointer-events-none disabled:opacity-50",
              selected
                ? "border-accent/50 bg-accent/10 text-accent"
                : "border-black/10 bg-white text-plum-500 hover:border-black/20 hover:bg-sand-50 hover:text-plum-700 dark:border-white/10 dark:bg-card dark:text-foreground dark:hover:border-white/20 dark:hover:bg-accent/10 dark:hover:text-accent",
            )}
          >
            <Icon className="h-3.5 w-3.5" aria-hidden="true" />
            {option.label}
          </button>
        )
      })}
    </div>
  )
}
