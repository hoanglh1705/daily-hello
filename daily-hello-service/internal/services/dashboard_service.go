package services

import (
	"daily-hello-service/internal/models"
	"daily-hello-service/internal/repositories"
	"log"
	"time"
)

type DashboardService struct {
	repo repositories.DashboardRepository
}

func NewDashboardService(repo repositories.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetOverview(branchIDs []uint, branchID *int64, date time.Time) (*models.DashboardOverviewResponse, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	prevStartOfDay := startOfDay.AddDate(0, 0, -1)
	prevEndOfDay := endOfDay.AddDate(0, 0, -1)

	totalEmployee, _ := s.repo.GetTotalEmployee(branchIDs, branchID)
	onTimeCount, _ := s.repo.GetOnTimeCount(branchIDs, branchID, startOfDay, endOfDay)
	lateCount, _ := s.repo.GetLateCount(branchIDs, branchID, startOfDay, endOfDay)
	totalCheckIn, _ := s.repo.GetTotalCheckIn(branchIDs, branchID, startOfDay, endOfDay)

	prevOnTimeCount, _ := s.repo.GetOnTimeCount(branchIDs, branchID, prevStartOfDay, prevEndOfDay)
	prevLateCount, _ := s.repo.GetLateCount(branchIDs, branchID, prevStartOfDay, prevEndOfDay)
	prevTotalCheckIn, _ := s.repo.GetTotalCheckIn(branchIDs, branchID, prevStartOfDay, prevEndOfDay)

	var onTimePercentage, prevOnTimePercentage float64
	if totalCheckIn > 0 {
		onTimePercentage = float64(onTimeCount) / float64(totalCheckIn) * 100
	}
	if prevTotalCheckIn > 0 {
		prevOnTimePercentage = float64(prevOnTimeCount) / float64(prevTotalCheckIn) * 100
	}

	onTimeTrend := onTimePercentage - prevOnTimePercentage
	lateTrend := float64(lateCount - prevLateCount)

	sevenDaysAgo := startOfDay.AddDate(0, 0, -6)
	trends, err := s.repo.GetAttendanceTrends(branchIDs, branchID, sevenDaysAgo, endOfDay)
	if err != nil {
		log.Printf("Error getting trends: %v", err)
		trends = []models.AttendanceTrend{}
	}

	pendingApproval, _ := s.repo.GetPendingDeviceApproval(branchIDs)
	activeBranches, _ := s.repo.GetActiveBranches(branchIDs)

	onTimeP := &onTimePercentage
	lateC := &lateCount

	return &models.DashboardOverviewResponse{
		Summary: models.DashboardSummary{
			TotalEmployee: totalEmployee,
			OnTime: models.MetricSummary{
				Percentage: onTimeP,
				Trend:      onTimeTrend,
			},
			LateArrival: models.MetricSummary{
				Count: lateC,
				Trend: lateTrend,
			},
		},
		AttendanceTrends: trends,
		QuickStats: models.DashboardQuickStat{
			CheckedInToday:  totalCheckIn,
			PendingApproval: pendingApproval,
			ActiveBranches:  activeBranches,
		},
	}, nil
}

func (s *DashboardService) GetRecentActivities(branchIDs []uint, branchID *int64, date time.Time, limit int) (*models.DashboardRecentActivityResponse, error) {
	if limit <= 0 {
		limit = 10
	}
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	items, err := s.repo.GetRecentActivities(branchIDs, branchID, startOfDay, endOfDay, limit)
	if err != nil {
		return nil, err
	}

	if items == nil {
		items = []models.RecentActivityItem{}
	}

	return &models.DashboardRecentActivityResponse{
		Items: items,
	}, nil
}
