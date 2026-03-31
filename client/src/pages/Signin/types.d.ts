export interface SignInFormData {
  email: string
  password: string
}

export interface SignInFormErrors {
  email?: string
  password?: string
}

export interface SignInFormProps {
  onSubmit: (data: SignInFormData) => void
  isLoading?: boolean
  error?: string | null
}