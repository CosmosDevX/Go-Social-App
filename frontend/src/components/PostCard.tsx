import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Post, likePost, dislikePost } from '../api/post'
import { getErrorMessage } from '../api/client'

interface PostCardProps {
  post: Post
  authorUsername?: string
  onLikeChange?: (postId: number, likes: number, isLiked: boolean) => void
}

export function PostCard({ post, authorUsername, onLikeChange }: PostCardProps) {
  const [likes, setLikes] = useState(post.likes)
  const [isLiked, setIsLiked] = useState(post.is_liked)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleLikeToggle = async () => {
    if (loading) return
    setLoading(true)
    setError(null)

    try {
      if (isLiked) {
        const res = await dislikePost(post.post_id)
        setLikes(res.likes)
        setIsLiked(false)
        onLikeChange?.(post.post_id, res.likes, false)
      } else {
        const res = await likePost(post.post_id)
        setLikes(res.likes)
        setIsLiked(true)
        onLikeChange?.(post.post_id, res.likes, true)
      }
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <article className="glass p-5 hover:border-violet-500/30 transition-all duration-300 group">
      <div className="flex items-start justify-between gap-3 mb-3">
        <div>
          <h3 className="font-semibold text-lg text-white group-hover:text-violet-200 transition-colors">
            {post.post_name}
          </h3>
          {authorUsername && (
            <Link
              to={`/profile/${authorUsername}`}
              className="text-sm text-cyan-400/80 hover:text-cyan-300 transition-colors"
            >
              @{authorUsername}
            </Link>
          )}
        </div>
      </div>

      <p className="text-white/70 text-sm leading-relaxed whitespace-pre-wrap mb-4">
        {post.post_description}
      </p>

      <div className="flex items-center gap-3">
        <button
          onClick={handleLikeToggle}
          disabled={loading}
          className={`
            flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium
            transition-all duration-200
            ${isLiked
              ? 'bg-pink-500/20 text-pink-400 border border-pink-500/30'
              : 'bg-white/5 text-white/60 hover:bg-white/10 hover:text-white border border-transparent'
            }
            disabled:opacity-50
          `}
        >
          <span className={isLiked ? 'scale-110' : ''}>
            {isLiked ? '❤️' : '🤍'}
          </span>
          <span>{likes}</span>
        </button>

        {error && (
          <span className="text-xs text-red-400">{error}</span>
        )}
      </div>
    </article>
  )
}
