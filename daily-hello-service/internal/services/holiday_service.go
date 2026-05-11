package services

import (
	"context"
	"time"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"
)

type HolidayService struct {
	repo repositories.HolidayRepository
}

func NewHolidayService(repo repositories.HolidayRepository) *HolidayService {
	return &HolidayService{repo: repo}
}

func (s *HolidayService) Create(ctx context.Context, req models.CreateHolidayRequest, createdBy uint) (*models.Holiday, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, appErrors.ErrInvalidInput
	}

	holiday := &models.Holiday{
		Name:        req.Name,
		Date:        date,
		Description: req.Description,
		CreatedBy:   createdBy,
	}

	if err := s.repo.Create(ctx, holiday); err != nil {
		return nil, appErrors.ErrInternal
	}

	return holiday, nil
}

func (s *HolidayService) GetByID(ctx context.Context, id uint) (*models.Holiday, error) {
	holiday, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrHolidayNotFound
	}
	return holiday, nil
}

func (s *HolidayService) Update(ctx context.Context, id uint, req models.UpdateHolidayRequest) (*models.Holiday, error) {
	holiday, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrHolidayNotFound
	}

	if req.Name != "" {
		holiday.Name = req.Name
	}
	if req.Date != "" {
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, appErrors.ErrInvalidInput
		}
		holiday.Date = date
	}
	if req.Description != "" {
		holiday.Description = req.Description
	}

	if err := s.repo.Update(ctx, holiday); err != nil {
		return nil, appErrors.ErrInternal
	}

	return holiday, nil
}

func (s *HolidayService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return appErrors.ErrHolidayNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *HolidayService) List(ctx context.Context, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
	items, total, err := s.repo.List(ctx, pq)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return &models.PaginatedResponse{
		Items: items,
		Meta: models.PaginationMeta{
			Page:  pq.GetPage(),
			Limit: pq.GetLimit(),
			Total: total,
		},
	}, nil
}

func (s *HolidayService) GetByMonth(ctx context.Context, year, month int) ([]models.Holiday, error) {
	if year <= 0 || month < 1 || month > 12 {
		return nil, appErrors.ErrInvalidInput
	}

	from := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, -1) // last day of month

	items, err := s.repo.ListByDateRange(ctx, from, to)
	if err != nil {
		return nil, appErrors.ErrInternal
	}

	return items, nil
}
