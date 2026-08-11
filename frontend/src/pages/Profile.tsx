import { useEffect, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { getUserPostsByUsername, getCurrentUserPosts, Post } from '../api/post'
import { getErrorMessage, isNotFoundError } from '../api/client'
import { useAuth } from '../context/AuthContext'
import { PostCard } from '../components/PostCard'
import { Spinner } from '../components/Spinner'

export function Profile() {
  const { username: paramUsername } = useParams<{ username: string }>()
  const { username: currentUsername, isAuthenticated } = useAuth()

  const [posts, setPosts] = useState<Post[]>([])
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const isOwnProfile = isAuthenticated && paramUsername === currentUsername

  // при смене пользователя — сброс на первую страницу
  useEffect(() => {
    setPage(1)
  }, [paramUsername])

  useEffect(() => {
    if (!paramUsername) return

    let cancelled = false
    setLoading(true)
    setError(null)

    const load = async () => {
      try {
        const data = isOwnProfile
          ? await getCurrentUserPosts(page)
          : await getUserPostsByUsername(paramUsername, page)

        if (!cancelled) {
          setPosts(data.posts ?? [])
          setPageSize(data.page_size || 10)
          // если бэк вернул page — синхронизируем
          if (data.page && data.page !== page) {
            setPage(data.page)
          }
        }
      } catch (err) {
        if (!cancelled) {
          if (isNotFoundError(err)) {
            setPosts([])
            setError(null)
          } else {
            setError(getErrorMessage(err))
            setPosts([])
          }
        }
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    load()
    return () => {
      cancelled = true
    }
  }, [paramUsername, isOwnProfile, page])

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

  // has_more: на странице пришло полный page_size → скорее всего есть ещё
  const hasNext = posts.length >= pageSize
  const hasPrev = page > 1

  if (!paramUsername) {
    return (
      <div className="text-center py-20 text-white/50">
        Пользователь не указан
      </div>
    )
  }

  return (
    <div className="max-w-3xl mx-auto px-4 py-10">
      {/* Header */}
      <div className="glass-strong p-6 sm:p-8 mb-8 flex flex-col sm:flex-row sm:items-center gap-5">
        <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-violet-500 to-cyan-400 flex items-center justify-center text-2xl font-bold shadow-lg shadow-violet-500/30">
          {paramUsername[0]?.toUpperCase()}
        </div>
        <div className="flex-1">
          <h1 className="text-2xl font-bold">@{paramUsername}</h1>
          <p className="text-white/50 text-sm mt-0.5">
            {isOwnProfile ? 'Твой профиль' : 'Профиль пользователя'}
          </p>
        </div>
        {isOwnProfile && (
          <Link to="/create" className="btn-primary text-sm self-start sm:self-center">
            + Новый пост
          </Link>
        )}
      </div>

      {/* Posts */}
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold text-white/80">
          Посты
        </h2>
        {!loading && posts.length > 0 && (
          <span className="text-sm text-white/40">Страница {page}</span>
        )}
      </div>

      {loading ? (
        <div className="py-16">
          <Spinner />
        </div>
      ) : error ? (
        <div className="glass p-8 text-center">
          <p className="text-red-300 mb-2">{error}</p>
          <p className="text-white/40 text-sm">
            Возможно, пользователь не существует или произошла ошибка
          </p>
        </div>
      ) : posts.length === 0 ? (
        <div className="glass p-12 text-center">
          <div className="text-4xl mb-3 opacity-50">🪐</div>
          <p className="text-white/50">
            {page > 1
              ? 'На этой странице постов нет'
              : isOwnProfile
                ? 'У тебя пока нет постов. Создай первый!'
                : 'У этого пользователя пока нет постов'}
          </p>
          {page > 1 ? (
            <button
              type="button"
              className="btn-ghost border border-white/10 mt-5 text-sm"
              onClick={() => setPage((p) => Math.max(1, p - 1))}
            >
              ← Предыдущая страница
            </button>
          ) : (
            isOwnProfile && (
              <Link to="/create" className="btn-primary inline-block mt-5 text-sm">
                Создать пост
              </Link>
            )
          )}
        </div>
      ) : (
        <>
          <div className="space-y-4">
            {posts.map((post) => (
              <PostCard
                key={post.post_id}
                post={post}
                authorUsername={paramUsername}
                onLikeChange={handleLikeChange}
                onPostDelete={handlePostDelete}
              />
            ))}
          </div>

          {/* Pagination */}
          <div className="mt-8 flex items-center justify-center gap-3">
            <button
              type="button"
              className="btn-ghost border border-white/10 text-sm px-4 py-2 disabled:opacity-30 disabled:cursor-not-allowed"
              disabled={!hasPrev || loading}
              onClick={() => setPage((p) => Math.max(1, p - 1))}
            >
              ← Предыдущая
            </button>
            <span className="text-sm text-white/50 min-w-[5rem] text-center">
              {page}
            </span>
            <button
              type="button"
              className="btn-ghost border border-white/10 text-sm px-4 py-2 disabled:opacity-30 disabled:cursor-not-allowed"
              disabled={!hasNext || loading}
              onClick={() => setPage((p) => p + 1)}
            >
              Следующая →
            </button>
          </div>
        </>
      )}
    </div>
  )
}
