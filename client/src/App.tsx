import { useEffect } from "react"
import { useSelector } from "react-redux"
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import type { RootState } from "./store/store"
import { ThemeProvider } from "./components/providers/ThemeProvider"
import { DashboardLayout } from "./components/Layout/DashboardLayout"
import { useInitializeApp } from "./hooks/useInitializeApp"
import { Toaster } from "react-hot-toast"
import LoadingSpinner from "./components/UI/LoadingSpinner"
import LoginPage from "./pages/Login"

function MainApp() {
  const auth = useSelector((s: RootState) => s.auth)
  const {
    isCheckingAuth,
    initialized,
    error,
  } = useInitializeApp()

  useEffect(() => {
    if (error && initialized) {
      console.error("Initialization error:", error)
    }
  }, [error, initialized])

  if (isCheckingAuth) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-[#fff]">
        <LoadingSpinner />
      </div>
    )
  }

  if (auth.authStatus === "notAuth" || !auth.user) {
    return <LoginPage />
  }

  return (
    <ThemeProvider>
      <Toaster
        position="top-right"
        toastOptions={{
          duration: 4000,
        }}
      />
      <div className="min-h-screen transition-all duration-300">
        <Routes>
          {auth.authStatus === "auth" && (
            <Route element={<DashboardLayout />}>
              {(auth.user!.role === "admin") && (
                <>
                  <Route path="/" element={<>hello</>} />
                </>
              )}
              <Route path="*" element={<h1>NOT FOUND</h1>} />
            </Route>
          )}
        </Routes>
      </div>
    </ThemeProvider>
  )
}

export default function App() {
  return (
    <Router>
      <MainApp />
    </Router>
  )
}