import axios from 'axios'
import type { LoginRequest, LoginResponse } from './types'

const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:8282/api'

export const login = (data: LoginRequest): Promise<LoginResponse> => {
  return axios.post(`${baseURL}/v1/auth/login`, data).then((res) => res.data)
}

export const refreshToken = (refresh_token: string): Promise<LoginResponse> => {
  return axios
    .post(`${baseURL}/v1/auth/refresh-token`, { refresh_token })
    .then((res) => res.data)
}
