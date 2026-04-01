package repositories

import (
	"context"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type (
	BranchRepository interface {
		Create(ctx context.Context, branch *models.Branch) error
		FindByID(ctx context.Context, id uint) (*models.Branch, error)
		FindByParentBranchCode(ctx context.Context, parentBranchCode string) ([]models.Branch, error)
		Update(ctx context.Context, branch *models.Branch) error
		Delete(ctx context.Context, id uint) error
		List(ctx context.Context, pq models.PaginationQuery) ([]models.Branch, int64, error)
		ListByIDs(ctx context.Context, branchIDs []uint, pq models.PaginationQuery) ([]models.Branch, int64, error)
	}

	branchRepository struct {
		db *gorm.DB
	}
)

func NewBranchRepository(db *gorm.DB) BranchRepository {
	return &branchRepository{db: db}
}

func (r *branchRepository) Create(ctx context.Context, branch *models.Branch) error {
	return r.db.WithContext(ctx).Create(branch).Error
}

func (r *branchRepository) FindByID(ctx context.Context, id uint) (*models.Branch, error) {
	var branch models.Branch
	err := r.db.WithContext(ctx).Preload("WifiList").First(&branch, id).Error
	if err != nil {
		return nil, err
	}
	return &branch, nil
}

func (r *branchRepository) FindByParentBranchCode(ctx context.Context, parentBranchCode string) ([]models.Branch, error) {
	var branches []models.Branch
	err := r.db.WithContext(ctx).
		Where("parent_branch_code = ?", parentBranchCode).
		Find(&branches).Error
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (r *branchRepository) Update(ctx context.Context, branch *models.Branch) error {
	return r.db.WithContext(ctx).Save(branch).Error
}

func (r *branchRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Branch{}, id).Error
}

func (r *branchRepository) List(ctx context.Context, pq models.PaginationQuery) ([]models.Branch, int64, error) {
	return r.list(ctx, nil, pq)
}

func (r *branchRepository) ListByIDs(ctx context.Context, branchIDs []uint, pq models.PaginationQuery) ([]models.Branch, int64, error) {
	return r.list(ctx, branchIDs, pq)
}

func (r *branchRepository) list(ctx context.Context, branchIDs []uint, pq models.PaginationQuery) ([]models.Branch, int64, error) {
	var items []models.Branch
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Branch{})
	if len(branchIDs) > 0 {
		query = query.Where("id IN ?", branchIDs)
	}

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
