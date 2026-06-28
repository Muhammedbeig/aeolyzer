import type { Metadata } from "next"
import { Geist, Rokkitt } from "next/font/google"
import "./globals.css"

// Instantiating Next/font at module level caches the font payload to avoid FOUT and layout shifts during navigation.
const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
})

const rokkitt = Rokkitt({
  subsets: ["latin"],
  variable: "--font-display",
  weight: ["400", "500", "600"],
})

// Static metadata export enables build-time evaluation, preventing runtime overhead for SEO tag generation.
export const metadata: Metadata = {
  title: "AEOlyzer — Answer Engine Visibility",
  description: "Understand and improve how AI answer engines see your brand.",
}

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  // Injecting CSS variables directly into HTML node forces style recalculation upfront, isolating font loading state from child components.
  return (
    <html lang="en" className={`${geist.variable} ${rokkitt.variable}`}>
      <body>{children}</body>
    </html>
  )
}
