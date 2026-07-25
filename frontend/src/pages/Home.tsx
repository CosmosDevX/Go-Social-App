import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

export function Home() {
  const { isAuthenticated, username } = useAuth()

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4">
      <div className="text-center max-w-lg">
        <div className="text-7xl mb-6 animate-float">🌌</div>
        <h1 className="text-4xl sm:text-5xl font-bold mb-4">
          <span className="bg-gradient-to-r from-violet-400 via-fuchsia-400 to-cyan-400 bg-clip-text text-transparent">
            Space Posts
          </span>
        </h1>
        <p className="text-white/60 text-lg mb-8 leading-relaxed">
          Космическая платформа для постов.
          <br />
          <span className="text-white/40 text-base">
            Общая лента появится позже
          </span>
        </p>

        {isAuthenticated ? (
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Link to={`/profile/${username}`} className="btn-primary">
              Мой профиль
            </Link>
            <Link to="/create" className="btn-ghost border border-white/10">
              Создать пост
            </Link>
          </div>
        ) : (
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Link to="/register" className="btn-primary">
              Начать
            </Link>
            <Link to="/login" className="btn-ghost border border-white/10">
              Войти
            </Link>
          </div>
        )}
      </div>
    </div>
  )
}
