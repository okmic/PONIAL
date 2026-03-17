import { useSelector } from 'react-redux'
import type { RootState } from '../../store/store'
import { navMobileItems } from './Navigation.contants'
import { NavLink, useLocation } from 'react-router-dom'
import type { FC } from 'react'
import { Projector, User } from 'lucide-react'
import { useTheme } from '../providers/ThemeProvider'
import Logo from '../Logo/Logo'

type PropsType = {
  onOpenSettings: () => void
}

export const NavMobileHeader: FC<PropsType> = ({ onOpenSettings }) => {
  const user = useSelector((s: RootState) => s.auth.user!)
  const { colors } = useTheme()

  return (
    <div 
      className='w-full mt-2 flex justify-between items-center px-4 py-3'
    >
      <NavLink
        onClick={onOpenSettings}
        to="/settings"
        className="cursor-pointer transition-all hover:scale-105"
      >
        <div 
          className="w-10 h-10 rounded-full flex items-center justify-center"
          style={{ backgroundColor: `${colors.mint}20` }}
        >
          <User className="w-5 h-5" style={{ color: colors.mint }} />
        </div>
      </NavLink>

      <div className='flex flex-col gap-1 justify-center items-center'>
        <Logo 
          animated={true}
          size='sm'
        />
        <span 
          className='text-xs font-medium'
          style={{ color: colors.textSecondary }}
        >
          {user.email}
        </span>
      </div>

      <NavLink
        to="/venue"
        className="cursor-pointer transition-all hover:scale-105"
      >
        <div 
          className="w-10 h-10 rounded-full flex items-center justify-center"
          style={{ backgroundColor: `${colors.mint}20` }}
        >
          <Projector className="w-5 h-5" style={{ color: colors.mint }} />
        </div>
      </NavLink>
    </div>
  )
}

export const NavMobileFooter = () => {
  const user = useSelector((s: RootState) => s.auth.user!)
  const location = useLocation()
  const { colors } = useTheme()

  const isLinkActive = (itemPath: string) => {
    return location.pathname === itemPath
  }

  return (
    <div 
      className='fixed left-1/2 -translate-x-1/2 bottom-4 min-w-[314px] rounded-[31px] flex justify-evenly p-1 backdrop-blur-md'
      style={{ 
        backgroundColor: `${colors.surface}CC`,
        border: `1px solid ${colors.border}`,
        boxShadow: `0 4px 20px ${colors.glow}`
      }}
    >
      {navMobileItems[user.role].map((n) => (
        <NavLink
          key={n.path}
          to={n.path}
          style={({ isActive }) => {
            return {
              backgroundColor: isActive ? `${colors.mint}20` : 'transparent',
              borderRadius: '27px',
              border: isActive ? `1px solid ${colors.mint}` : '1px solid transparent'
            }
          }}
        >
          <div style={{ color: isLinkActive(n.path) ? colors.mint : colors.textSecondary }}>
            {n.icon}
          </div>
          <span 
            className='text-[10px] font-bold uppercase tracking-wide'
            style={{ color: isLinkActive(n.path) ? colors.mint : colors.textSecondary }}
          >
            {n.label}
          </span>
          {isLinkActive(n.path) && (
            <div 
              className="w-1 h-1 rounded-full animate-pulse absolute -bottom-1"
              style={{ backgroundColor: colors.mint }}
            />
          )}
        </NavLink>
      ))}
    </div>
  )
}