import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { getFeed, Post } from '../api/post'
import { getErrorMessage, isNotFoundError } from '../api/client'
import { PostCard } from '../components/PostCard'
import { Spinner } from '../components/Spinner'

export function Home() {
  const { isAuthenticated, username } = useAuth()
  const [posts, setPosts] = useState<Post[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!isAuthenticated) return

    let cancelled = false
    setLoading(true)
    setError(null)

    getFeed()
      .then((data) => {
        if (!cancelled) setPosts(data)
      })
      .catch((err) => {
        if (!cancelled) {
          if (isNotFoundError(err)) {
            setPosts([])
            setError(null)
          } else {
            setError(getErrorMessage(err))
          }
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false)
      })

    return () => {
      cancelled = true
    }
  }, [isAuthenticated])

  const handleLikeChange = (postId: number, likes: number, isLiked: boolean) => {
    setPosts((prev) =>
      prev.map((p) =>
        p.post_id === postId ? { ...p, likes, is_liked: isLiked } : p
      )
    )
  }

  const handlePostDelete = (postId: number) => {
    setPosts((prev) => prev.filter((p) => p.post_id !== postId))
  }

  // Не авторизован — приветственный экран
  if (!isAuthenticated) {
    return (
      <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4">
        <div className="text-center max-w-lg">
          <div className="text-7xl mb-6 animate-float">🌌</div>
          <h1 className="text-4xl sm:text-5xl font-bold mb-4">
            <span className="bg-gradient-to-r from-violet-400 via-fuchsia-400 to-cyan-400 bg-clip-text text-transparent">
              КосмоСеть
            </span>
          </h1>
          <p className="text-white/60 text-lg mb-8 leading-relaxed">
            Космическая социальная сеть.
          </p>
          <div className="flex flex-col sm:flex-row gap-3 justify-center">
            <Link to="/register" className="btn-primary">
              Начать
            </Link>
            <Link to="/login" className="btn-ghost border border-white/10">
              Войти
            </Link>
          </div>
        </div>
      </div>
    )
  }

  // Авторизован — лента
  return (
    <div className="max-w-3xl mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent">
            Лента
          </h1>
          <p className="text-sm text-white/40 mt-0.5">Случайные посты из космоса</p>
        </div>
        <div className="flex gap-2">
          <Link to={`/profile/${username}`} className="btn-ghost text-sm border border-white/10">
            Мой профиль
          </Link>
          <Link to="/create" className="btn-primary text-sm py-2">
            + Пост
          </Link>
        </div>
      </div>

      {loading ? (
        <div className="py-20">
          <Spinner />
        </div>
      ) : error ? (
        <div className="glass p-8 text-center">
          <p className="text-red-300 mb-2">{error}</p>
          <button
            onClick={() => window.location.reload()}
            className="btn-ghost text-sm mt-2"
          >
            Обновить
          </button>
        </div>
      ) : posts.length === 0 ? (
        <div className="glass p-12 text-center">
          <div className="text-4xl mb-3 opacity-50">🪐</div>
          <p className="text-white/50">Пока нет постов в ленте</p>
          <Link to="/create" className="btn-primary inline-block mt-5 text-sm">
            Создать первый
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => (
            <PostCard
              key={post.post_id}
              post={post}
              onLikeChange={handleLikeChange}
              onPostDelete={handlePostDelete}
            />
          ))}
        </div>
      )}
    </div>
  )
}
