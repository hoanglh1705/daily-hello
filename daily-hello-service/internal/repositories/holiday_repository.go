package repositories

import (
	"context"
	"time"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type (
	HolidayRepository interface {
		Create(ctx context.Context, holiday *models.Holiday) error
		FindByID(ctx context.Context, id uint) (*models.Holiday, error)
		Update(ctx context.Context, holiday *models.Holiday) error
		Delete(ctx context.Context, id uint) error
		List(ctx context.Context, pq models.PaginationQuery) ([]models.Holiday, int64, error)
		ListByDateRange(ctx context.Context, from, to time.Time) ([]models.Holiday, error)
	}

	holidayRepository struct {
		db *gorm.DB
	}
)

func NewHolidayRepository(db *gorm.DB) HolidayRepository {
	return &holidayRepository{db: db}
}

func (r *holidayRepository) Create(ctx context.Context, holiday *models.Holiday) error {
	return r.db.WithContext(ctx).Create(holiday).Error
}

func (r *holidayRepository) FindByID(ctx context.Context, id uint) (*models.Holiday, error) {
	var holiday models.Holiday
	err := r.db.WithContext(ctx).First(&holiday, id).Error
	if err != nil {
		return nil, err
	}
	return &holiday, nil
}

func (r *holidayRepository) Update(ctx context.Context, holiday *models.Holiday) error {
	return r.db.WithContext(ctx).Save(holiday).Error
}

func (r *holidayRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Holiday{}, id).Error
}

func (r *holidayRepository) List(ctx context.Context, pq models.PaginationQuery) ([]models.Holiday, int64, error) {
	var items []models.Holiday
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Holiday{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("date DESC").
		Offset(pq.GetOffset()).
		Limit(pq.GetLimit()).
		Find(&items).Error

	return items, total, err
}

func (r *holidayRepository) ListByDateRange(ctx context.Context, from, to time.Time) ([]models.Holiday, error) {
	var items []models.Holiday
	err := r.db.WithContext(ctx).
		Where("date >= ? AND date <= ?", from, to).
		Order("date ASC").
		Find(&items).Error
	return items, err
}
