import React from 'react'
import { Outlet } from 'react-router-dom'
import { NavMobileFooter, NavMobileHeader } from '../Navigation/Navigation.mobile'
import { NavDesktop } from '../Navigation/Navigation.desktop'
import { useTheme } from '../providers/ThemeProvider'

export const DashboardLayout: React.FC = () => {
  const { colors } = useTheme()

  return (
    <div 
      className="relative min-h-screen w-full font-inter flex flex-col lg:flex-row"
      style={{ 
        backgroundColor: colors.background,
        color: colors.text
      }}
    >
      <div 
        className="fixed inset-0 pointer-events-none opacity-5"
        style={{
          backgroundImage: `radial-gradient(circle at 24px 24px, ${colors.mint} 2px, transparent 0)`,
          backgroundSize: '48px 48px'
        }}
      />

      <div 
        className="hidden lg:block fixed left-0 top-0 h-full w-[300px]"
        style={{ 
          backgroundColor: colors.surface,
          borderRight: `1px solid ${colors.border}`
        }}
      >
        <NavDesktop onOpenSettings={() => {}} />
      </div>

      <div className="flex-1 flex flex-col lg:ml-[300px] w-full lg:w-[calc(100%-300px)] relative">
        <div className="pt-[env(safe-area-inset-top)]" />
        
        <div className="lg:hidden">
          <NavMobileHeader onOpenSettings={() => {}} />
        </div>

        <main className="flex-1">
          <div 
            className="min-w-fit px-4 md:px-6 lg:px-8 py-6 md:py-8 lg:py-10"
            style={{
              paddingLeft: 'max(env(safe-area-inset-left, 16px), 1rem)',
              paddingRight: 'max(env(safe-area-inset-right, 16px), 1rem)'
            }}
          >
            <Outlet />
          </div>
        </main>

        <div className="lg:hidden">
          <NavMobileFooter />
        </div>
      </div>
    </div>
  )
}