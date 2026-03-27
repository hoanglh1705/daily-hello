package repositories

import (
	"context"
	"time"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type AttendanceSummaryRepository struct {
	db *gorm.DB
}

func NewAttendanceSummaryRepository(db *gorm.DB) *AttendanceSummaryRepository {
	return &AttendanceSummaryRepository{db: db}
}

func (r *AttendanceSummaryRepository) Create(ctx context.Context, summary *models.AttendanceSummary) error {
	return r.db.WithContext(ctx).Create(summary).Error
}

func (r *AttendanceSummaryRepository) FindByID(ctx context.Context, id uint) (*models.AttendanceSummary, error) {
	var summary models.AttendanceSummary
	err := r.db.WithContext(ctx).Preload("User").First(&summary, id).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *AttendanceSummaryRepository) FindByUserIDAndDate(ctx context.Context, userID uint, date time.Time) (*models.AttendanceSummary, error) {
	var summary models.AttendanceSummary
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", userID, date.Format("2006-01-02")).
		First(&summary).Error
	if err != nil {
		return nil, err
	}
	return &summary, nil
}

func (r *AttendanceSummaryRepository) Update(ctx context.Context, summary *models.AttendanceSummary) error {
	return r.db.WithContext(ctx).Save(summary).Error
}

func (r *AttendanceSummaryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AttendanceSummary{}, id).Error
}
