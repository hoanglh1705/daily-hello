import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { clearTokens } from '@/services/tokenStorage'

const navItems = [
  { to: '/', label: 'Dashboard' },
  { to: '/branches', label: 'Chi nhanh' },
  { to: '/wifi', label: 'WiFi' },
  { to: '/attendance', label: 'Cham cong' },
]

export default function AdminLayout() {
  const navigate = useNavigate()

  const handleLogout = () => {
    clearTokens()
    navigate('/login', { replace: true })
  }

  return (
    <div className="admin-layout">
      <aside className="sidebar">
        <div className="sidebar-header">
          <h2>Daily Hello</h2>
        </div>
        <nav>
          {navItems.map((item) => (
            <NavLink key={item.to} to={item.to} className="nav-link">
              {item.label}
            </NavLink>
          ))}
        </nav>
        <div className="sidebar-footer">
          <button className="logout-btn" onClick={handleLogout}>
            Dang xuat
          </button>
        </div>
      </aside>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  )
}
