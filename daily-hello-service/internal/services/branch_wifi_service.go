package services

import (
	"slices"
	"context"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"
)

type BranchWifiService struct {
	repo       repositories.BranchWifiRepository
	branchRepo repositories.BranchRepository
}

func NewBranchWifiService(repo repositories.BranchWifiRepository, branchRepo repositories.BranchRepository) *BranchWifiService {
	return &BranchWifiService{repo: repo, branchRepo: branchRepo}
}

func (s *BranchWifiService) Create(ctx context.Context, req models.CreateBranchWifiRequest) (*models.BranchWifi, error) {
	// Validate branch exists
	_, err := s.branchRepo.FindByID(ctx, req.BranchID)
	if err != nil {
		return nil, appErrors.ErrBranchNotFound
	}

	wifi := &models.BranchWifi{
		Code:     req.Code,
		Name:     req.Name,
		BranchID: req.BranchID,
		SSID:     req.SSID,
		BSSID:    req.BSSID,
	}

	if err := s.repo.Create(ctx, wifi); err != nil {
		return nil, appErrors.ErrInternal
	}

	return wifi, nil
}

func (s *BranchWifiService) GetByID(ctx context.Context, id uint) (*models.BranchWifi, error) {
	wifi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}
	return wifi, nil
}

func (s *BranchWifiService) GetByBranchID(ctx context.Context, branchID uint, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
	items, total, err := s.repo.FindByBranchID(ctx, branchID, pq)
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

func (s *BranchWifiService) GetMyList(ctx context.Context, role string, branchID *uint, queryBranchID *uint, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
	var (
		items []models.BranchWifi
		total int64
		err   error
	)

	switch models.Role(role) {
	case models.RoleAdmin:
		if queryBranchID != nil {
			items, total, err = s.repo.FindByBranchID(ctx, *queryBranchID, pq)
		} else {
			items, total, err = s.repo.FindAll(ctx, pq)
		}
	case models.RoleManager:
		if branchID == nil {
			return nil, appErrors.ErrForbidden
		}

		branch, findErr := s.branchRepo.FindByID(ctx, *branchID)
		if findErr != nil {
			return nil, appErrors.ErrBranchNotFound
		}

		branchIDs := []uint{branch.ID}
		if branch.BranchCode != "" {
			children, childErr := s.branchRepo.FindByParentBranchCode(ctx, branch.BranchCode)
			if childErr != nil {
				return nil, appErrors.ErrInternal
			}
			for _, child := range children {
				branchIDs = append(branchIDs, child.ID)
			}
		}

		if queryBranchID != nil {
			allowed := slices.Contains(branchIDs, *queryBranchID)
			if !allowed {
				return nil, appErrors.ErrForbidden
			}
			items, total, err = s.repo.FindByBranchID(ctx, *queryBranchID, pq)
		} else {
			items, total, err = s.repo.FindByBranchIDs(ctx, branchIDs, pq)
		}
	default:
		if branchID == nil {
			return nil, appErrors.ErrForbidden
		}
		if queryBranchID != nil && *queryBranchID != *branchID {
			return nil, appErrors.ErrForbidden
		}
		items, total, err = s.repo.FindByBranchID(ctx, *branchID, pq)
	}

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

func (s *BranchWifiService) Update(ctx context.Context, id uint, req models.UpdateBranchWifiRequest) (*models.BranchWifi, error) {
	wifi, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrNotFound
	}

	if req.Code != "" {
		wifi.Code = req.Code
	}
	if req.Name != "" {
		wifi.Name = req.Name
	}
	if req.SSID != "" {
		wifi.SSID = req.SSID
	}
	if req.BSSID != "" {
		wifi.BSSID = req.BSSID
	}

	if err := s.repo.Update(ctx, wifi); err != nil {
		return nil, appErrors.ErrInternal
	}

	return wifi, nil
}

func (s *BranchWifiService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return appErrors.ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}
