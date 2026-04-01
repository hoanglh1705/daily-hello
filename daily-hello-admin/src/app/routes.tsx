import { createBrowserRouter, Navigate } from 'react-router-dom'
import AdminLayout from '@/layouts/AdminLayout'
import LoginPage from '@/features/auth/LoginPage'
import DashboardPage from '@/features/dashboard/DashboardPage'
import BranchPage from '@/features/branch/BranchPage'
import WifiPage from '@/features/wifi/WifiPage'
import AttendancePage from '@/features/attendance/AttendancePage'
import UserPage from '@/features/user/UserPage'
import DevicePage from '@/features/device/DevicePage'
import { isAuthenticated } from '@/services/tokenStorage'

function RequireAuth({ children }: { children: React.ReactNode }) {
  if (!isAuthenticated()) {
    return <Navigate to="/login" replace />
  }
  return <>{children}</>
}

function GuestOnly({ children }: { children: React.ReactNode }) {
  if (isAuthenticated()) {
    return <Navigate to="/" replace />
  }
  return <>{children}</>
}

export const router = createBrowserRouter([
  {
    path: '/login',
    element: (
      <GuestOnly>
        <LoginPage />
      </GuestOnly>
    ),
  },
  {
    element: (
      <RequireAuth>
        <AdminLayout />
      </RequireAuth>
    ),
    children: [
      { path: '/', element: <DashboardPage /> },
      { path: '/users', element: <UserPage /> },
      { path: '/branches', element: <BranchPage /> },
      { path: '/wifi', element: <WifiPage /> },
      { path: '/attendance', element: <AttendancePage /> },
      { path: '/devices', element: <DevicePage /> },
    ],
  },
])
