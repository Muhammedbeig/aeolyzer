// Force client-side execution to isolate SVG animation rendering.
"use client"

import { cn } from "@/lib/utils"
import React, { useState, useEffect } from "react"

interface AeolyzerLogoProps {
  className?: string
  size?: number
  animate?: boolean
}

// Pre-computed SVG path to prevent recalculation overhead during render.
const SPARKLE_PATH = "M 16 0 C 16.8 12, 20 15.5, 29 16 C 20 16.5, 16.8 20, 16 32 C 15.2 20, 12 16.5, 3 16 C 12 15.5, 15.2 12, 16 0 Z"

// Stateless UI component. Uses pure functional rendering to avoid unnecessary reconciliations.
export function AeolyzerLogo({ className, size = 32, animate = false }: AeolyzerLogoProps) {
  return (
    <svg
      viewBox="0 0 32 32"
      width={size}
      height={size}
      className={className}
      fill="none"
    >
      {/* Conditional CSS injection. Dynamically scoped to avoid global stylesheet pollution. */}
      {animate && (
        <style>{`
          @keyframes sparkle-breathe {
            0%, 100% { transform: scale(1); opacity: 1; }
            50% { transform: scale(1.08); opacity: 0.8; }
          }
        `}</style>
      )}

      <defs>
        {/* Static linear gradient definition for SVG fill optimizations. */}
        <linearGradient id="aeo-star-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#ff9a6f" />
          <stop offset="50%" stopColor="#e07b53" />
          <stop offset="100%" stopColor="#c45e2e" />
        </linearGradient>
      </defs>

      <path
        d={SPARKLE_PATH}
        fill="url(#aeo-star-grad)"
        // Inline styles used here since animation is conditionally applied at runtime.
        style={animate ? {
          animation: "sparkle-breathe 3s ease-in-out infinite",
          transformOrigin: "16px 16px",
        } : undefined}
      />
    </svg>
  )
}

// Animated spinner component encapsulating isolated state to prevent parent tree re-renders.
export function AeolyzerLogoAnimated({ className, size = 32 }: { className?: string; size?: number }) {
  // Local state for gradient index to isolate re-render scope strictly to the spinner.
  const [gradState, setGradState] = useState(0)

  // Side effect for continuous animation tick.
  // Returns cleanup function to clear interval on unmount, preventing memory leaks and orphaned timers.
  useEffect(() => {
    const interval = setInterval(() => {
      setGradState((prev) => (prev + 1) % 3)
    }, 1500)
    return () => clearInterval(interval)
  }, [])

  // Static constant array defined within component scope but lightweight enough to bypass memoization overhead.
  const gradients = [
    { start: "#ff2d55", end: "#ff4500" },
    { start: "#ff4500", end: "#25f4ee" },
    { start: "#25f4ee", end: "#ff2d55" },
  ]

  const current = gradients[gradState]

  return (
    <svg
      viewBox="0 0 50 50"
      width={size}
      height={size}
      className={cn(className)}
      fill="none"
    >
      {/* Inline styles for scoped keyframes, avoiding CSS module overhead for pure UI artifacts. */}
      <style>{`
        @keyframes aeo-spinner-rotate {
          100% { transform: rotate(360deg); }
        }
        
        @keyframes aeo-spinner-dash {
          0% {
            stroke-dasharray: 1, 150;
            stroke-dashoffset: 0;
          }
          50% {
            stroke-dasharray: 90, 150;
            stroke-dashoffset: -35;
          }
          100% {
            stroke-dasharray: 90, 150;
            stroke-dashoffset: -124;
          }
        }
      `}</style>

      <defs>
        {/* Hardware-accelerated color transitions. Values bound to local state. */}
        <linearGradient id="aeo-dynamic-spinner" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor={current.start} style={{ transition: 'stop-color 1.5s ease-in-out' }} />
          <stop offset="100%" stopColor={current.end} style={{ transition: 'stop-color 1.5s ease-in-out' }} />
        </linearGradient>
      </defs>

      {/* Hardware-accelerated transforms for smooth rendering without layout thrashing. */}
      <g style={{ animation: "aeo-spinner-rotate 2s linear infinite", transformOrigin: "25px 25px" }}>
        <circle
          cx="25"
          cy="25"
          r="20"
          stroke="url(#aeo-dynamic-spinner)"
          strokeWidth="4"
          strokeLinecap="round"
          fill="none"
          style={{ animation: "aeo-spinner-dash 1.5s ease-in-out infinite" }}
        />
      </g>
    </svg>
  )
}
