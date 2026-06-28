// Forces client-side boundary to prevent hydration mismatch for theme selection, ensuring isolated UI state.
'use client'

import * as React from 'react'
import {
  ThemeProvider as NextThemesProvider,
  type ThemeProviderProps,
} from 'next-themes'

// Wraps application root to isolate theme context state.
// Prevents full app re-renders when theme toggles by delegating DOM mutation to next-themes.
export function ThemeProvider({ children, ...props }: ThemeProviderProps) {
  // Prop-drilling explicitly passed through to NextThemesProvider.
  // Children are passed untouched to prevent unnecessary React re-renders within the provider.
  return <NextThemesProvider {...props}>{children}</NextThemesProvider>
}
