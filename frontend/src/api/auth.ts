import { api } from './client'

export interface AuthResponse {
  access_token: string
}

export async function login(username: string, password: string): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>('/auth', { username, password })
  return data
}

export async function register(username: string, password: string): Promise<{ user_id: number }> {
  const { data } = await api.post<{ user_id: number }>('/user/create', { username, password })
  return data
}

export async function refreshToken(): Promise<AuthResponse> {
  const { data } = await api.post<AuthResponse>('/refresh')
  return data
}
