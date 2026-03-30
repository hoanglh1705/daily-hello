import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from './api'
import { saveTokens } from '@/services/tokenStorage'

export default function LoginPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const res = await login({ email, password })
      if (res.success) {
        saveTokens(res.data.access_token, res.data.expires_in, res.data.refresh_token)
        navigate('/', { replace: true })
      } else {
        setError(res.error_message || 'Dang nhap that bai')
      }
    } catch {
      setError('Khong the ket noi den server')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <form className="login-form" onSubmit={handleSubmit}>
        <h1>Daily Hello</h1>
        <p className="login-subtitle">Dang nhap he thong quan tri</p>

        {error && <div className="login-error">{error}</div>}

        <div>
          <label>Email</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="admin@example.com"
            required
          />
        </div>
        <div>
          <label>Mat khau</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="••••••••"
            required
          />
        </div>
        <button type="submit" className="login-btn" disabled={loading}>
          {loading ? 'Dang xu ly...' : 'Dang nhap'}
        </button>
      </form>
    </div>
  )
}
