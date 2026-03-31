export interface SignUpFormData {
  vin: string
  name: string
  email: string
  password: string
  confirmPassword: string
}

export interface SignUpFormErrors {
  vin?: string
  name?: string
  email?: string
  password?: string
  confirmPassword?: string
}

export interface SignUpFormProps {
  onSubmit: (data: Omit<SignUpFormData, 'confirmPassword'>) => void
  isLoading?: boolean
  error?: string | null
}
