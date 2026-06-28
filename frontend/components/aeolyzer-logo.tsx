"use client" // -> tells Next.js this code must run in the user's browser, not on the server

import { cn } from "@/lib/utils" // -> imports a helper function 'cn' that safely merges CSS class names together
import React, { useState, useEffect } from "react" // -> import react hooks for the loading spinner

// -> defines the shape of the data (props) the AeolyzerLogo block expects to receive
// - `className` -> optional extra CSS classes for custom styling
// - `size` -> optional number to set how big the logo is
// - `animate` -> optional true/false switch to decide if the logo should breathe/shimmer
interface AeolyzerLogoProps {
  className?: string
  size?: number
  animate?: boolean
}

/*
  AEOlyzer North Star logo — a unique sparkle-inspired brand mark.

  The shape is a vertically elongated 4-pointed star drawn with cubic bezier curves.
  It's inspired by modern sparkle aesthetics but differentiated in several key ways:

  1. Vertically elongated — top/bottom tips stretch further than left/right (portrait, not square)
  2. Thinner arms — the points taper sharply, creating a sleek, needle-like silhouette
  3. Inner accent dot — a small circle at the center that anchors the mark

  👉 Think: A compass north star or guiding light — fitting for an "Answer Engine" tool
  that helps users navigate search visibility.

  The gradient flows diagonally from warm peach (#ff9a6f) through the brand
  terracotta (#e07b53) down to a deep amber (#c45e2e), giving it depth and warmth.
*/

// -> The SVG path data for AEOlyzer's north-star sparkle
// 🧩 How the path creates thin, elegant arms:
//   - Top tip: (16, 0)   — stretches to the very top edge
//   - Right tip: (29, 16) — horizontal reach
//   - Bottom tip: (16, 32) — stretches to the very bottom edge
//   - Left tip: (3, 16)   — horizontal reach
//   - Control points (16.8/15.2) hug the center axis tightly, making the arms
//     very thin and needle-like near the body before blooming at the tips
//   Z → close the path
const SPARKLE_PATH = "M 16 0 C 16.8 12, 20 15.5, 29 16 C 20 16.5, 16.8 20, 16 32 C 15.2 20, 12 16.5, 3 16 C 12 15.5, 15.2 12, 16 0 Z"

// -> This is the main logo component — AEOlyzer's north-star sparkle mark
export function AeolyzerLogo({ className, size = 32, animate = false }: AeolyzerLogoProps) {
  return (
    <svg
      viewBox="0 0 32 32"
      width={size}
      height={size}
      className={className}
      fill="none"
    >
      {/* -> Breathing animation keyframes (only injected when animate is on) */}
      {animate && (
        <style>{`
          @keyframes sparkle-breathe {
            0%, 100% { transform: scale(1); opacity: 1; }
            50% { transform: scale(1.08); opacity: 0.8; }
          }
        `}</style>
      )}

      <defs>
        {/* -> A flowing diagonal gradient from warm peach to deep terracotta */}
        <linearGradient id="aeo-star-grad" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#ff9a6f" />
          <stop offset="50%" stopColor="#e07b53" />
          <stop offset="100%" stopColor="#c45e2e" />
        </linearGradient>
      </defs>

      {/* -> The north-star shape, filled with the gradient */}
      <path
        d={SPARKLE_PATH}
        fill="url(#aeo-star-grad)"
        style={animate ? {
          animation: "sparkle-breathe 3s ease-in-out infinite",
          transformOrigin: "16px 16px",
        } : undefined}
      />
    </svg>
  )
}

// -> A loading spinner inspired by the QuickNews app — a rotating circle with a stretching
//    dash stroke that alternates between 3 color gradients every 1.5 seconds.
// 👉 Think: Like the sleek loading spinner in the QuickNews repo that shifts colors 
//    while it spins, acting as a mesmerizing "thinking" indicator
export function AeolyzerLogoAnimated({ className, size = 32 }: { className?: string; size?: number }) {
  // -> State to track which of the 3 gradients is currently active (0, 1, or 2)
  const [gradState, setGradState] = useState(0)

  // -> Cycle through the 3 gradient states every 1.5 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      setGradState((prev) => (prev + 1) % 3)
    }, 1500)
    return () => clearInterval(interval)
  }, [])

  // -> The 3 gradient color pairs that cycle (matching the QuickNews aesthetic but with AEOlyzer colors)
  // We use warm tones that fit our brand
  const gradients = [
    { start: "#ff2d55", end: "#ff4500" },   // Pink/Orange (like State One)
    { start: "#ff4500", end: "#25f4ee" },   // Orange/Cyan (like State Two)
    { start: "#25f4ee", end: "#ff2d55" },   // Cyan/Pink (like State Three)
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
      <style>{`
        /* The main rotation animation for the whole SVG group */
        @keyframes aeo-spinner-rotate {
          100% { transform: rotate(360deg); }
        }
        
        /* The dash animation that makes the stroke stretch and squish */
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
        {/* We update the stop colors dynamically based on the current state */}
        <linearGradient id="aeo-dynamic-spinner" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor={current.start} style={{ transition: 'stop-color 1.5s ease-in-out' }} />
          <stop offset="100%" stopColor={current.end} style={{ transition: 'stop-color 1.5s ease-in-out' }} />
        </linearGradient>
      </defs>

      {/* The rotating group that holds our circle */}
      <g style={{ animation: "aeo-spinner-rotate 2s linear infinite", transformOrigin: "25px 25px" }}>
        {/* The circle with the stretching dash stroke */}
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
