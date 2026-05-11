import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { clearTokens } from '@/services/tokenStorage'

const navItems = [
  { to: '/', label: 'Dashboard' },
  { to: '/users', label: 'Quản lý nhân viên' },
  { to: '/branches', label: 'Chi nhánh' },
  { to: '/wifi', label: 'WiFi' },
  { to: '/attendance', label: 'Chấm công' },
  { to: '/devices', label: 'Thiết bị' },
  { to: '/holidays', label: 'Ngày nghỉ' },
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
