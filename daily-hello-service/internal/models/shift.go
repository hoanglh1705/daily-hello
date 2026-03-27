package models

import "time"

type Shift struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	BranchID  *uint     `json:"branch_id" gorm:"index"`
	Branch    *Branch   `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	StartTime string    `json:"start_time" gorm:"type:time"` // TIME type
	EndTime   string    `json:"end_time" gorm:"type:time"`   // TIME type
	CreatedAt time.Time `json:"created_at"`
}

type CreateShiftRequest struct {
	BranchID  *uint  `json:"branch_id"`
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

type UpdateShiftRequest struct {
	BranchID  *uint  `json:"branch_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
