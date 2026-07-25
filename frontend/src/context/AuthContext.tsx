import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from 'react'
import { login as apiLogin, register as apiRegister, refreshToken } from '../api/auth'
import { getCurrentProfile } from '../api/user'
import { getErrorMessage } from '../api/client'

interface AuthState {
  username: string | null
  isAuthenticated: boolean
  isLoading: boolean
}

interface AuthContextValue extends AuthState {
  login: (username: string, password: string) => Promise<void>
  register: (username: string, password: string) => Promise<void>
  logout: () => void
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>({
    username: null,
    isAuthenticated: false,
    isLoading: true,
  })

  const logout = useCallback(() => {
    localStorage.removeItem('access_token')
    setState({
      username: null,
      isAuthenticated: false,
      isLoading: false,
    })
  }, [])

  const refreshUser = useCallback(async () => {
    try {
      const profile = await getCurrentProfile()
      setState({
        username: profile.username,
        isAuthenticated: true,
        isLoading: false,
      })
    } catch {
      // try refresh
      try {
        const { access_token } = await refreshToken()
        localStorage.setItem('access_token', access_token)
        const profile = await getCurrentProfile()
        setState({
          username: profile.username,
          isAuthenticated: true,
          isLoading: false,
        })
      } catch {
        logout()
      }
    }
  }, [logout])

  useEffect(() => {
    const token = localStorage.getItem('access_token')
    if (token) {
      refreshUser()
    } else {
      setState((s) => ({ ...s, isLoading: false }))
    }
  }, [refreshUser])

  const login = async (username: string, password: string) => {
    const { access_token } = await apiLogin(username, password)
    localStorage.setItem('access_token', access_token)
    const profile = await getCurrentProfile()
    setState({
      username: profile.username,
      isAuthenticated: true,
      isLoading: false,
    })
  }

  const register = async (username: string, password: string) => {
    await apiRegister(username, password)
    // after register → auto login
    await login(username, password)
  }

  return (
    <AuthContext.Provider
      value={{
        ...state,
        login,
        register,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}

export { getErrorMessage }
