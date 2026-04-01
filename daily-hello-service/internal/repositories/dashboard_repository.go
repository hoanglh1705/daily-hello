package repositories

import (
	"daily-hello-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTotalEmployee(role models.Role, currentBranchID *uint, branchID *int64) (int, error)
	GetOnTimeCount(branchID *int64, startOfDay time.Time, endOfDay time.Time) (int, error)
	GetLateCount(branchID *int64, startOfDay time.Time, endOfDay time.Time) (int, error)
	GetTotalCheckIn(branchID *int64, startOfDay time.Time, endOfDay time.Time) (int, error)
	GetAttendanceTrends(branchID *int64, fromDate time.Time, toDate time.Time) ([]models.AttendanceTrend, error)
	GetPendingDeviceApproval() (int, error)
	GetActiveBranches() (int, error)
	GetRecentActivities(branchID *int64, startOfDay time.Time, endOfDay time.Time, limit int) ([]models.RecentActivityItem, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetTotalEmployee(role models.Role, currentBranchID *uint, branchID *int64) (int, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("role <> ?", models.RoleAdmin)

	switch role {
	case models.RoleAdmin:
		if branchID != nil && *branchID > 0 {
			query = query.Where("branch_id = ?", *branchID)
		}
	case models.RoleManager:
		if currentBranchID == nil || *currentBranchID == 0 {
			return 0, nil
		}

		accessibleBranches := r.db.Model(&models.Branch{}).
			Select("id").
			Where("id = ?", *currentBranchID).
			Or("parent_branch_code = (?)",
				r.db.Model(&models.Branch{}).
					Select("branch_code").
					Where("id = ? AND status = ?", *currentBranchID, "active"),
			)

		query = query.Where("branch_id IN (?)", accessibleBranches)
		if branchID != nil && *branchID > 0 {
			query = query.Where("branch_id = ?", *branchID)
		}
	default:
		if currentBranchID == nil || *currentBranchID == 0 {
			return 0, nil
		}
		query = query.Where("branch_id = ?", *currentBranchID)
	}

	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetOnTimeCount(branchID *int64, start time.Time, end time.Time) (int, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).Where("check_in_time BETWEEN ? AND ?", start, end).Where("check_in_status = ?", "on_time")
	if branchID != nil && *branchID > 0 {
		query = query.Where("branch_id = ?", *branchID)
	}
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetLateCount(branchID *int64, start time.Time, end time.Time) (int, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).Where("check_in_time BETWEEN ? AND ?", start, end).Where("check_in_status = ?", "late")
	if branchID != nil && *branchID > 0 {
		query = query.Where("branch_id = ?", *branchID)
	}
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetTotalCheckIn(branchID *int64, start time.Time, end time.Time) (int, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).Where("check_in_time BETWEEN ? AND ?", start, end)
	if branchID != nil && *branchID > 0 {
		query = query.Where("branch_id = ?", *branchID)
	}
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetAttendanceTrends(branchID *int64, fromDate time.Time, toDate time.Time) ([]models.AttendanceTrend, error) {
	var results []models.AttendanceTrend

	query := r.db.Table("attendances").
		Select("TO_CHAR(DATE(check_in_time), 'YYYY-MM-DD') as date, TRIM(TO_CHAR(DATE(check_in_time), 'Dy')) as day, count(id) as present_count").
		Where("check_in_time BETWEEN ? AND ?", fromDate, toDate)

	if branchID != nil && *branchID > 0 {
		query = query.Where("branch_id = ?", *branchID)
	}

	err := query.Group("DATE(check_in_time)").Order("DATE(check_in_time) ASC").Scan(&results).Error
	return results, err
}

func (r *dashboardRepository) GetPendingDeviceApproval() (int, error) {
	var count int64
	err := r.db.Model(&models.Device{}).Where("status = ?", "pending").Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetActiveBranches() (int, error) {
	var count int64
	err := r.db.Model(&models.Branch{}).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetRecentActivities(branchID *int64, start time.Time, end time.Time, limit int) ([]models.RecentActivityItem, error) {
	var results []models.RecentActivityItem

	query := r.db.Table("attendances").
		Select("attendances.id, users.name as user_name, "+
			"SUBSTRING(users.name FROM 1 FOR 1) as avatar_text, "+
			"CASE "+
			"  WHEN check_out_time IS NOT NULL THEN "+
			"    CASE WHEN check_out_time::time > '17:00:00'::time THEN 'Depart soon' ELSE 'Check-out' END "+
			"  WHEN check_in_time::time > '08:00:00'::time THEN 'Late arrival' "+
			"  ELSE 'Check-in' "+
			"END as action_type, "+
			"TO_CHAR(COALESCE(check_out_time, check_in_time), 'HH24:MI') as time, "+
			"COALESCE(check_out_time, check_in_time) as timestamp").
		Joins("JOIN users ON users.id = attendances.user_id").
		Where("check_in_time BETWEEN ? AND ?", start, end)

	if branchID != nil && *branchID > 0 {
		query = query.Where("attendances.branch_id = ?", *branchID)
	}

	err := query.Order("timestamp DESC").Limit(limit).Scan(&results).Error
	return results, err
}
