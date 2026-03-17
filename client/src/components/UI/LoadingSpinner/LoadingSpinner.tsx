import React from 'react'
import Logo from '../../Logo/Logo'

type SpinnerSize = 'sm' | 'md' | 'lg' | 'xl'
type SpinnerTheme = 'dark' | 'light'

interface LoadingSpinnerProps {
  size?: SpinnerSize
  theme?: SpinnerTheme
  text?: string
  fullScreen?: boolean
  className?: string
}

const sizeMap = {
  sm: { container: 'w-24 h-24', logo: 'sm', text: 'text-sm' },
  md: { container: 'w-32 h-32', logo: 'md', text: 'text-base' },
  lg: { container: 'w-48 h-48', logo: 'lg', text: 'text-lg' },
  xl: { container: 'w-64 h-64', logo: 'xl', text: 'text-xl' }
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  theme = 'dark',
  text = 'ПОНЯЛ',
  fullScreen = false,
  className = ''
}) => {
  const spinnerContent = (
    <div className={`
      flex flex-col items-center justify-center gap-6
      ${sizeMap[size].container}
    `}>
      <div className="relative">
        <div className="absolute inset-0 rounded-full bg-[#00E5B0] blur-xl opacity-20 animate-pulse" />
        <Logo 
          theme={theme}
          size={sizeMap[size].logo as any}
          animated
        />
      </div>

      <div className="text-center">
        <p className={`
          font-black tracking-widest uppercase
          ${sizeMap[size].text}
          ${theme === 'dark' ? 'text-[#F5F0E8]' : 'text-[#0B1A33]'}
        `}>
          {text}
          <span className="animate-pulse">...</span>
        </p>
        <p className={`
          text-xs mt-2 opacity-40
          ${theme === 'dark' ? 'text-[#F5F0E8]' : 'text-[#0B1A33]'}
        `}>
          ищу деталь
        </p>
      </div>

      <div className="flex gap-1">
        {[1,2,3].map(i => (
          <div
            key={i}
            className="w-1 h-4 bg-[#00E5B0] rounded-full animate-pulse"
            style={{ animationDelay: `${i * 0.2}s` }}
          />
        ))}
      </div>
    </div>
  )

  if (fullScreen) {
    return (
      <div className={`
        fixed inset-0 z-50
        flex items-center justify-center
        ${theme === 'dark' ? 'bg-[#0B1A33]' : 'bg-[#F5F0E8]'}
        ${className}
      `}>
        <div className="absolute inset-0 opacity-5">
          <div className="w-full h-full" style={{
            backgroundImage: 'radial-gradient(circle at 24px 24px, #00E5B0 2px, transparent 0)',
            backgroundSize: '48px 48px'
          }} />
        </div>
        {spinnerContent}
      </div>
    )
  }

  return spinnerContent
}
