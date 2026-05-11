package models

import "time"

type Holiday struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null"`
	Date        time.Time `json:"date" gorm:"type:date;not null;index"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedBy   uint      `json:"created_by" gorm:"index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateHolidayRequest struct {
	Name        string `json:"name" validate:"required"`
	Date        string `json:"date" validate:"required"` // format: 2006-01-02
	Description string `json:"description"`
}

type UpdateHolidayRequest struct {
	Name        string `json:"name"`
	Date        string `json:"date"` // format: 2006-01-02
	Description string `json:"description"`
}

type HolidayFilter struct {
	Year  int `query:"year"`
	Month int `query:"month"`
}
