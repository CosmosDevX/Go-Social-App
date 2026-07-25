import { api } from './client'

export async function getCurrentProfile(): Promise<{ username: string }> {
  const { data } = await api.get<{ username: string }>('/user/current/profile')
  return data
}

export async function getUsernameById(userId: number | string): Promise<{ username: string }> {
  const { data } = await api.get<{ username: string }>(`/user/get_username_by_id/${userId}`)
  return data
}
