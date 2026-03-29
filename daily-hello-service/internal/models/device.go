package models

import "time"

type DeviceStatus string

const (
	DeviceStatusPending  DeviceStatus = "pending"
	DeviceStatusApproved DeviceStatus = "approved"
	DeviceStatusRejected DeviceStatus = "rejected"
)

type Device struct {
	ID         uint         `json:"id" gorm:"primaryKey"`
	UserID     *uint        `json:"user_id" gorm:"index"`
	User       *User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DeviceID   string       `json:"device_id" gorm:"type:varchar(100)"`
	DeviceName string       `json:"device_name" gorm:"type:varchar(200)"`
	Platform   string       `json:"platform" gorm:"type:varchar(20)"`
	Model      string       `json:"model" gorm:"type:varchar(100)"`
	Status     DeviceStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	ApprovedBy *uint        `json:"approved_by,omitempty" gorm:"index"`
	ApprovedAt *time.Time   `json:"approved_at,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
}

type RegisterDeviceRequest struct {
	DeviceID   string `json:"device_id" validate:"required"`
	DeviceName string `json:"device_name"`
	Platform   string `json:"platform"`
	Model      string `json:"model"`
}

type DeviceStatusQuery struct {
	DeviceID string `query:"device_id" validate:"required"`
}

type DeviceListQuery struct {
	Status string `query:"status"`
}
