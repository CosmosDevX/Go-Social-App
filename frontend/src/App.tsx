import { Routes, Route } from 'react-router-dom'
import { Navbar } from './components/Navbar'
import { Stars } from './components/Stars'
import { ProtectedRoute } from './components/ProtectedRoute'
import { Home } from './pages/Home'
import { Login } from './pages/Login'
import { Register } from './pages/Register'
import { CreatePost } from './pages/CreatePost'
import { Profile } from './pages/Profile'

export default function App() {
  return (
    <div className="relative min-h-screen">
      <Stars />
      <div className="relative z-10">
        <Navbar />
        <main>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route
              path="/create"
              element={
                <ProtectedRoute>
                  <CreatePost />
                </ProtectedRoute>
              }
            />
            <Route
              path="/profile/:username"
              element={
                <ProtectedRoute>
                  <Profile />
                </ProtectedRoute>
              }
            />
            {/* fallback */}
            <Route path="*" element={<Home />} />
          </Routes>
        </main>
      </div>
    </div>
  )
}
