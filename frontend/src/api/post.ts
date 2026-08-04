import { api } from './client'

export interface Post {
  post_id: number
  post_name: string
  post_description: string
  creator_id: number
  creator_name: string
  likes: number
  is_liked: boolean
  comments_count: number
  image_name: string | null
  created_at: string // TIMESTAMPTZ from backend (ISO string)
}

/** Человекочитаемое время создания поста */
export function formatPostDate(iso: string | null | undefined): string {
  if (!iso) return ''
  const date = new Date(iso)
  if (Number.isNaN(date.getTime())) return ''

  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)

  if (diffSec < 60) return 'только что'
  if (diffMin < 60) return `${diffMin} мин. назад`
  if (diffHour < 24) return `${diffHour} ч. назад`
  if (diffDay < 7) return `${diffDay} дн. назад`

  return date.toLocaleString('ru-RU', {
    day: 'numeric',
    month: 'short',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined,
    hour: '2-digit',
    minute: '2-digit',
  })
}

export async function createPost(
  post_name: string,
  post_description: string,
  imageFile?: File | null
): Promise<{ post_id: number }> {
  const form = new FormData()
  form.append('post_name', post_name)
  form.append('post_description', post_description)
  // post_image — только если пользователь выбрал файл
  if (imageFile) {
    form.append('post_image', imageFile)
  }

  const { data } = await api.post<{ post_id: number }>('/post/create', form)
  return data
}

export async function likePost(postId: number): Promise<{ likes: number }> {
  const { data } = await api.post<{ likes: number }>(`/post/like/${postId}`)
  return data
}

export async function dislikePost(postId: number): Promise<{ likes: number }> {
  const { data } = await api.post<{ likes: number }>(`/post/dislike/${postId}`)
  return data
}

export async function deletePost(postId: number): Promise<{ message: string }> {
  const { data } = await api.delete<{ message: string }>(`/post/${postId}`)
  return data
}

export async function getCurrentUserPosts(): Promise<Post[]> {
  const { data } = await api.get<Post[]>('/post/current_user/all')
  return data
}

export async function getUserPostsByUsername(username: string): Promise<Post[]> {
  const { data } = await api.get<Post[]>(`/post/${username}/all`)
  return data
}

export async function getFeed(): Promise<Post[]> {
  const { data } = await api.get<Post[]>('/post/feed')
  return data
}
