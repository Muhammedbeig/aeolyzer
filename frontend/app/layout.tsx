import type { Metadata } from "next"
import { Geist, Rokkitt } from "next/font/google"
import "./globals.css"

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
})

const rokkitt = Rokkitt({
  subsets: ["latin"],
  variable: "--font-display",
  weight: ["400", "500", "600"],
})

export const metadata: Metadata = {
  title: "AEOlyzer — Answer Engine Visibility",
  description: "Understand and improve how AI answer engines see your brand.",
}

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className={`${geist.variable} ${rokkitt.variable}`}>
      <body>{children}</body>
    </html>
  )
}
