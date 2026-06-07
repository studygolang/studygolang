"use client"

import { createContext, useContext, useEffect, useState, useCallback } from "react"
import type { Me } from "@/lib/types"

interface AuthState {
  user: Me | null
  isLoading: boolean
  isLoggedIn: boolean
}

interface AuthContextValue extends AuthState {
  login: (user: Me) => void
  logout: () => Promise<void>
  refresh: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue>({
  user: null,
  isLoading: true,
  isLoggedIn: false,
  login: () => {},
  logout: async () => {},
  refresh: async () => {},
})

export function useAuth() {
  return useContext(AuthContext)
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [state, setState] = useState<AuthState>({
    user: null,
    isLoading: true,
    isLoggedIn: false,
  })

  const refresh = useCallback(async () => {
    try {
      const res = await fetch("/api/v1/user/me", { credentials: "include" })
      const json = await res.json()
      if (json.code === 0 && json.data?.user) {
        setState({ user: json.data.user, isLoading: false, isLoggedIn: true })
      } else {
        setState({ user: null, isLoading: false, isLoggedIn: false })
      }
    } catch {
      setState({ user: null, isLoading: false, isLoggedIn: false })
    }
  }, [])

  useEffect(() => {
    refresh()
  }, [refresh])

  const login = useCallback((user: Me) => {
    setState({ user, isLoading: false, isLoggedIn: true })
  }, [])

  const logout = useCallback(async () => {
    try {
      await fetch("/api/v1/user/logout", { method: "POST", credentials: "include" })
    } catch {
      // ignore
    }
    setState({ user: null, isLoading: false, isLoggedIn: false })
  }, [])

  return (
    <AuthContext.Provider value={{ ...state, login, logout, refresh }}>
      {children}
    </AuthContext.Provider>
  )
}
