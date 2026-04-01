package repositories

import (
	"context"
	"time"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) *AttendanceRepository {
	return &AttendanceRepository{db: db}
}

func (r *AttendanceRepository) Create(ctx context.Context, att *models.Attendance) error {
	return r.db.WithContext(ctx).Create(att).Error
}

func (r *AttendanceRepository) FindByID(ctx context.Context, id uint) (*models.Attendance, error) {
	var att models.Attendance
	err := r.db.WithContext(ctx).Preload("User").Preload("Branch").First(&att, id).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *AttendanceRepository) FindTodayCheckIn(ctx context.Context, userID uint, from, to time.Time) (*models.Attendance, error) {
	var att models.Attendance

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND check_in_time >= ? AND check_in_time < ?",
			userID, from, to).
		First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *AttendanceRepository) FindTodayByUserID(ctx context.Context, userID uint, from, to time.Time) (*models.Attendance, error) {
	var att models.Attendance

	err := r.db.WithContext(ctx).
		Preload("User").Preload("Branch").
		Where("user_id = ? AND check_in_time >= ? AND check_in_time < ?",
			userID, from, to).
		First(&att).Error
	if err != nil {
		return nil, err
	}
	return &att, nil
}

func (r *AttendanceRepository) UpdateCheckOut(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Attendance{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AttendanceRepository) UpdateStatus(ctx context.Context, id uint, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&models.Attendance{}).Where("id = ?", id).Updates(updates).Error
}

func (r *AttendanceRepository) List(ctx context.Context, filter models.AttendanceFilter, pq models.PaginationQuery) ([]models.Attendance, int64, error) {
	var items []models.Attendance
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Attendance{})

	if filter.UserID > 0 {
		query = query.Where("user_id = ?", filter.UserID)
	}
	if filter.BranchID > 0 {
		query = query.Where("branch_id = ?", filter.BranchID)
	} else if len(filter.BranchIDs) > 0 {
		query = query.Where("branch_id IN ?", filter.BranchIDs)
	}
	if filter.DateFrom != "" {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("created_at <= ?", filter.DateTo)
	}
	if filter.Status != "" {
		query = query.Where("checkin_status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("User").Preload("Branch").
		Order("created_at DESC").
		Offset(pq.GetOffset()).
		Limit(pq.GetLimit()).
		Find(&items).Error

	return items, total, err
}
