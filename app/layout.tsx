import type { Metadata } from 'next'
import { Geist, Geist_Mono, Rokkitt } from 'next/font/google'
import { Analytics } from '@vercel/analytics/next'
import './globals.css'

const _geist = Geist({ subsets: ["latin"] });
const _geistMono = Geist_Mono({ subsets: ["latin"] });
// Using Rokkitt as closest Google Fonts alternative to Bogue Slab
const _rokkitt = Rokkitt({ 
  subsets: ["latin"],
  variable: "--font-display",
  weight: ["100", "200", "300", "400"]
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

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" className={`dark ${_rokkitt.variable}`}>
      <body className="font-sans antialiased" style={{ backgroundColor: "#2b2a27" }}>
        {children}
        <Analytics />
      </body>
    </html>
  )
}
