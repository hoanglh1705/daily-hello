package models

import "time"

type AttendanceStatus string

const (
	StatusOnTime AttendanceStatus = "on_time"
	StatusLate   AttendanceStatus = "late"
	StatusAbsent AttendanceStatus = "absent"
)

type Attendance struct {
	ID           uint             `json:"id" gorm:"primaryKey"`
	UserID       uint             `json:"user_id" gorm:"index;not null"`
	User         *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	BranchID     uint             `json:"branch_id" gorm:"index;not null"`
	Branch       *Branch          `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	CheckInTime  *time.Time       `json:"check_in_time" gorm:"index"`
	CheckOutTime *time.Time       `json:"check_out_time"`
	CheckInLat   *float64         `json:"check_in_lat" gorm:"type:double precision"`
	CheckInLng   *float64         `json:"check_in_lng" gorm:"type:double precision"`
	CheckOutLat  *float64         `json:"check_out_lat" gorm:"type:double precision"`
	CheckOutLng  *float64         `json:"check_out_lng" gorm:"type:double precision"`
	WifiBSSID    string           `json:"wifi_bssid" gorm:"column:wifi_bssid;type:varchar(100)"`
	DeviceID     string           `json:"device_id" gorm:"type:varchar(100)"`
	Status       AttendanceStatus `json:"status" gorm:"type:varchar(20)"`
	CreatedAt    time.Time        `json:"created_at"`
}

type AttendanceRequest struct {
	Lat       float64 `json:"lat" validate:"required"`
	Lng       float64 `json:"lng" validate:"required"`
	WifiBSSID string  `json:"wifi_bssid"`
	WifiSSID  string  `json:"wifi_ssid"`
	DeviceID  string  `json:"device_id" validate:"required"`
	BranchID  uint    `json:"branch_id"`
}

type AttendanceFilter struct {
	UserID   uint   `query:"user_id"`
	BranchID uint   `query:"branch_id"`
	DateFrom string `query:"date_from"`
	DateTo   string `query:"date_to"`
	Status   string `query:"status"`
}
