import axios, { AxiosError, InternalAxiosRequestConfig } from 'axios'

export const API_BASE = 'http://localhost:8080'

export const api = axios.create({
  baseURL: API_BASE,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
})

// Attach access token; for FormData let browser set Content-Type with boundary
api.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = localStorage.getItem('access_token')
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`
  }
  if (config.data instanceof FormData && config.headers) {
    delete config.headers['Content-Type']
  }
  return config
})

let isRefreshing = false
let failedQueue: Array<{
  resolve: (token: string) => void
  reject: (err: unknown) => void
}> = []

const processQueue = (error: unknown, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error)
    } else if (token) {
      prom.resolve(token)
    }
  })
  failedQueue = []
}

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        })
          .then((token) => {
            if (originalRequest.headers) {
              originalRequest.headers.Authorization = `Bearer ${token}`
            }
            return api(originalRequest)
          })
          .catch((err) => Promise.reject(err))
      }

      originalRequest._retry = true
      isRefreshing = true

      try {
        const { data } = await axios.post(
          `${API_BASE}/refresh`,
          {},
          { withCredentials: true }
        )
        const newToken = data.access_token
        localStorage.setItem('access_token', newToken)
        processQueue(null, newToken)

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${newToken}`
        }
        return api(originalRequest)
      } catch (refreshError) {
        processQueue(refreshError, null)
        localStorage.removeItem('access_token')
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data

    // JSON: { "code": "...", "message": "..." }
    if (data && typeof data === 'object' && !Array.isArray(data)) {
      const obj = data as Record<string, unknown>
      if (typeof obj.message === 'string' && obj.message.trim()) {
        return obj.message
      }
      if (typeof obj.error === 'string' && obj.error.trim()) {
        return obj.error
      }
    }

    // plain text (на всякий случай)
    if (typeof data === 'string' && data.trim()) {
      return data
    }

    if (error.response?.status === 401) return 'Unauthorized'
    if (error.response?.status === 400) return 'Bad Request'
    if (error.response?.status === 404) return 'Not Found'
    if (error.response?.status === 504) return 'Gateway Timeout'
    return error.message || 'Network error'
  }
  if (error instanceof Error) return error.message
  return 'Something went wrong'
}

/** URL картинки поста (или null, если нет) */
export function getPostImageUrl(imageName: string | null | undefined): string | null {
  if (!imageName || !imageName.trim()) return null
  return `${API_BASE}/uploads/${imageName}`
}
