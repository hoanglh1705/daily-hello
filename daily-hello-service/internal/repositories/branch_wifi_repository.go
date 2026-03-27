package repositories

import (
	"context"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type BranchWifiRepository struct {
	db *gorm.DB
}

func NewBranchWifiRepository(db *gorm.DB) *BranchWifiRepository {
	return &BranchWifiRepository{db: db}
}

func (r *BranchWifiRepository) Create(ctx context.Context, wifi *models.BranchWifi) error {
	return r.db.WithContext(ctx).Create(wifi).Error
}

func (r *BranchWifiRepository) FindByID(ctx context.Context, id uint) (*models.BranchWifi, error) {
	var wifi models.BranchWifi
	err := r.db.WithContext(ctx).Preload("Branch").First(&wifi, id).Error
	if err != nil {
		return nil, err
	}
	return &wifi, nil
}

func (r *BranchWifiRepository) FindByBranchID(ctx context.Context, branchID uint) ([]models.BranchWifi, error) {
	var items []models.BranchWifi
	err := r.db.WithContext(ctx).Where("branch_id = ?", branchID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *BranchWifiRepository) FindByBSSID(ctx context.Context, bssid string) (*models.BranchWifi, error) {
	var wifi models.BranchWifi
	err := r.db.WithContext(ctx).Where("bssid = ?", bssid).First(&wifi).Error
	if err != nil {
		return nil, err
	}
	return &wifi, nil
}

func (r *BranchWifiRepository) Update(ctx context.Context, wifi *models.BranchWifi) error {
	return r.db.WithContext(ctx).Save(wifi).Error
}

func (r *BranchWifiRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.BranchWifi{}, id).Error
}
