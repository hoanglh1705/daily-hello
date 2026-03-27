package models

import "time"

type BranchWifi struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code" gorm:"type:varchar(100);not null"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	BranchID  uint      `json:"branch_id" gorm:"index;not null"`
	Branch    *Branch   `json:"-" gorm:"foreignKey:BranchID" swaggerignore:"true"`
	SSID      string    `json:"ssid" gorm:"column:ssid;type:varchar(100)"`
	BSSID     string    `json:"bssid" gorm:"column:bssid;type:varchar(100)"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBranchWifiRequest struct {
	Code     string `json:"code" validate:"required"`
	Name     string `json:"name" validate:"required"`
	BranchID uint   `json:"branch_id" validate:"required"`
	SSID     string `json:"ssid"`
	BSSID    string `json:"bssid"`
}

type UpdateBranchWifiRequest struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	SSID  string `json:"ssid"`
	BSSID string `json:"bssid"`
}
