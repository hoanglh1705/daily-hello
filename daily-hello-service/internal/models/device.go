package models

import "time"

type Device struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    *uint     `json:"user_id" gorm:"index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DeviceID  string    `json:"device_id" gorm:"type:varchar(100)"`
	IsTrusted *bool     `json:"is_trusted" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateDeviceRequest struct {
	UserID   uint   `json:"user_id" validate:"required"`
	DeviceID string `json:"device_id" validate:"required"`
}
