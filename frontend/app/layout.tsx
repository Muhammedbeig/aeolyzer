import type { Metadata } from 'next'
import { Geist, Geist_Mono, Rokkitt, Outfit, Plus_Jakarta_Sans } from 'next/font/google'
import { Analytics } from '@vercel/analytics/next'
import './globals.css'

const _geist = Geist({ 
  subsets: ["latin"],
  variable: "--font-geist-sans"
});
const _geistMono = Geist_Mono({ 
  subsets: ["latin"],
  variable: "--font-geist-mono"
});
// Using Rokkitt as closest Google Fonts alternative to Bogue Slab
const _rokkitt = Rokkitt({ 
  subsets: ["latin"],
  variable: "--font-display",
  weight: ["100", "200", "300", "400"]
});
const _outfit = Outfit({
  subsets: ["latin"],
  variable: "--font-outfit",
});
const _jakarta = Plus_Jakarta_Sans({
  subsets: ["latin"],
  variable: "--font-jakarta",
});

export const metadata: Metadata = {
  title: 'AEOlyzer',
  description: 'Talk to AEOlyzer, your intelligent AI assistant',
  generator: 'v0.app',
  icons: {
    icon: [
      {
        url: '/icon-light-32x32.png',
        media: '(prefers-color-scheme: light)',
      },
      {
        url: '/icon-dark-32x32.png',
        media: '(prefers-color-scheme: dark)',
      },
      {
        url: '/icon.svg',
        type: 'image/svg+xml',
      },
    ],
    apple: '/apple-icon.png',
  },
}

import { ThemeProvider } from "@/components/theme-provider"

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className={`${_geist.variable} ${_geistMono.variable} ${_rokkitt.variable} ${_outfit.variable} ${_jakarta.variable}`} suppressHydrationWarning>
      <body className="font-sans antialiased bg-background">
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {children}
          <Analytics />
        </ThemeProvider>
      </body>
    </html>
  )
}
