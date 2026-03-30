const ACCESS_TOKEN_KEY = 'access_token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const TOKEN_EXPIRES_AT_KEY = 'token_expires_at'

export function saveTokens(accessToken: string, expiresIn: number, refreshToken: string) {
  const expiresAt = Date.now() + expiresIn * 1000
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
  localStorage.setItem(TOKEN_EXPIRES_AT_KEY, String(expiresAt))
}

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY)
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}

export function isTokenExpired(): boolean {
  const expiresAt = localStorage.getItem(TOKEN_EXPIRES_AT_KEY)
  if (!expiresAt) return true
  return Date.now() >= Number(expiresAt)
}

export function clearTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY)
  localStorage.removeItem(REFRESH_TOKEN_KEY)
  localStorage.removeItem(TOKEN_EXPIRES_AT_KEY)
}

export function isAuthenticated(): boolean {
  return !!getAccessToken()
}
