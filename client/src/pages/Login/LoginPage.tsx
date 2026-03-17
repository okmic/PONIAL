import React, { useState } from 'react'
import { LoginForm } from './LoginForm'
import type { LoginFormData } from './types'
import Logo from '../../components/Logo/Logo'
import apiAuthService from '../../libs/api/api.auth.service'
import { useNotifications } from '../../hooks/useNotifications'

type AuthMode = 'signin' | 'signup'

export const LoginPage: React.FC = () => {
    const [isLoading, setIsLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [mode, setMode] = useState<AuthMode>('signin')
    const notifications = useNotifications()

    const handleLogin = async (data: LoginFormData) => {
        setIsLoading(true)
        setError(null)

        try {
            if (mode === 'signin') {
                await apiAuthService.signin(data)
            } else {
                await apiAuthService.signup({ name: data.name!, email: data.email, password: data.password, role: "user" })
            }
        } catch (err) {
            notifications.error(err instanceof Error ? err.message : 'Ошибка')
            setError(err instanceof Error ? err.message : 'Ошибка')
        } finally {
            setIsLoading(false)
        }
    }

    const toggleMode = () => {
        setMode(mode === 'signin' ? 'signup' : 'signin')
        setError(null)
    }

    return (
        <div className="min-h-screen w-full relative overflow-hidden bg-[#0B1A33] transition-colors duration-500">
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
                            theme="dark"
                            size="lg"
                            animated
                            onClick={() => console.log('ПОНЯЛ!')}
                            className="justify-center"
                        />
                        <p className="mt-4 text-sm uppercase tracking-[0.3em] text-[#F5F0E8]/50">
                            {mode === 'signin' ? 'Вход в систему' : 'Регистрация'}
                        </p>
                    </div>

                    <div className="flex gap-2 p-1 mb-8 rounded-full bg-[#1A1A1A] border border-[#00E5B0]/20">
                        <button
                            onClick={() => setMode('signin')}
                            className={`
                                cursor-pointer flex-1 py-3 px-6 rounded-full font-black tracking-wide transition-all duration-300
                                ${mode === 'signin' 
                                    ? 'bg-[#00E5B0] text-[#0B1A33] shadow-[0_0_20px_rgba(0,229,176,0.3)]' 
                                    : 'text-[#F5F0E8]/60 hover:text-[#F5F0E8]'
                                }
                            `}
                        >
                            ВХОД
                        </button>
                        <button
                            onClick={() => setMode('signup')}
                            className={`
                                cursor-pointer flex-1 py-3 px-6 rounded-full font-black tracking-wide transition-all duration-300
                                ${mode === 'signup' 
                                    ? 'bg-[#00E5B0] text-[#0B1A33] shadow-[0_0_20px_rgba(0,229,176,0.3)]' 
                                    : 'text-[#F5F0E8]/60 hover:text-[#F5F0E8]'
                                }
                            `}
                        >
                            РЕГИСТРАЦИЯ
                        </button>
                    </div>

                    <LoginForm
                        onSubmit={handleLogin}
                        isLoading={isLoading}
                        error={error}
                        mode={mode}
                    />

                    <div className="mt-8 text-center">
                        <button
                            onClick={toggleMode}
                            className="text-sm text-[#F5F0E8]/60 hover:text-[#00E5B0] transition-colors duration-300 underline decoration-dotted underline-offset-4"
                        >
                            {mode === 'signin' 
                                ? 'Нет аккаунта? Создать' 
                                : 'Уже есть аккаунт? Войти'}
                        </button>
                    </div>

                    <div className="mt-12 text-center">
                        <p className="text-xs italic text-[#F5F0E8]/40">
                            "Понял. Нашел. Поехали."
                        </p>
                    </div>
                </div>
            </div>
        </div>
    )
}