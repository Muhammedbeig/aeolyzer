"use client"

import { cn } from "@/lib/utils"

interface AeolyzerLogoProps {
  className?: string
  size?: number
  animate?: boolean
}

// Agentic Cursor Logo
export function AeolyzerLogo({ className, size = 32, animate = false }: AeolyzerLogoProps) {
  return (
    <svg
      viewBox="0 0 32 32"
      width={size}
      height={size}
      className={cn(animate && "animate-pulse", className)}
      style={animate ? { animationDuration: "2s" } : undefined}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      {/* AI Agent Cursor */}
      <path
        d="M 6 6 L 11 26 L 14.5 16.5 L 24 13 Z"
        fill="var(--accent)"
        fillOpacity="0.2"
        stroke="var(--accent)"
        strokeWidth="2.5"
        strokeLinejoin="round"
        strokeLinecap="round"
      />
      {/* AI Spark */}
      <path
        d="M 25 3 C 25 6 27 8 30 8 C 27 8 25 10 25 13 C 25 10 23 8 20 8 C 23 8 25 6 25 3 Z"
        fill="var(--accent)"
      />
    </svg>
  )
}

export function AeolyzerLogoAnimated({ className, size = 32 }: { className?: string; size?: number }) {
  return (
    <svg
      viewBox="0 0 32 32"
      width={size}
      height={size}
      className={cn("animate-pulse", className)}
      style={{ animationDuration: "1.5s" }}
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M 6 6 L 11 26 L 14.5 16.5 L 24 13 Z"
        fill="var(--accent)"
        fillOpacity="0.2"
        stroke="var(--accent)"
        strokeWidth="2.5"
        strokeLinejoin="round"
        strokeLinecap="round"
      />
      <path
        d="M 25 3 C 25 6 27 8 30 8 C 27 8 25 10 25 13 C 25 10 23 8 20 8 C 23 8 25 6 25 3 Z"
        fill="var(--accent)"
      />
    </svg>
  )
}
