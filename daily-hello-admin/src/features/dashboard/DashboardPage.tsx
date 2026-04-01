import { useEffect, useState } from 'react'
import { getDashboardOverview, getRecentActivities } from './api'
import { getBranches } from '@/features/branch/api'
import type { DashboardOverviewResponse, RecentActivityItem } from './types'
import type { Branch } from '@/features/branch/types'
import { getCurrentBranchId } from '@/services/tokenStorage'

function formatDate(date: Date): string {
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export default function DashboardPage() {
  const [overview, setOverview] = useState<DashboardOverviewResponse | null>(null)
  const [activities, setActivities] = useState<RecentActivityItem[]>([])
  const [branches, setBranches] = useState<Branch[]>([])
  const [selectedBranch, setSelectedBranch] = useState<number | null>(getCurrentBranchId())
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    getBranches({ page: 1, limit: 100 })
      .then((res) => setBranches(res.data.items))
      .catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    const params = selectedBranch ? { branch_id: selectedBranch } : undefined
    
    Promise.all([
      getDashboardOverview(params),
      getRecentActivities(params)
    ])
    .then(([overviewRes, activitiesRes]) => {
      setOverview(overviewRes.data)
      setActivities(activitiesRes.data.items || [])
    })
    .catch(() => {
      // Handle error gracefully if needed
    })
    .finally(() => setLoading(false))
  }, [selectedBranch])

  const trends = overview?.attendance_trends || []
  const maxCount = Math.max(...trends.map((t) => t.present_count), 1)
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
              {loading ? '—' : overview?.summary.total_employee.toLocaleString() || '0'}
            </span>
          </div>
        </div>
      </div>

      <div className="dashboard-stats-row">
        <div className="stats-card stats-card-success">
          <span className="stats-card-label">On Time</span>
          <div className="stats-card-value-row">
            <span className="stats-card-value">{loading ? '—' : `${overview?.summary.on_time.percentage?.toFixed(1) || 0}%`}</span>
            {!!overview?.summary.on_time.trend && (
              <span className={`stats-change ${overview.summary.on_time.trend > 0 ? 'positive' : 'negative'}`}>
                {overview.summary.on_time.trend > 0 ? '+' : ''}{overview.summary.on_time.trend.toFixed(1)}%
              </span>
            )}
          </div>
        </div>
        <div className="stats-card stats-card-danger">
          <span className="stats-card-label">Late Arrival</span>
          <div className="stats-card-value-row">
            <span className="stats-card-value">{loading ? '—' : (overview?.summary.late_arrival.count || 0)}</span>
            {!!overview?.summary.late_arrival.trend && (
              <span className={`stats-change ${overview.summary.late_arrival.trend < 0 ? 'positive' : 'negative'}`}>
                {overview.summary.late_arrival.trend > 0 ? '+' : ''}{overview.summary.late_arrival.trend}
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
            <div className="chart-bar-group" key={t.date}>
              <div className="chart-bar-wrapper">
                <div
                  className={`chart-bar ${i === trends.findIndex((x) => x.present_count === maxCount) ? 'chart-bar-highlight' : ''}`}
                  style={{ height: `${maxCount > 0 ? (t.present_count / maxCount) * 100 : 0}%` }}
                />
              </div>
              <span className="chart-label">{t.day}</span>
            </div>
          ))}
          {trends.length === 0 && !loading && (
            <div className="text-gray-400 text-sm mt-4 text-center" style={{ width: '100%', position: 'absolute' }}>No trend data</div>
          )}
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
            {activities.map((a) => (
              <div className="activity-item" key={a.id}>
                <div className="activity-avatar">
                  {a.avatar_text || a.user_name.charAt(0)}
                </div>
                <div className="activity-info">
                  <span className="activity-name">{a.user_name}</span>
                  <span className="activity-action">{a.action_type}</span>
                </div>
                <span className="activity-time">{a.time}</span>
              </div>
            ))}
            {activities.length === 0 && !loading && (
              <div style={{ color: '#9ca3af', fontSize: '0.875rem', marginTop: '1rem', textAlign: 'center' }}>No recent activities</div>
            )}
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
              <span className="quick-stat-value">{loading ? '—' : overview?.quick_stats?.checked_in_today || 0}</span>
            </div>
            <div className="quick-stat-item">
              <span className="quick-stat-label">Pending Approval</span>
              <span className="quick-stat-value">{loading ? '—' : overview?.quick_stats?.pending_approval || 0}</span>
            </div>
            <div className="quick-stat-item">
              <span className="quick-stat-label">Active Branches</span>
              <span className="quick-stat-value">{loading ? '—' : overview?.quick_stats?.active_branches || 0}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
