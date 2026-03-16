import React, { useState } from 'react'
import { LoginForm } from './LoginForm'
import type { LoginFormData } from './types'
import Logo from '../../components/Logo/Logo'

export const LoginPage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleLogin = async (data: LoginFormData) => {
    setIsLoading(true)
    setError(null)

    try {
      await new Promise(resolve => setTimeout(resolve, 1500))
      
      if (data.email === 'test@ponial.ru' && data.password === 'ponial2024') {
        console.log('Успешный вход:', data)
      } else {
        throw new Error('Неверный email или пароль')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Ошибка входа')
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className={`
      min-h-screen w-full relative overflow-hidden
      bg-[#0B1A33]
      transition-colors duration-500
    `}>
      <div className="absolute inset-0 opacity-5">
        <div className="w-full h-full" style={{
          backgroundImage: `radial-gradient(circle at 24px 24px, #00E5B0 2px, transparent 0)`,
          backgroundSize: '48px 48px'
        }} />
      </div>
      <div className="absolute top-20 left-10 w-32 h-32 opacity-10">
        <div className="w-full h-full rounded-full bg-[#00E5B0] blur-3xl" />
      </div>
      <div className="absolute bottom-20 right-10 w-64 h-64 opacity-10">
        <div className="w-full h-full rounded-full bg-[#00B8FF] blur-3xl" />
      </div>
      <div className="relative z-10 min-h-screen flex items-center justify-center p-4">
        <div className="w-full max-w-md">
          <div className="text-center mb-12">
            <Logo
              theme={"dark"} 
              size="lg" 
              animated 
              onClick={() => console.log('ПОНЯЛ!')}
              className="justify-center"
            />
            <p className={`
              mt-4 text-sm uppercase tracking-[0.3em]
              text-[#F5F0E8]/50
            `}>
              Вход в систему
            </p>
          </div>
          <LoginForm
            onSubmit={handleLogin}
            isLoading={isLoading}
            error={error}
          />
          <div className="mt-12 text-center">
            <p className={`
              text-xs italic text-[#F5F0E8]/40
            `}>
              "Понял. Нашел. Поехали."
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
