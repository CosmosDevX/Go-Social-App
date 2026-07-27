import { useState, FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { Post, likePost, dislikePost, deletePost } from '../api/post'
import { Comment, createComment, getCommentsByPostId, deleteComment } from '../api/comment'
import { getErrorMessage, getPostImageUrl } from '../api/client'
import { useAuth } from '../context/AuthContext'

interface PostCardProps {
  post: Post
  /** если не передан — берётся post.creator_name */
  authorUsername?: string
  onLikeChange?: (postId: number, likes: number, isLiked: boolean) => void
  onPostDelete?: (postId: number) => void
}

export function PostCard({ post, authorUsername, onLikeChange, onPostDelete }: PostCardProps) {
  const { username: currentUsername } = useAuth()
  const displayName = authorUsername || post.creator_name
  const isOwnPost = !!currentUsername && currentUsername === displayName
  const imageUrl = getPostImageUrl(post.image_name)

  const [likes, setLikes] = useState(post.likes)
  const [isLiked, setIsLiked] = useState(post.is_liked)
  const [likeLoading, setLikeLoading] = useState(false)
  const [likeError, setLikeError] = useState<string | null>(null)

  const [deletingPost, setDeletingPost] = useState(false)

  // Comments
  const [showComments, setShowComments] = useState(false)
  const [comments, setComments] = useState<Comment[]>([])
  const [commentsCount, setCommentsCount] = useState(post.comments_count ?? 0)
  const [commentsLoading, setCommentsLoading] = useState(false)
  const [commentsError, setCommentsError] = useState<string | null>(null)
  const [commentsLoaded, setCommentsLoaded] = useState(false)

  // New comment form
  const [newComment, setNewComment] = useState('')
  const [commentSubmitting, setCommentSubmitting] = useState(false)
  const [commentError, setCommentError] = useState<string | null>(null)

  const handleLikeToggle = async () => {
    if (likeLoading) return
    setLikeLoading(true)
    setLikeError(null)

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
      setLikeError(getErrorMessage(err))
    } finally {
      setLikeLoading(false)
    }
  }

  const handleDeletePost = async () => {
    if (deletingPost) return
    if (!confirm('Удалить этот пост?')) return

    setDeletingPost(true)
    try {
      await deletePost(post.post_id)
      onPostDelete?.(post.post_id)
    } catch (err) {
      setLikeError(getErrorMessage(err))
      setDeletingPost(false)
    }
  }

  const loadComments = async () => {
    if (commentsLoaded) return
    setCommentsLoading(true)
    setCommentsError(null)
    try {
      const data = await getCommentsByPostId(post.post_id)
      setComments(data)
      setCommentsCount(data.length)
      setCommentsLoaded(true)
    } catch (err) {
      setCommentsError(getErrorMessage(err))
    } finally {
      setCommentsLoading(false)
    }
  }

  const toggleComments = () => {
    const next = !showComments
    setShowComments(next)
    if (next && !commentsLoaded) {
      loadComments()
    }
  }

  const handleSubmitComment = async (e: FormEvent) => {
    e.preventDefault()
    const text = newComment.trim()
    if (text.length < 1 || text.length > 250) {
      setCommentError('Комментарий: от 1 до 250 символов')
      return
    }

    setCommentSubmitting(true)
    setCommentError(null)

    try {
      await createComment(post.post_id, text)
      setNewComment('')
      const data = await getCommentsByPostId(post.post_id)
      setComments(data)
      setCommentsCount(data.length)
      setCommentsLoaded(true)
    } catch (err) {
      setCommentError(getErrorMessage(err))
    } finally {
      setCommentSubmitting(false)
    }
  }

  const handleDeleteComment = async (commentId: number) => {
    if (!confirm('Удалить комментарий?')) return
    try {
      await deleteComment(commentId)
      setComments((prev) => prev.filter((c) => c.comment_id !== commentId))
      setCommentsCount((c) => Math.max(0, c - 1))
    } catch (err) {
      setCommentError(getErrorMessage(err))
    }
  }

  return (
    <article className="glass p-5 hover:border-violet-500/30 transition-all duration-300 group">
      {/* Header */}
      <div className="flex items-start justify-between gap-3 mb-3">
        <div className="min-w-0">
          <h3 className="font-semibold text-lg text-white group-hover:text-violet-200 transition-colors">
            {post.post_name}
          </h3>
          {displayName && (
            <Link
              to={`/profile/${displayName}`}
              className="text-sm text-cyan-400/80 hover:text-cyan-300 transition-colors"
            >
              @{displayName}
            </Link>
          )}
        </div>

        {isOwnPost && (
          <button
            onClick={handleDeletePost}
            disabled={deletingPost}
            className="shrink-0 text-xs text-white/40 hover:text-red-400 px-2 py-1 rounded-lg hover:bg-red-500/10 transition-colors disabled:opacity-50"
            title="Удалить пост"
          >
            {deletingPost ? '...' : '🗑'}
          </button>
        )}
      </div>

      {/* Body */}
      <p className="text-white/70 text-sm leading-relaxed whitespace-pre-wrap mb-4">
        {post.post_description}
      </p>

      {/* Image */}
      {imageUrl && (
        <div className="mb-4 rounded-xl overflow-hidden border border-white/10">
          <img
            src={imageUrl}
            alt={post.post_name}
            className="w-full max-h-96 object-cover"
            loading="lazy"
          />
        </div>
      )}

      {/* Actions */}
      <div className="flex items-center gap-3 flex-wrap">
        <button
          onClick={handleLikeToggle}
          disabled={likeLoading}
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

        <button
          onClick={toggleComments}
          className={`
            flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium
            transition-all duration-200
            ${showComments
              ? 'bg-cyan-500/20 text-cyan-300 border border-cyan-500/30'
              : 'bg-white/5 text-white/60 hover:bg-white/10 hover:text-white border border-transparent'
            }
          `}
        >
          💬
          <span>{commentsCount}</span>
        </button>

        {likeError && (
          <span className="text-xs text-red-400">{likeError}</span>
        )}
      </div>

      {/* Comments section */}
      {showComments && (
        <div className="mt-5 pt-4 border-t border-white/10">
          {commentsLoading ? (
            <div className="flex justify-center py-4">
              <div className="w-6 h-6 border-2 border-cyan-500/30 border-t-cyan-500 rounded-full animate-spin" />
            </div>
          ) : commentsError ? (
            <p className="text-sm text-red-400 py-2">{commentsError}</p>
          ) : (
            <>
              {comments.length === 0 ? (
                <p className="text-sm text-white/40 py-2">Пока нет комментариев</p>
              ) : (
                <ul className="space-y-3 mb-4">
                  {comments.map((c) => {
                    const isOwnComment = !!currentUsername && currentUsername === c.creator_username
                    return (
                      <li
                        key={c.comment_id}
                        className="bg-space-900/50 rounded-xl px-3.5 py-2.5"
                      >
                        <div className="flex items-center justify-between gap-2 mb-1">
                          <Link
                            to={`/profile/${c.creator_username}`}
                            className="text-sm font-medium text-cyan-400/90 hover:text-cyan-300"
                          >
                            @{c.creator_username}
                          </Link>
                          {isOwnComment && (
                            <button
                              onClick={() => handleDeleteComment(c.comment_id)}
                              className="text-xs text-white/30 hover:text-red-400 transition-colors"
                              title="Удалить комментарий"
                            >
                              🗑
                            </button>
                          )}
                        </div>
                        <p className="text-sm text-white/75 whitespace-pre-wrap leading-relaxed">
                          {c.comment_text}
                        </p>
                      </li>
                    )
                  })}
                </ul>
              )}

              {/* Add comment form */}
              <form onSubmit={handleSubmitComment} className="space-y-2">
                <textarea
                  className="input-field min-h-[72px] text-sm resize-y"
                  placeholder="Написать комментарий... (1–250 символов)"
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  maxLength={250}
                  disabled={commentSubmitting}
                />
                <div className="flex items-center justify-between gap-3">
                  <span className="text-xs text-white/40">
                    {newComment.length}/250
                  </span>
                  <button
                    type="submit"
                    className="btn-primary text-sm py-2 px-4"
                    disabled={commentSubmitting || newComment.trim().length < 1}
                  >
                    {commentSubmitting ? 'Отправка...' : 'Отправить'}
                  </button>
                </div>
                {commentError && (
                  <p className="text-xs text-red-400">{commentError}</p>
                )}
              </form>
            </>
          )}
        </div>
      )}
    </article>
  )
}
