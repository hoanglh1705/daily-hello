package models

import "time"

type AttendanceSummary struct {
	ID            uint             `json:"id" gorm:"primaryKey"`
	UserID        *uint            `json:"user_id" gorm:"index"`
	User          *User            `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Date          *time.Time       `json:"date" gorm:"type:date"`
	TotalHours    *float64         `json:"total_hours"`
	OvertimeHours *float64         `json:"overtime_hours"`
	Status        AttendanceStatus `json:"status" gorm:"type:varchar(20)"`
}
