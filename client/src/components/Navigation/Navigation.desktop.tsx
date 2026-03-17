import { User } from 'lucide-react'
import { useSelector } from 'react-redux'
import type { RootState } from '../../store/store'
import { NavLink, useLocation } from 'react-router-dom'
import { navDesktopItems } from './Navigation.contants'
import { type FC } from 'react'
import Logo from '../Logo/Logo'
import { useTheme } from '../providers/ThemeProvider'

type PropsType = {
  onOpenSettings: () => void
}

export const NavDesktop: FC<PropsType> = ({ onOpenSettings }) => {
  const user = useSelector((s: RootState) => s.auth.user!)
  const location = useLocation()
  const { colors } = useTheme()

  return (
    <nav 
      className="hidden lg:flex fixed top-0 left-0 h-screen w-[300px] flex-col z-10"
      style={{ 
        backgroundColor: colors.surface,
        borderRight: `1px solid ${colors.border}`
      }}
    >
      <div className="px-8 pt-12 flex flex-col items-start gap-1.5">
        <div className="flex items-center gap-3">
          <Logo size='lg' theme="dark" />
        </div>
      </div>

      <div className="mb-8" />

      <ul className="flex flex-col flex-1 overflow-y-auto px-5 gap-2">
        {navDesktopItems[user.role].map((item) => (
          <li key={item.path}>
            <NavLink
              to={item.path}
              className="transition-all duration-200"
              style={({ isActive }) => {
                const shouldBeActive = item.path === '/reception'
                  ? isActive || location.pathname === '/'
                  : isActive

                return {
                  backgroundColor: shouldBeActive ? `${colors.mint}20` : 'transparent',
                  border: shouldBeActive ? `2px solid ${colors.mint}` : `1px solid transparent`,
                  boxShadow: shouldBeActive ? `0 0 20px ${colors.glow}` : 'none',
                  color: shouldBeActive ? colors.mint : colors.text
                }
              }}
            >
              <div className="w-5 h-5 flex items-center justify-center" style={{ color: 'inherit' }}>
                {item.icon}
              </div>
              <span className="text-base font-bold tracking-wide">{item.label}</span>
            </NavLink>
          </li>
        ))}
      </ul>

      <div className="p-4 text-center" style={{ borderColor: `${colors.text}1A` }}>
         <NavLink
          onClick={onOpenSettings}
          to="/settings"
          className="flex items-center ml-5 gap-3 w-full cursor-pointer transition-all duration-200 group"
        >
          <div 
            className="w-9 h-9 rounded-lg flex items-center justify-center flex-shrink-0 transition-all duration-200"
            style={{ backgroundColor: colors.mint }}
          >
            <User className="w-5 h-5" style={{ color: colors.navy }} />
          </div>
          <div className="flex flex-col min-w-0">
        <p className="text-xs uppercase tracking-wider" style={{ color: colors.textSecondary }}>
          Понял. Нашел. Поехали.
        </p>
          </div>
          <div className="ml-auto w-1 h-6 rounded-full opacity-0 transition-opacity" style={{ backgroundColor: colors.mint }} />
        </NavLink>
      </div>
    </nav>
  )
}
