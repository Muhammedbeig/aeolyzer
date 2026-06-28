'use client' // -> tells Next.js this code runs in the user's browser, allowing it to switch themes dynamically

import * as React from 'react' // -> brings in React so we can build user interface pieces
import {
  ThemeProvider as NextThemesProvider, // -> imports a tool from 'next-themes' and renames it so it doesn't clash with our own ThemeProvider name
  type ThemeProviderProps, // -> brings in the rules for what data (props) the provider is allowed to expect
} from 'next-themes' // -> an external package that automatically handles dark mode and light mode switching

// -> This is a wrapper component. Anything placed inside it (its "children") will be able to know if the app is in dark mode or light mode
export function ThemeProvider({ children, ...props }: ThemeProviderProps) {
  // - `<NextThemesProvider>` -> the actual tool doing the work of managing the theme
  // - `{...props}` -> whatever extra settings were passed to us, pass them straight through to NextThemesProvider (like a relay runner handing off a baton)
  // - `{children}` -> the rest of our app goes inside here
  return <NextThemesProvider {...props}>{children}</NextThemesProvider>
}
