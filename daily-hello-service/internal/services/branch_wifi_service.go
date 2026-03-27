package services

import (
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

func (s *BranchWifiService) GetByBranchID(ctx context.Context, branchID uint) ([]models.BranchWifi, error) {
	items, err := s.repo.FindByBranchID(ctx, branchID)
	if err != nil {
		return nil, appErrors.ErrInternal
	}
	return items, nil
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
