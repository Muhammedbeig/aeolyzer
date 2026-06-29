"use client"

import { cn } from "@/lib/utils"

interface AeolyzerLogoProps {
  className?: string
  size?: number
  animate?: boolean
}

// Claude-style sunburst logo
export function AeolyzerLogo({ className, size = 32, animate = false }: AeolyzerLogoProps) {
  const rayCount = 8
  const rays = Array.from({ length: rayCount }, (_, i) => {
    const angle = (i * 360) / rayCount
    return angle
  })

  return (
    <svg
      viewBox="0 0 32 32"
      width={size}
      height={size}
      className={cn(animate && "animate-spin", className)}
      style={{ animationDuration: "3s" }}
      fill="none"
    >
      {rays.map((angle, i) => (
        <line
          key={i}
          x1="16"
          y1="16"
          x2="16"
          y2={i % 2 === 0 ? "4" : "6"}
          stroke="#e07b53"
          strokeWidth={i % 2 === 0 ? "2.5" : "2"}
          strokeLinecap="round"
          transform={`rotate(${angle} 16 16)`}
        />
      ))}
    </svg>
  )
}

export function AeolyzerLogoAnimated({ className, size = 32 }: { className?: string; size?: number }) {
  const rayCount = 8
  const rays = Array.from({ length: rayCount }, (_, i) => {
    const angle = (i * 360) / rayCount
    return angle
  })

  return (
    <svg
      viewBox="0 0 32 32"
      width={size}
      height={size}
      className={cn("animate-spin", className)}
      style={{ animationDuration: "2s" }}
      fill="none"
    >
      {rays.map((angle, i) => (
        <line
          key={i}
          x1="16"
          y1="16"
          x2="16"
          y2={i % 2 === 0 ? "4" : "6"}
          stroke="#e07b53"
          strokeWidth={i % 2 === 0 ? "2.5" : "2"}
          strokeLinecap="round"
          transform={`rotate(${angle} 16 16)`}
        />
      ))}
    </svg>
  )
}
