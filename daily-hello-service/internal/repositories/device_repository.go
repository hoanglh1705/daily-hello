package repositories

import (
	"context"

	"daily-hello-service/internal/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

func (r *DeviceRepository) Create(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Create(device).Error
}

func (r *DeviceRepository) FindByID(ctx context.Context, id uint) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).Preload("User").First(&device, id).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Device, error) {
	var items []models.Device
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DeviceRepository) FindByDeviceID(ctx context.Context, deviceID string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) FindByUserIDAndDeviceID(ctx context.Context, userID uint, deviceID string) (*models.Device, error) {
	var device models.Device
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND device_id = ?", userID, deviceID).
		First(&device).Error
	if err != nil {
		return nil, err
	}
	return &device, nil
}

func (r *DeviceRepository) FindByStatus(ctx context.Context, status string) ([]models.Device, error) {
	var items []models.Device
	err := r.db.WithContext(ctx).Preload("User").Where("status = ?", status).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *DeviceRepository) Update(ctx context.Context, device *models.Device) error {
	return r.db.WithContext(ctx).Save(device).Error
}

func (r *DeviceRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Device{}, id).Error
}
