import { useEffect, useState } from "react"
import { useDispatch, useSelector } from "react-redux"
import apiAuthService from "../libs/api/api.auth.service"
import type { RootState } from "../store/store"
import {  } from "../store/slices/taxonomy.slice"
import { authStatus } from "../store/slices/auth.slice"
import type { IUser } from "../pkg/types/user"

interface UseInitializeAppReturn {
  isCheckingAuth: boolean
  initialized: boolean
  user: IUser | null
  error: Error | null
}

export const useInitializeApp = (): UseInitializeAppReturn => {
  const dispatch = useDispatch()
  const [state, setState] = useState<UseInitializeAppReturn>({
    isCheckingAuth: true,
    initialized: false,
    user: null,
    error: null,
  })
  const auth = useSelector((s: RootState) => s.auth)

  useEffect(() => {
    const initializeApp = async () => {
      try {
        const [me] = await Promise.all([
          apiAuthService.geMyInfo(),
        ])
        if (me) {
          dispatch(authStatus({ status: "auth", user: me }))
          setState({
            isCheckingAuth: false,
            initialized: true,
            user: me,
            error: null,
          })
        } else {
          dispatch(authStatus({ status: "notAuth", user: null }))
          setState({
            isCheckingAuth: false,
            initialized: true,
            user: null,
            error: null,
          })
        }
      } catch (error) {
        dispatch(authStatus({ status: "notAuth", user: null }))
        setState({
          isCheckingAuth: false,
          initialized: true,
          user: null,
          error: error as Error,
        })
      }
    }
    initializeApp()
  }, [dispatch, auth.authStatus])

  return state
}