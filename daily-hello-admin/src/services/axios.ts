import axios from 'axios'
import {
  getAccessToken,
  getRefreshToken,
  isTokenExpired,
  saveTokens,
  clearTokens,
} from './tokenStorage'
import { refreshToken as refreshTokenApi } from '@/features/auth/api'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8282/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

let isRefreshing = false
let pendingRequests: Array<(token: string) => void> = []

function onTokenRefreshed(token: string) {
  pendingRequests.forEach((cb) => cb(token))
  pendingRequests = []
}

function redirectToLogin() {
  clearTokens()
  window.location.href = '/login'
}

instance.interceptors.request.use(async (config) => {
  let token = getAccessToken()

  if (token && isTokenExpired()) {
    const rt = getRefreshToken()
    if (!rt) {
      redirectToLogin()
      return Promise.reject(new Error('Session expired'))
    }

    if (!isRefreshing) {
      isRefreshing = true
      try {
        const res = await refreshTokenApi(rt)
        if (res.success) {
          saveTokens(res.data.access_token, res.data.expires_in, res.data.refresh_token)
          token = res.data.access_token
          onTokenRefreshed(token)
        } else {
          redirectToLogin()
          return Promise.reject(new Error('Session expired'))
        }
      } catch {
        redirectToLogin()
        return Promise.reject(new Error('Session expired'))
      } finally {
        isRefreshing = false
      }
    } else {
      token = await new Promise<string>((resolve) => {
        pendingRequests.push(resolve)
      })
    }
  }

  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    if (error.response?.status === 401) {
      const rt = getRefreshToken()
      if (!rt) {
        redirectToLogin()
        return Promise.reject(error)
      }

      if (!isRefreshing) {
        isRefreshing = true
        try {
          const res = await refreshTokenApi(rt)
          if (res.success) {
            saveTokens(res.data.access_token, res.data.expires_in, res.data.refresh_token)
            onTokenRefreshed(res.data.access_token)
            error.config.headers.Authorization = `Bearer ${res.data.access_token}`
            return instance.request(error.config)
          }
        } catch {
          // refresh failed
        } finally {
          isRefreshing = false
        }
      } else {
        const token = await new Promise<string>((resolve) => {
          pendingRequests.push(resolve)
        })
        error.config.headers.Authorization = `Bearer ${token}`
        return instance.request(error.config)
      }

      redirectToLogin()
    }
    return Promise.reject(error)
  },
)

export default instance
