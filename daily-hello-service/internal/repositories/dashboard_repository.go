package repositories

import (
	"daily-hello-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTotalEmployee(branchIDs []uint, branchID *int64) (int, error)
	GetOnTimeCount(branchIDs []uint, branchID *int64, startOfDay time.Time, endOfDay time.Time) (int, error)
	GetLateCount(branchIDs []uint, branchID *int64, startOfDay time.Time, endOfDay time.Time) (int, error)
	GetTotalCheckIn(branchIDs []uint, branchID *int64, startOfDay time.Time, endOfDay time.Time) (int, error)
	GetAttendanceTrends(branchIDs []uint, branchID *int64, fromDate time.Time, toDate time.Time) ([]models.AttendanceTrend, error)
	GetPendingDeviceApproval(branchIDs []uint) (int, error)
	GetActiveBranches(branchIDs []uint) (int, error)
	GetRecentActivities(branchIDs []uint, branchID *int64, startOfDay time.Time, endOfDay time.Time, limit int) ([]models.RecentActivityItem, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func applyBranchScope(query *gorm.DB, column string, branchIDs []uint, branchID *int64) *gorm.DB {
	if branchID != nil && *branchID > 0 {
		return query.Where(column+" = ?", *branchID)
	}
	if len(branchIDs) > 0 {
		return query.Where(column+" IN ?", branchIDs)
	}
	return query
}

func (r *dashboardRepository) GetTotalEmployee(branchIDs []uint, branchID *int64) (int, error) {
	var count int64
	query := r.db.Model(&models.User{}).Where("role <> ?", models.RoleAdmin)
	query = applyBranchScope(query, "branch_id", branchIDs, branchID)

	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetOnTimeCount(branchIDs []uint, branchID *int64, start time.Time, end time.Time) (int, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).Where("check_in_time BETWEEN ? AND ?", start, end).Where("check_in_status = ?", "on_time")
	query = applyBranchScope(query, "branch_id", branchIDs, branchID)
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetLateCount(branchIDs []uint, branchID *int64, start time.Time, end time.Time) (int, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).Where("check_in_time BETWEEN ? AND ?", start, end).Where("check_in_status = ?", "late")
	query = applyBranchScope(query, "branch_id", branchIDs, branchID)
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetTotalCheckIn(branchIDs []uint, branchID *int64, start time.Time, end time.Time) (int, error) {
	var count int64
	query := r.db.Model(&models.Attendance{}).Where("check_in_time BETWEEN ? AND ?", start, end)
	query = applyBranchScope(query, "branch_id", branchIDs, branchID)
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetAttendanceTrends(branchIDs []uint, branchID *int64, fromDate time.Time, toDate time.Time) ([]models.AttendanceTrend, error) {
	var results []models.AttendanceTrend

	query := r.db.Table("attendances").
		Select("TO_CHAR(DATE(check_in_time), 'YYYY-MM-DD') as date, TRIM(TO_CHAR(DATE(check_in_time), 'Dy')) as day, count(id) as present_count").
		Where("check_in_time BETWEEN ? AND ?", fromDate, toDate)

	query = applyBranchScope(query, "branch_id", branchIDs, branchID)

	err := query.Group("DATE(check_in_time)").Order("DATE(check_in_time) ASC").Scan(&results).Error
	return results, err
}

func (r *dashboardRepository) GetPendingDeviceApproval(branchIDs []uint) (int, error) {
	var count int64
	query := r.db.Model(&models.Device{}).Where("status = ?", "pending")
	if len(branchIDs) > 0 {
		query = query.Joins("JOIN users ON users.id = devices.user_id").Where("users.branch_id IN ?", branchIDs)
	}
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetActiveBranches(branchIDs []uint) (int, error) {
	var count int64
	query := r.db.Model(&models.Branch{})
	if len(branchIDs) > 0 {
		query = query.Where("id IN ?", branchIDs)
	}
	err := query.Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) GetRecentActivities(branchIDs []uint, branchID *int64, start time.Time, end time.Time, limit int) ([]models.RecentActivityItem, error) {
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

	query = applyBranchScope(query, "attendances.branch_id", branchIDs, branchID)

	err := query.Order("timestamp DESC").Limit(limit).Scan(&results).Error
	return results, err
}
