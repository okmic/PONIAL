export const getAuthHeader = () => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export const handlerError = (e: any) => {
  if (e.response?.data?.message) throw new Error(e.response?.data?.message)
  else if (e.response?.data?.error?.message) throw new Error(e.response?.data?.error?.message)
  else throw e
}