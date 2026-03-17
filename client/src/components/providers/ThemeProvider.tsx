import React, { createContext, useContext, type ReactNode } from 'react'

interface ThemeContextType {
  colors: {
    navy: '#0B1A33'
    mint: '#00E5B0'
    paper: '#F5F0E8'
    rubber: '#1A1A1A'
    red: '#FF3B30'
    background: '#0B1A33'
    surface: '#1A1A1A'
    text: '#F5F0E8'
    textSecondary: '#F5F0E8B3'
    border: '#00E5B033'
    glow: '#00E5B04D'
  }
  radius: {
    card: '24px'
    full: '9999px'
  }
  font: {
    family: 'Inter, system-ui, sans-serif'
    weight: {
      normal: 400
      bold: 700
      black: 900
    }
    tracking: {
      normal: '0'
      wide: '0.05em'
      wider: '0.1em'
    }
  }
  animation: {
    duration: {
      fast: '150ms'
      normal: '300ms'
      slow: '500ms'
    }
    easing: 'ease-in-out'
  }
}

const theme: ThemeContextType = {
  colors: {
    navy: '#0B1A33',
    mint: '#00E5B0',
    paper: '#F5F0E8',
    rubber: '#1A1A1A',
    red: '#FF3B30',
    background: '#0B1A33',
    surface: '#1A1A1A',
    text: '#F5F0E8',
    textSecondary: '#F5F0E8B3',
    border: '#00E5B033',
    glow: '#00E5B04D'
  },
  radius: {
    card: '24px',
    full: '9999px'
  },
  font: {
    family: 'Inter, system-ui, sans-serif',
    weight: {
      normal: 400,
      bold: 700,
      black: 900
    },
    tracking: {
      normal: '0',
      wide: '0.05em',
      wider: '0.1em'
    }
  },
  animation: {
    duration: {
      fast: '150ms',
      normal: '300ms',
      slow: '500ms'
    },
    easing: 'ease-in-out'
  }
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined)

export const ThemeProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  return (
    <ThemeContext.Provider value={theme}>
      {children}
    </ThemeContext.Provider>
  )
}

export const useTheme = () => {
  const context = useContext(ThemeContext)
  if (!context) {
    throw new Error('useTheme must be used within ThemeProvider')
  }
  return context
}