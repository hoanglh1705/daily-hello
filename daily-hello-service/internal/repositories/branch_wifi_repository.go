package repositories

import (
	"context"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type (
	BranchWifiRepository interface {
		Create(ctx context.Context, wifi *models.BranchWifi) error
		FindByID(ctx context.Context, id uint) (*models.BranchWifi, error)
		FindAll(ctx context.Context, pq models.PaginationQuery) ([]models.BranchWifi, int64, error)
		FindByBranchID(ctx context.Context, branchID uint, pq models.PaginationQuery) ([]models.BranchWifi, int64, error)
		FindByBranchIDs(ctx context.Context, branchIDs []uint, pq models.PaginationQuery) ([]models.BranchWifi, int64, error)
		FindByBSSID(ctx context.Context, bssid string) (*models.BranchWifi, error)
		Update(ctx context.Context, wifi *models.BranchWifi) error
		Delete(ctx context.Context, id uint) error
	}

	branchWifiRepository struct {
		db *gorm.DB
	}
)

func NewBranchWifiRepository(db *gorm.DB) BranchWifiRepository {
	return &branchWifiRepository{db: db}
}

func (r *branchWifiRepository) Create(ctx context.Context, wifi *models.BranchWifi) error {
	return r.db.WithContext(ctx).Create(wifi).Error
}

func (r *branchWifiRepository) FindByID(ctx context.Context, id uint) (*models.BranchWifi, error) {
	var wifi models.BranchWifi
	err := r.db.WithContext(ctx).First(&wifi, id).Preload("Branch").Error
	if err != nil {
		return nil, err
	}
	return &wifi, nil
}

func (r *branchWifiRepository) FindAll(ctx context.Context, pq models.PaginationQuery) ([]models.BranchWifi, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&models.BranchWifi{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.BranchWifi
	err := query.Preload("Branch").
		Order("created_at DESC").
		Offset(pq.GetOffset()).
		Limit(pq.GetLimit()).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *branchWifiRepository) FindByBranchID(ctx context.Context, branchID uint, pq models.PaginationQuery) ([]models.BranchWifi, int64, error) {
	return r.FindByBranchIDs(ctx, []uint{branchID}, pq)
}

func (r *branchWifiRepository) FindByBranchIDs(ctx context.Context, branchIDs []uint, pq models.PaginationQuery) ([]models.BranchWifi, int64, error) {
	var total int64
	query := r.db.WithContext(ctx).Model(&models.BranchWifi{}).Where("branch_id IN ?", branchIDs)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []models.BranchWifi
	err := query.Preload("Branch").
		Order("created_at DESC").
		Offset(pq.GetOffset()).
		Limit(pq.GetLimit()).
		Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *branchWifiRepository) FindByBSSID(ctx context.Context, bssid string) (*models.BranchWifi, error) {
	var wifi models.BranchWifi
	err := r.db.WithContext(ctx).Where("bssid = ?", bssid).First(&wifi).Error
	if err != nil {
		return nil, err
	}
	return &wifi, nil
}

func (r *branchWifiRepository) Update(ctx context.Context, wifi *models.BranchWifi) error {
	return r.db.WithContext(ctx).Save(wifi).Error
}

func (r *branchWifiRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.BranchWifi{}, id).Error
}
