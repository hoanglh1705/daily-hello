package repositories

import (
	"context"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type BranchRepository struct {
	db *gorm.DB
}

func NewBranchRepository(db *gorm.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

func (r *BranchRepository) Create(ctx context.Context, branch *models.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *BranchRepository) FindByID(ctx context.Context, id uint) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.WithContext(ctx).Preload("WifiList").First(&branch, id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *BranchRepository) Update(ctx context.Context, branch *models.Branch) error {
	return r.db.WithContext(ctx).Save(branch).Error
}

func (r *BranchRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Branch{}, id).Error
}

func (r *BranchRepository) List(ctx context.Context, pq models.PaginationQuery) ([]models.Branch, int64, error) {
	var items []models.Branch
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Branch{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("WifiList").
		Order("created_at DESC").
		Offset(pq.GetOffset()).
		Limit(pq.GetLimit()).
		Find(&items).Error

	return items, total, err
}
