package models

import "time"

// DashboardOverviewResponse contains summary statistics
type DashboardOverviewResponse struct {
	Summary          DashboardSummary   `json:"summary"`
	AttendanceTrends []AttendanceTrend  `json:"attendance_trends"`
	QuickStats       DashboardQuickStat `json:"quick_stats"`
}

// DashboardSummary holds aggregated data
type DashboardSummary struct {
	TotalEmployee int           `json:"total_employee"`
	OnTime        MetricSummary `json:"on_time"`
	LateArrival   MetricSummary `json:"late_arrival"`
}

// MetricSummary holds count/percentage and its trend
type MetricSummary struct {
	Percentage *float64 `json:"percentage,omitempty"`
	Count      *int     `json:"count,omitempty"`
	Trend      float64  `json:"trend"`
}

// AttendanceTrend holds data for a specific day
type AttendanceTrend struct {
	Day          string `json:"day"`
	Date         string `json:"date"`
	PresentCount int    `json:"present_count"`
}

// DashboardQuickStat holds generic dashboard quick stats
type DashboardQuickStat struct {
	CheckedInToday  int `json:"checked_in_today"`
	PendingApproval int `json:"pending_approval"`
	ActiveBranches  int `json:"active_branches"`
}

// DashboardRecentActivityResponse lists recent activities
type DashboardRecentActivityResponse struct {
	Items []RecentActivityItem `json:"items"`
}

// RecentActivityItem holds data about an event (Check-in, Check-out, Late arrival)
type RecentActivityItem struct {
	ID          int64     `json:"id"`
	UserName    string    `json:"user_name"`
	AvatarText  string    `json:"avatar_text"`
	ActionType  string    `json:"action_type"`
	Time        string    `json:"time"`
	Timestamp   time.Time `json:"timestamp"`
}
