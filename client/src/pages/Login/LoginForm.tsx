import React, { useState, useEffect } from 'react'
import type { LoginFormProps, LoginFormData, LoginFormErrors } from './types'

export const LoginForm: React.FC<LoginFormProps> = ({
  onSubmit,
  isLoading = false,
  error = null,
  mode = 'signin'
}) => {
  const [formData, setFormData] = useState<LoginFormData>({
    email: '',
    password: '',
    name: '',
    confirmPassword: '',
    remember: false
  })

  const [errors, setErrors] = useState<LoginFormErrors>({})
  const [touched, setTouched] = useState<Set<string>>(new Set())
  const [shake, setShake] = useState(false)
  const [focusedField, setFocusedField] = useState<string | null>(null)

  const validateField = (name: string, value: string): string | undefined => {
    switch (name) {
      case 'email':
        if (!value) return 'Email обязателен'
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
          return 'Неверный формат email'
        }
        break
      case 'password':
        if (!value) return 'Пароль обязателен'
        if (value.length < 6) return 'Минимум 6 символов'
        break
      case 'name':
        if (mode === 'signup' && !value) return 'Имя обязательно'
        break
      case 'confirmPassword':
        if (mode === 'signup' && !value) return 'Подтвердите пароль'
        if (mode === 'signup' && value !== formData.password) return 'Пароли не совпадают'
        break
    }
    return undefined
  }

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value, type, checked } = e.target
    const newValue = type === 'checkbox' ? checked : value

    setFormData(prev => ({
      ...prev,
      [name]: newValue
    }))

    if (touched.has(name)) {
      const error = validateField(name, value)
      setErrors(prev => ({
        ...prev,
        [name]: error
      }))
    }

    if (name === 'password' && touched.has('confirmPassword') && mode === 'signup') {
      const confirmError = validateField('confirmPassword', formData.confirmPassword!)
      setErrors(prev => ({
        ...prev,
        confirmPassword: confirmError
      }))
    }
  }

  const handleBlur = (e: React.FocusEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setTouched(prev => new Set(prev).add(name))
    setFocusedField(null)

    const error = validateField(name, value)
    setErrors(prev => ({
      ...prev,
      [name]: error
    }))
  }

  const handleFocus = (e: React.FocusEvent<HTMLInputElement>) => {
    setFocusedField(e.target.name)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    const newErrors: LoginFormErrors = {}
    const fieldsToValidate = mode === 'signin' 
      ? ['email', 'password']
      : ['name', 'email', 'password', 'confirmPassword']

    fieldsToValidate.forEach(key => {
      const error = validateField(key, formData[key as keyof LoginFormData] as string)
      if (error) newErrors[key as keyof LoginFormErrors] = error
    })

    setErrors(newErrors)

    if (Object.keys(newErrors).length > 0) {
      setShake(true)
      setTimeout(() => setShake(false), 500)
      return
    }

    const submitData = mode === 'signin' 
      ? { email: formData.email, password: formData.password, remember: formData.remember }
      : { name: formData.name, email: formData.email, password: formData.password, remember: formData.remember }

    onSubmit(submitData)
  }

  useEffect(() => {
    setErrors({})
    setTouched(new Set())
  }, [mode])

  useEffect(() => {
    if (error) {
      setShake(true)
      setTimeout(() => setShake(false), 500)
    }
  }, [error])

  return (
    <form 
      onSubmit={handleSubmit} 
      className={`
        w-full space-y-6
        transition-all duration-300
        ${shake ? 'animate-shake' : ''}
      `}
    >
      {mode === 'signup' && (
        <div className="space-y-2">
          <label 
            htmlFor="name" 
            className={`
              block text-sm font-bold uppercase tracking-wide
              transition-colors duration-300
              ${focusedField === 'name' ? 'text-[#00E5B0]' : 'text-[#F5F0E8]'}
              ${errors.name ? 'text-[#FF3B30]' : ''}
              ${!focusedField && !errors.name ? 'opacity-80' : ''}
            `}
          >
            Имя
          </label>
          <div className="relative">
            <input
              type="text"
              id="name"
              name="name"
              value={formData.name}
              onChange={handleChange}
              onFocus={handleFocus}
              onBlur={handleBlur}
              disabled={isLoading}
              placeholder="Иван Иванов"
              className={`
                w-full px-6 py-4 rounded-full
                font-medium text-[#F5F0E8] placeholder-[#F5F0E8]/40
                transition-all duration-300
                outline-none border-2
                ${errors.name 
                  ? 'border-[#FF3B30] bg-[#FF3B30]/10' 
                  : focusedField === 'name'
                    ? 'border-[#00E5B0] bg-white/10'
                    : 'border-transparent bg-white/10'
                }
                ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}
                ${!errors.name && !focusedField ? 'hover:border-[#00E5B0]/30 hover:bg-white/15' : ''}
              `}
            />
            <div className="absolute right-6 top-1/2 transform -translate-y-1/2">
              <div className={`
                w-2 h-2 rounded-full
                transition-all duration-300
                ${formData.name && !errors.name 
                  ? 'bg-[#00E5B0] animate-pulse' 
                  : errors.name 
                    ? 'bg-[#FF3B30]' 
                    : 'bg-[#F5F0E8]/30'
                }
              `} />
            </div>
          </div>
          {errors.name && (
            <p className="text-sm text-[#FF3B30] mt-2 px-6 animate-fadeIn">
              {errors.name}
            </p>
          )}
        </div>
      )}

      <div className="space-y-2">
        <label 
          htmlFor="email" 
          className={`
            block text-sm font-bold uppercase tracking-wide
            transition-colors duration-300
            ${focusedField === 'email' ? 'text-[#00E5B0]' : 'text-[#F5F0E8]'}
            ${errors.email ? 'text-[#FF3B30]' : ''}
            ${!focusedField && !errors.email ? 'opacity-80' : ''}
          `}
        >
          Email
        </label>
        <div className="relative">
          <input
            type="email"
            id="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            onFocus={handleFocus}
            onBlur={handleBlur}
            disabled={isLoading}
            placeholder="example@ponial.ru"
            className={`
              w-full px-6 py-4 rounded-full
              font-medium text-[#F5F0E8] placeholder-[#F5F0E8]/40
              transition-all duration-300
              outline-none border-2
              ${errors.email 
                ? 'border-[#FF3B30] bg-[#FF3B30]/10' 
                : focusedField === 'email'
                  ? 'border-[#00E5B0] bg-white/10'
                  : 'border-transparent bg-white/10'
              }
              ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}
              ${!errors.email && !focusedField ? 'hover:border-[#00E5B0]/30 hover:bg-white/15' : ''}
            `}
          />
          <div className="absolute right-6 top-1/2 transform -translate-y-1/2">
            <div className={`
              w-2 h-2 rounded-full
              transition-all duration-300
              ${formData.email && !errors.email 
                ? 'bg-[#00E5B0] animate-pulse' 
                : errors.email 
                  ? 'bg-[#FF3B30]' 
                  : 'bg-[#F5F0E8]/30'
              }
            `} />
          </div>
        </div>
        {errors.email && (
          <p className="text-sm text-[#FF3B30] mt-2 px-6 animate-fadeIn">
            {errors.email}
          </p>
        )}
      </div>

      <div className="space-y-2">
        <label 
          htmlFor="password" 
          className={`
            block text-sm font-bold uppercase tracking-wide
            transition-colors duration-300
            ${focusedField === 'password' ? 'text-[#00E5B0]' : 'text-[#F5F0E8]'}
            ${errors.password ? 'text-[#FF3B30]' : ''}
            ${!focusedField && !errors.password ? 'opacity-80' : ''}
          `}
        >
          Пароль
        </label>
        <div className="relative">
          <input
            type="password"
            id="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            onFocus={handleFocus}
            onBlur={handleBlur}
            disabled={isLoading}
            placeholder="••••••••"
            className={`
              w-full px-6 py-4 rounded-full
              font-medium text-[#F5F0E8] placeholder-[#F5F0E8]/40
              transition-all duration-300
              outline-none border-2
              ${errors.password 
                ? 'border-[#FF3B30] bg-[#FF3B30]/10' 
                : focusedField === 'password'
                  ? 'border-[#00E5B0] bg-white/10'
                  : 'border-transparent bg-white/10'
              }
              ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}
              ${!errors.password && !focusedField ? 'hover:border-[#00E5B0]/30 hover:bg-white/15' : ''}
            `}
          />
          <div className="absolute right-6 top-1/2 transform -translate-y-1/2 flex gap-1">
            {[1, 2, 3].map(i => (
              <div
                key={i}
                className={`
                  w-1 h-4 rounded-full
                  transition-all duration-300
                  ${formData.password.length >= i * 2 
                    ? formData.password.length > 8 
                      ? 'bg-[#00E5B0]' 
                      : 'bg-yellow-400'
                    : 'bg-[#F5F0E8]/20'
                  }
                `}
              />
            ))}
          </div>
        </div>
        {errors.password && (
          <p className="text-sm text-[#FF3B30] mt-2 px-6 animate-fadeIn">
            {errors.password}
          </p>
        )}
      </div>

      {mode === 'signup' && (
        <div className="space-y-2">
          <label 
            htmlFor="confirmPassword" 
            className={`
              block text-sm font-bold uppercase tracking-wide
              transition-colors duration-300
              ${focusedField === 'confirmPassword' ? 'text-[#00E5B0]' : 'text-[#F5F0E8]'}
              ${errors.confirmPassword ? 'text-[#FF3B30]' : ''}
              ${!focusedField && !errors.confirmPassword ? 'opacity-80' : ''}
            `}
          >
            Подтверждение пароля
          </label>
          <div className="relative">
            <input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              onFocus={handleFocus}
              onBlur={handleBlur}
              disabled={isLoading}
              placeholder="••••••••"
              className={`
                w-full px-6 py-4 rounded-full
                font-medium text-[#F5F0E8] placeholder-[#F5F0E8]/40
                transition-all duration-300
                outline-none border-2
                ${errors.confirmPassword 
                  ? 'border-[#FF3B30] bg-[#FF3B30]/10' 
                  : focusedField === 'confirmPassword'
                    ? 'border-[#00E5B0] bg-white/10'
                    : 'border-transparent bg-white/10'
                }
                ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}
                ${!errors.confirmPassword && !focusedField ? 'hover:border-[#00E5B0]/30 hover:bg-white/15' : ''}
              `}
            />
            <div className="absolute right-6 top-1/2 transform -translate-y-1/2">
              <div className={`
                w-2 h-2 rounded-full
                transition-all duration-300
                ${formData.confirmPassword && !errors.confirmPassword && formData.password === formData.confirmPassword
                  ? 'bg-[#00E5B0] animate-pulse' 
                  : errors.confirmPassword 
                    ? 'bg-[#FF3B30]' 
                    : 'bg-[#F5F0E8]/30'
                }
              `} />
            </div>
          </div>
          {errors.confirmPassword && (
            <p className="text-sm text-[#FF3B30] mt-2 px-6 animate-fadeIn">
              {errors.confirmPassword}
            </p>
          )}
        </div>
      )}

      {error && (
        <div className="
          p-4 rounded-full
          bg-[#FF3B30]/10 border-2 border-[#FF3B30]/30
          text-[#FF3B30] text-sm text-center
          animate-fadeIn
        ">
          {error}
        </div>
      )}

      <button
        type="submit"
        disabled={isLoading}
        className={`
          cursor-pointer w-full py-5 px-8 rounded-full
          font-black text-lg uppercase tracking-wider
          bg-[#00E5B0] text-[#0B1A33]
          transition-all duration-300
          transform hover:scale-105 active:scale-95
          ${isLoading ? 'opacity-50 cursor-not-allowed' : 'hover:shadow-[0_0_30px_rgba(0,229,176,0.3)]'}
          relative overflow-hidden group
        `}
      >
        <span className="relative z-10">
          {isLoading ? 'ПОДОЖДИТЕ...' : mode === 'signin' ? 'ПОНЯЛ, ВОЙТИ' : 'ПОНЯЛ, СОЗДАТЬ'}
        </span>
        
        {isLoading && (
          <div className="absolute inset-0">
            <div className="absolute inset-0 bg-white/30 animate-pulse" />
          </div>
        )}
        <div className="absolute left-0 top-1/2 transform -translate-y-1/2 w-1 h-8 bg-[#0B1A33] opacity-0 group-hover:opacity-20 transition-opacity" />
        <div className="absolute right-0 top-1/2 transform -translate-y-1/2 w-1 h-8 bg-[#0B1A33] opacity-0 group-hover:opacity-20 transition-opacity" />
      </button>
    </form>
  )
}