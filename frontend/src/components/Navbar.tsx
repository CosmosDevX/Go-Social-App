import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Navbar() {
  const { isAuthenticated, username, logout } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <header className="sticky top-0 z-50 border-b border-white/10 bg-space-950/80 backdrop-blur-xl">
      <div className="max-w-5xl mx-auto px-4 h-16 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 group">
          <span className="text-2xl">🌌</span>
          <span className="font-semibold text-lg tracking-tight bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent group-hover:from-violet-300 group-hover:to-cyan-300 transition-all">
            КосмоСеть
          </span>
        </Link>

        <nav className="flex items-center gap-2">
          {isAuthenticated ? (
            <>
              <Link to="/create" className="btn-ghost text-sm">
                + Создать
              </Link>
              <Link
                to={`/profile/${username}`}
                className="btn-ghost text-sm flex items-center gap-1.5"
              >
                <span className="w-7 h-7 rounded-full bg-gradient-to-br from-violet-500 to-cyan-400 flex items-center justify-center text-xs font-bold">
                  {username?.[0]?.toUpperCase()}
                </span>
                <span className="hidden sm:inline">{username}</span>
              </Link>
              <button onClick={handleLogout} className="btn-ghost text-sm text-white/50 hover:text-red-400">
                Выйти
              </button>
            </>
          ) : (
            <>
              <Link to="/login" className="btn-ghost text-sm">
                Войти
              </Link>
              <Link to="/register" className="btn-primary text-sm py-2 px-4">
                Регистрация
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  )
}
