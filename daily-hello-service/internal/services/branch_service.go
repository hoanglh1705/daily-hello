package services

import (
	"context"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"
)

type BranchService struct {
	repo repositories.BranchRepository
}

func NewBranchService(repo repositories.BranchRepository) *BranchService {
	return &BranchService{repo: repo}
}

func (s *BranchService) Create(ctx context.Context, req models.CreateBranchRequest) (*models.Branch, error) {
	branch := &models.Branch{
		BranchCode:       req.BranchCode,
		ParentBranchCode: req.ParentBranchCode,
		Name:             req.Name,
		Address:          req.Address,
		Lat:              req.Lat,
		Lng:              req.Lng,
		Radius:           req.Radius,
	}

	if err := s.repo.Create(ctx, branch); err != nil {
		return nil, appErrors.ErrInternal
	}

	return branch, nil
}

func (s *BranchService) GetByID(ctx context.Context, id uint) (*models.Branch, error) {
	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrBranchNotFound
	}
	return branch, nil
}

func (s *BranchService) Update(ctx context.Context, id uint, req models.UpdateBranchRequest) (*models.Branch, error) {
	branch, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrBranchNotFound
	}

	if req.ParentBranchCode != "" {
		branch.ParentBranchCode = req.ParentBranchCode
	}
	if req.Name != "" {
		branch.Name = req.Name
	}
	if req.Address != "" {
		branch.Address = req.Address
	}
	if req.Lat != nil {
		branch.Lat = req.Lat
	}
	if req.Lng != nil {
		branch.Lng = req.Lng
	}
	if req.Radius != nil {
		branch.Radius = req.Radius
	}
	if req.Status != "" {
		branch.Status = req.Status
	}

	if err := s.repo.Update(ctx, branch); err != nil {
		return nil, appErrors.ErrInternal
	}

	return branch, nil
}

func (s *BranchService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return appErrors.ErrBranchNotFound
	}

	return s.repo.Delete(ctx, id)
}

func (s *BranchService) List(ctx context.Context, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
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
