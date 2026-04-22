"use client"

import { AuthProvider } from "@/lib/auth-context"

export function Providers({ children }: { children: React.ReactNode }) {
  return <AuthProvider>{children}</AuthProvider>
}
