package services

import (
	"context"
	"time"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/repositories"

	"gorm.io/gorm"
)

type DeviceService struct {
	repo *repositories.DeviceRepository
}

func NewDeviceService(repo *repositories.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

// Register registers a device for the current user.
// If the device is already registered for that user, returns the existing record.
// Otherwise creates a new record with status=pending.
func (s *DeviceService) Register(ctx context.Context, userID uint, req models.RegisterDeviceRequest) (*models.Device, error) {
	existing, err := s.repo.FindByUserIDAndDeviceID(ctx, userID, req.DeviceID)
	if err == nil {
		return existing, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, appErrors.ErrInternal
	}

	device := &models.Device{
		UserID:     &userID,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		Platform:   req.Platform,
		Model:      req.Model,
		Status:     models.DeviceStatusPending,
	}

	if err := s.repo.Create(ctx, device); err != nil {
		return nil, appErrors.ErrInternal
	}

	return device, nil
}

// GetStatus returns the device record for a given user and device_id.
func (s *DeviceService) GetStatus(ctx context.Context, userID uint, deviceID string) (*models.Device, error) {
	device, err := s.repo.FindByUserIDAndDeviceID(ctx, userID, deviceID)
	if err != nil {
		return nil, appErrors.ErrDeviceNotFound
	}
	return device, nil
}

// GetByUserID lists all devices belonging to a user.
func (s *DeviceService) GetByUserID(ctx context.Context, userID uint, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
	items, total, err := s.repo.FindByUserID(ctx, userID, pq)
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

// ListByStatus returns all devices filtered by status (admin use).
func (s *DeviceService) ListByStatus(ctx context.Context, status string, branchID *uint, pq models.PaginationQuery) (*models.PaginatedResponse, error) {
	items, total, err := s.repo.FindByStatus(ctx, status, branchID, pq)
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

// Approve sets a device status to approved.
func (s *DeviceService) Approve(ctx context.Context, adminID uint, id uint) (*models.Device, error) {
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrDeviceNotFound
	}

	now := time.Now()
	device.Status = models.DeviceStatusApproved
	device.ApprovedBy = &adminID
	device.ApprovedAt = &now

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, appErrors.ErrInternal
	}

	return device, nil
}

// Reject sets a device status to rejected.
func (s *DeviceService) Reject(ctx context.Context, adminID uint, id uint) (*models.Device, error) {
	device, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, appErrors.ErrDeviceNotFound
	}

	now := time.Now()
	device.Status = models.DeviceStatusRejected
	device.ApprovedBy = &adminID
	device.ApprovedAt = &now

	if err := s.repo.Update(ctx, device); err != nil {
		return nil, appErrors.ErrInternal
	}

	return device, nil
}
