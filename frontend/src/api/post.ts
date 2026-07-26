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
}

export async function createPost(post_name: string, post_description: string): Promise<{ post_id: number }> {
  const { data } = await api.post<{ post_id: number }>('/post/create', {
    post_name,
    post_description,
  })
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
