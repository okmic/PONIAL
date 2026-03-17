import { toast } from 'react-hot-toast'

export const useNotifications = () => {
  const success = (message: string) => {
    toast.success(message, {
      style: {
        background: '#3cb550',
        color: 'white',
      },
    })
  }

  const error = (message: string) => {
    toast.error(message, {
      style: {
        background: '#F15D55',
        color: 'white',
      },
    })
  }

  return { success, error }
}