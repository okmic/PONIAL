export interface LoginFormData {
  email: string
  password: string
  remember: boolean
}

export interface LoginFormErrors {
  email?: string
  password?: string
}

export interface LoginPageProps {
  onLogin?: (data: LoginFormData) => void
  isLoading?: boolean
  error?: string | null
}

export interface LoginFormProps extends LoginPageProps {
  onSubmit: (data: LoginFormData) => void
}