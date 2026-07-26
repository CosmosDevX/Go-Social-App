import { api } from './client'

export interface Comment {
  comment_id: number
  post_id: number
  comment_text: string
  creator_username: string
  creator_id: number
}

export async function createComment(postId: number, comment_text: string): Promise<{ comment_id: number }> {
  const { data } = await api.post<{ comment_id: number }>(`/comment/create/${postId}`, {
    comment_text,
  })
  return data
}

export async function getCommentsByPostId(postId: number): Promise<Comment[]> {
  const { data } = await api.get<Comment[]>(`/comment/all/${postId}`)
  return data
}

export async function deleteComment(commentId: number): Promise<{ message: string }> {
  const { data } = await api.delete<{ message: string }>(`/comment/${commentId}`)
  return data
}
