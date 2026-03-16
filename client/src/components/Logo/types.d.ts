export type LogoTheme = 'dark' | 'light'
export type LogoSize = 'sm' | 'md' | 'lg' | 'xl'

export interface LogoProps {
  theme?: LogoTheme
  size?: LogoSize
  animated?: boolean
  onClick?: () => void
  className?: string
}