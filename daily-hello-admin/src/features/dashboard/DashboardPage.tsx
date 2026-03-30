import { useEffect, useState } from 'react'
import { getDashboard } from './api'
import { getBranches } from '@/features/branch/api'
import type { DashboardData, AttendanceTrend } from './types'
import type { Branch } from '@/features/branch/types'

const MOCK_TRENDS: AttendanceTrend[] = [
  { day: '2024-10-21', label: 'Mon', check_in_count: 65, total: 100 },
  { day: '2024-10-22', label: 'Tue', check_in_count: 72, total: 100 },
  { day: '2024-10-23', label: 'Wed', check_in_count: 90, total: 100 },
  { day: '2024-10-24', label: 'Thu', check_in_count: 55, total: 100 },
  { day: '2024-10-25', label: 'Fri', check_in_count: 40, total: 100 },
  { day: '2024-10-26', label: 'Sat', check_in_count: 48, total: 100 },
  { day: '2024-10-27', label: 'Sun', check_in_count: 45, total: 100 },
]

const MOCK_DATA: DashboardData = {
  stats: {
    total_employees: 1284,
    on_time_percentage: 94,
    on_time_change: 2.4,
    late_arrivals: 12,
    late_arrivals_change: -3,
  },
  trends: MOCK_TRENDS,
  recent_activities: [
    { id: 1, user_name: 'Nguyen Van A', action: 'Check-in', time: '08:01' },
    { id: 2, user_name: 'Tran Thi B', action: 'Check-in', time: '08:15' },
    { id: 3, user_name: 'Le Van C', action: 'Late arrival', time: '09:30' },
  ],
}

function formatDate(date: Date): string {
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export default function DashboardPage() {
  const [data, setData] = useState<DashboardData>(MOCK_DATA)
  const [branches, setBranches] = useState<Branch[]>([])
  const [selectedBranch, setSelectedBranch] = useState<number | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getBranches({ page: 1, limit: 100 })
      .then((res) => setBranches(res.data.items))
      .catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    getDashboard(selectedBranch ? { branch_id: selectedBranch } : undefined)
      .then((res) => {
        if (res.data) setData(res.data)
      })
      .catch(() => {
        // Use mock data as fallback
      })
      .finally(() => setLoading(false))
  }, [selectedBranch])

  const { stats, trends } = data
  const maxCount = Math.max(...trends.map((t) => t.check_in_count), 1)
  const today = new Date()

  return (
    <div className="dashboard">
      <div className="dashboard-header">
        <div>
          <h1 className="dashboard-title">Daily Overview</h1>
          <p className="dashboard-date">Today, {formatDate(today)}</p>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="dashboard-stats-main">
        <div className="stats-total-card">
          <div className="stats-total-icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </div>
          <div>
            <span className="stats-total-label">Total Employee</span>
            <span className="stats-total-value">
              {loading ? '—' : stats.total_employees.toLocaleString()}
            </span>
          </div>
        </div>
      </div>

      <div className="dashboard-stats-row">
        <div className="stats-card stats-card-success">
          <span className="stats-card-label">On Time</span>
          <div className="stats-card-value-row">
            <span className="stats-card-value">{loading ? '—' : `${stats.on_time_percentage}%`}</span>
            {stats.on_time_change !== 0 && (
              <span className={`stats-change ${stats.on_time_change > 0 ? 'positive' : 'negative'}`}>
                {stats.on_time_change > 0 ? '+' : ''}{stats.on_time_change}%
              </span>
            )}
          </div>
        </div>
        <div className="stats-card stats-card-danger">
          <span className="stats-card-label">Late Arrival</span>
          <div className="stats-card-value-row">
            <span className="stats-card-value">{loading ? '—' : stats.late_arrivals}</span>
            {stats.late_arrivals_change !== 0 && (
              <span className={`stats-change ${stats.late_arrivals_change < 0 ? 'positive' : 'negative'}`}>
                {stats.late_arrivals_change > 0 ? '+' : ''}{stats.late_arrivals_change}
              </span>
            )}
          </div>
        </div>
      </div>

      {/* Branch Filter */}
      <div className="dashboard-filters">
        <div className="toolbar-filters">
          <button
            className={`toolbar-chip ${selectedBranch === null ? 'active' : ''}`}
            onClick={() => setSelectedBranch(null)}
          >
            All Branches
          </button>
          {branches.map((b) => (
            <button
              key={b.id}
              className={`toolbar-chip ${selectedBranch === b.id ? 'active' : ''}`}
              onClick={() => setSelectedBranch(b.id)}
            >
              {b.name}
            </button>
          ))}
        </div>
      </div>

      {/* Attendance Trends Chart */}
      <div className="dashboard-card">
        <div className="dashboard-card-header">
          <h2>Attendance Trends</h2>
          <span className="dashboard-card-subtitle">Last 7 Days</span>
        </div>
        <div className="chart-container">
          {trends.map((t, i) => (
            <div className="chart-bar-group" key={t.day}>
              <div className="chart-bar-wrapper">
                <div
                  className={`chart-bar ${i === trends.findIndex((x) => x.check_in_count === maxCount) ? 'chart-bar-highlight' : ''}`}
                  style={{ height: `${(t.check_in_count / maxCount) * 100}%` }}
                />
              </div>
              <span className="chart-label">{t.label}</span>
            </div>
          ))}
        </div>
      </div>

      {/* Bottom Cards */}
      <div className="dashboard-bottom-row">
        <div className="dashboard-card">
          <div className="dashboard-card-header">
            <div className="dashboard-card-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
                <line x1="16" y1="2" x2="16" y2="6" />
                <line x1="8" y1="2" x2="8" y2="6" />
                <line x1="3" y1="10" x2="21" y2="10" />
              </svg>
            </div>
            <h2>Recent Activity</h2>
          </div>
          <div className="activity-list">
            {data.recent_activities.map((a) => (
              <div className="activity-item" key={a.id}>
                <div className="activity-avatar">
                  {a.user_name.charAt(0)}
                </div>
                <div className="activity-info">
                  <span className="activity-name">{a.user_name}</span>
                  <span className="activity-action">{a.action}</span>
                </div>
                <span className="activity-time">{a.time}</span>
              </div>
            ))}
          </div>
        </div>

        <div className="dashboard-card">
          <div className="dashboard-card-header">
            <div className="dashboard-card-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M5 12.55a11 11 0 0 1 14.08 0" />
                <path d="M1.42 9a16 16 0 0 1 21.16 0" />
                <path d="M8.53 16.11a6 6 0 0 1 6.95 0" />
                <line x1="12" y1="20" x2="12.01" y2="20" />
              </svg>
            </div>
            <h2>Quick Stats</h2>
          </div>
          <div className="quick-stats-list">
            <div className="quick-stat-item">
              <span className="quick-stat-label">Checked In Today</span>
              <span className="quick-stat-value">{Math.round(stats.total_employees * stats.on_time_percentage / 100)}</span>
            </div>
            <div className="quick-stat-item">
              <span className="quick-stat-label">Pending Approval</span>
              <span className="quick-stat-value">{stats.late_arrivals}</span>
            </div>
            <div className="quick-stat-item">
              <span className="quick-stat-label">Active Branches</span>
              <span className="quick-stat-value">{branches.length}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
