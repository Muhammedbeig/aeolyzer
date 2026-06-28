import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  // clsx handles conditional grouping, while twMerge resolves Tailwind specificity conflicts (e.g., p-4 vs p-2) without needing explicit overrides.
  return twMerge(clsx(inputs))
}
