import { FormEvent, useState, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { createPost } from '../api/post'
import { getErrorMessage } from '../api/client'
import { useAuth } from '../context/AuthContext'

export function CreatePost() {
  const { username } = useAuth()
  const navigate = useNavigate()
  const fileInputRef = useRef<HTMLInputElement>(null)

  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [imageFile, setImageFile] = useState<File | null>(null)
  const [previewUrl, setPreviewUrl] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl)
    }
    if (file) {
      setImageFile(file)
      setPreviewUrl(URL.createObjectURL(file))
    } else {
      setImageFile(null)
      setPreviewUrl(null)
    }
  }

  const clearImage = () => {
    if (previewUrl) URL.revokeObjectURL(previewUrl)
    setImageFile(null)
    setPreviewUrl(null)
    if (fileInputRef.current) fileInputRef.current.value = ''
  }

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)

    if (name.length < 5 || name.length > 100) {
      setError('Название: от 5 до 100 символов')
      return
    }
    if (description.length < 1 || description.length > 900) {
      setError('Описание: от 1 до 900 символов')
      return
    }

    setLoading(true)
    try {
      await createPost(name.trim(), description.trim(), imageFile)
      navigate(`/profile/${username}`)
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl mx-auto px-4 py-10">
      <div className="mb-8">
        <h1 className="text-3xl font-bold bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent">
          Новый пост
        </h1>
        <p className="text-white/50 mt-1">Поделись чем-то интересным из космоса</p>
      </div>

      <form onSubmit={handleSubmit} className="glass-strong p-6 sm:p-8 space-y-5">
        {error && (
          <div className="px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-300 text-sm">
            {error}
          </div>
        )}

        <div>
          <label className="label" htmlFor="name">
            Название <span className="text-white/40">({name.length}/100)</span>
          </label>
          <input
            id="name"
            type="text"
            className="input-field"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Минимум 5 символов"
            required
            minLength={5}
            maxLength={100}
          />
        </div>

        <div>
          <label className="label" htmlFor="description">
            Описание <span className="text-white/40">({description.length}/900)</span>
          </label>
          <textarea
            id="description"
            className="input-field min-h-[140px] resize-y"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="О чём твой пост?"
            required
            minLength={1}
            maxLength={900}
          />
        </div>

        {/* Image upload */}
        <div>
          <label className="label">Картинка <span className="text-white/40">(необязательно)</span></label>

          {previewUrl ? (
            <div className="relative rounded-xl overflow-hidden border border-white/10">
              <img
                src={previewUrl}
                alt="Превью"
                className="w-full max-h-64 object-cover"
              />
              <button
                type="button"
                onClick={clearImage}
                className="absolute top-2 right-2 px-2.5 py-1 rounded-lg bg-black/60 text-white/90 text-xs hover:bg-red-500/80 transition-colors"
              >
                Убрать
              </button>
            </div>
          ) : (
            <label className="flex flex-col items-center justify-center gap-2 w-full h-32 rounded-xl border border-dashed border-white/15 bg-space-900/50 hover:border-violet-500/40 hover:bg-space-900/80 cursor-pointer transition-all">
              <span className="text-2xl opacity-60">🖼</span>
              <span className="text-sm text-white/50">Нажми, чтобы выбрать изображение</span>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/*"
                className="hidden"
                onChange={handleImageChange}
              />
            </label>
          )}
        </div>

        <div className="flex gap-3 pt-2">
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Публикация...' : 'Опубликовать'}
          </button>
          <button
            type="button"
            className="btn-ghost"
            onClick={() => navigate(-1)}
            disabled={loading}
          >
            Отмена
          </button>
        </div>
      </form>
    </div>
  )
}
