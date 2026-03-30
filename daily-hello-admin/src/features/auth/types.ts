export type LoginRequest = {
  email: string
  password: string
}

export type AuthTokens = {
  access_token: string
  expires_in: number
  refresh_token: string
  token_type: string
}

export type LoginResponse = {
  data: AuthTokens
  error_code: string
  error_message: string
  success: boolean
}
