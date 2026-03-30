export type DashboardStats = {
  total_employees: number
  on_time_percentage: number
  on_time_change: number
  late_arrivals: number
  late_arrivals_change: number
}

export type AttendanceTrend = {
  day: string
  label: string
  check_in_count: number
  total: number
}

export type RecentActivity = {
  id: number
  user_name: string
  action: string
  time: string
}

export type DashboardData = {
  stats: DashboardStats
  trends: AttendanceTrend[]
  recent_activities: RecentActivity[]
}
