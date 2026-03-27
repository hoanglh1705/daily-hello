package repositories

import (
	"context"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type ShiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) *ShiftRepository {
	return &ShiftRepository{db: db}
}

func (r *ShiftRepository) Create(ctx context.Context, shift *models.Shift) error {
	return r.db.WithContext(ctx).Create(shift).Error
}

func (r *ShiftRepository) FindByID(ctx context.Context, id uint) (*models.Shift, error) {
	var shift models.Shift
	err := r.db.WithContext(ctx).Preload("Branch").First(&shift, id).Error
	if err != nil {
		return nil, err
	}
	return &shift, nil
}

func (r *ShiftRepository) FindByBranchID(ctx context.Context, branchID uint) ([]models.Shift, error) {
	var items []models.Shift
	err := r.db.WithContext(ctx).Where("branch_id = ?", branchID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ShiftRepository) Update(ctx context.Context, shift *models.Shift) error {
	return r.db.WithContext(ctx).Save(shift).Error
}

func (r *ShiftRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Shift{}, id).Error
}
