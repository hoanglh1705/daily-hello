package repositories

import (
	"context"
	"strings"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Branch").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByCode(ctx context.Context, code string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(user).Error
}

func (r *UserRepository) List(ctx context.Context, q models.UserListQuery) ([]models.User, int64, error) {
	var items []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{}) //.Where("role = ?", models.RoleEmployee)

	if q.BranchID != nil {
		query = query.Where("branch_id = ?", *q.BranchID)
	}

	if keyword := strings.TrimSpace(q.Keyword); keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where(
			"name ILIKE ? OR code ILIKE ? OR email ILIKE ?",
			likeKeyword,
			likeKeyword,
			likeKeyword,
		)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Preload("Branch").
		Order("created_at DESC").
		Offset(q.GetOffset()).
		Limit(q.GetLimit()).
		Find(&items).Error

	return items, total, err
}
