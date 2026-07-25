import { FormEvent, useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth, getErrorMessage } from '../context/AuthContext'

export function Register() {
  const { register, isAuthenticated } = useAuth()
  const navigate = useNavigate()

  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/', { replace: true })
    }
  }, [isAuthenticated, navigate])

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)

    if (password !== confirm) {
      setError('Пароли не совпадают')
      return
    }

    if (username.length < 3 || username.length > 60) {
      setError('Имя пользователя: от 3 до 60 символов')
      return
    }

    if (password.length < 10 || password.length > 100) {
      setError('Пароль: от 10 до 100 символов')
      return
    }

    setLoading(true)

    try {
      await register(username.trim(), password)
      navigate('/', { replace: true })
    } catch (err) {
      setError(getErrorMessage(err))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-[calc(100vh-4rem)] flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="text-5xl mb-4 animate-float">✨</div>
          <h1 className="text-3xl font-bold bg-gradient-to-r from-violet-400 to-cyan-400 bg-clip-text text-transparent">
            Регистрация
          </h1>
          <p className="text-white/50 mt-2">Создай свой космический аккаунт</p>
        </div>

        <form onSubmit={handleSubmit} className="glass-strong p-8 space-y-5">
          {error && (
            <div className="px-4 py-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-300 text-sm">
              {error}
            </div>
          )}

          <div>
            <label className="label" htmlFor="username">Имя пользователя</label>
            <input
              id="username"
              type="text"
              className="input-field"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="от 3 до 60 символов"
              required
              minLength={3}
              maxLength={60}
              autoComplete="username"
            />
          </div>

          <div>
            <label className="label" htmlFor="password">Пароль</label>
            <input
              id="password"
              type="password"
              className="input-field"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="от 10 до 100 символов"
              required
              minLength={10}
              maxLength={100}
              autoComplete="new-password"
            />
          </div>

          <div>
            <label className="label" htmlFor="confirm">Подтвердите пароль</label>
            <input
              id="confirm"
              type="password"
              className="input-field"
              value={confirm}
              onChange={(e) => setConfirm(e.target.value)}
              placeholder="ещё раз"
              required
              minLength={10}
              maxLength={100}
              autoComplete="new-password"
            />
          </div>

          <button type="submit" className="btn-primary w-full" disabled={loading}>
            {loading ? 'Создание...' : 'Создать аккаунт'}
          </button>

          <p className="text-center text-sm text-white/50">
            Уже есть аккаунт?{' '}
            <Link to="/login" className="text-cyan-400 hover:text-cyan-300 transition-colors">
              Войти
            </Link>
          </p>
        </form>
      </div>
    </div>
  )
}
