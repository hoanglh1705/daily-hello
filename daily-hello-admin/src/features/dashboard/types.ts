export type DashboardOverviewResponse = {
  summary: {
    total_employee: number
    on_time: { percentage: number; trend: number }
    late_arrival: { count: number; trend: number }
  }
  attendance_trends: Array<{
    day: string
    date: string
    present_count: number
  }>
  quick_stats: {
    checked_in_today: number
    pending_approval: number
    active_branches: number
  }
}

export type RecentActivityItem = {
  id: number
  user_name: string
  avatar_text: string
  action_type: string
  time: string
  timestamp: string
}

export type DashboardRecentActivitiesResponse = {
  items: RecentActivityItem[]
}
