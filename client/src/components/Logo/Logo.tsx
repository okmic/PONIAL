import React, { useState, useEffect } from 'react'
import type { LogoProps, LogoTheme, LogoSize } from './types'

const sizeMap: Record<LogoSize, { text: string; bulb: string; gap: string }> = {
  sm: { text: 'text-2xl', bulb: 'w-6 h-8', gap: 'gap-0.5' },
  md: { text: 'text-4xl', bulb: 'w-8 h-10', gap: 'gap-1' },
  lg: { text: 'text-6xl', bulb: 'w-12 h-16', gap: 'gap-2' },
  xl: { text: 'text-8xl', bulb: 'w-16 h-20', gap: 'gap-3' }
}

export const Logo: React.FC<LogoProps> = ({
  theme = 'dark',
  size = 'lg',
  animated = true,
  onClick,
  className = ''
}) => {
  const [isHovered, setIsHovered] = useState(false)
  const [glowIntensity, setGlowIntensity] = useState(0)

  useEffect(() => {
    if (!animated) return

    const interval = setInterval(() => {
      setGlowIntensity(prev => (prev === 0 ? 1 : 0))
    }, 2000)

    return () => clearInterval(interval)
  }, [animated])

  const themeStyles: Record<LogoTheme, { text: string; bulb: string; glow: string; bg: string }> = {
    dark: {
      text: 'text-[#F5F0E8]',
      bulb: 'bg-gradient-to-br from-[#00E5B0] to-[#00B8FF]',
      glow: 'shadow-[0_0_20px_rgba(0,229,176,0.5)]',
      bg: 'bg-[#0B1A33]'
    },
    light: {
      text: 'text-[#0B1A33]', 
      bulb: 'bg-gradient-to-br from-[#00E5B0] to-[#00B8FF]',
      glow: 'shadow-[0_0_30px_rgba(0,229,176,0.3)]',
      bg: 'bg-[#F5F0E8]'
    }
  }

  const currentSize = sizeMap[size]
  const currentTheme = themeStyles[theme]

  return (
    <div
      className={`
        flex items-center ${currentSize.gap} select-none
        ${onClick ? 'cursor-pointer' : ''}
        ${className}
      `}
      onClick={onClick}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
      role={onClick ? 'button' : 'img'}
      aria-label="ПОНЯЛ логотип"
    >
      <span className={`
        font-black tracking-tight
        ${currentSize.text}
        ${currentTheme.text}
        transition-all duration-300
        drop-shadow-md
        ${isHovered ? 'scale-105' : 'scale-100'}
      `}>
        ПОН
      </span>

      <div className="relative flex items-center justify-center">
        <div
          className={`
            relative ${currentSize.bulb} rounded-t-full
            ${currentTheme.bulb}
            transition-all duration-500 ease-out
            ${isHovered || glowIntensity ? 'scale-110' : 'scale-100'}
            ${currentTheme.glow}
          `}
        >
          <div className="absolute inset-0 flex items-center justify-center">
            <div className={`
              w-1/2 h-1/2 rounded-full
              bg-[#0B1A33] bg-opacity-40
              flex items-center justify-center
              transition-all duration-300
              ${isHovered ? 'rotate-180' : 'rotate-0'}
            `}>
              <div className="w-1/2 h-1/2 bg-[#0B1A33] bg-opacity-60 rounded-full" />
            </div>
          </div>

          <div className={`
            absolute top-1 left-1 w-2 h-2
            bg-white rounded-full opacity-70
            transition-opacity duration-300
            ${isHovered ? 'opacity-90' : 'opacity-50'}
          `} />
        </div>
        {animated && isHovered && (
          <>
            <div className="absolute -top-2 -right-2 w-2 h-2 bg-[#00E5B0] rounded-full animate-ping" />
            <div className="absolute -bottom-1 -left-1 w-1 h-1 bg-[#00B8FF] rounded-full animate-pulse" />
          </>
        )}
      </div>

      <span className={`
        font-black tracking-tight
        ${currentSize.text}
        ${currentTheme.text}
        transition-all duration-300
        drop-shadow-md
        ${isHovered ? 'translate-y-[-2px]' : 'translate-y-0'}
      `}>
        Л
      </span>
    </div>
  )
}

export default Logo
