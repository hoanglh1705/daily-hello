package models

import "time"

type Branch struct {
	ID        uint         `json:"id" gorm:"primaryKey"`
	Name      string       `json:"name" gorm:"type:varchar(100);not null"`
	Address   string       `json:"address" gorm:"type:text"`
	Lat       *float64     `json:"lat" gorm:"type:double precision"`
	Lng       *float64     `json:"lng" gorm:"type:double precision"`
	Radius    *int         `json:"radius" gorm:"type:int"` // meters (geofence)
	WifiList  []BranchWifi `json:"wifi_list,omitempty" gorm:"foreignKey:BranchID"`
	Status    string       `json:"status" gorm:"type:varchar(20);default:'active'"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateBranchRequest struct {
	Name    string   `json:"name" validate:"required"`
	Address string   `json:"address"`
	Lat     *float64 `json:"lat"`
	Lng     *float64 `json:"lng"`
	Radius  *int     `json:"radius" validate:"omitempty,gt=0"`
}

type UpdateBranchRequest struct {
	Name    string   `json:"name"`
	Address string   `json:"address"`
	Lat     *float64 `json:"lat"`
	Lng     *float64 `json:"lng"`
	Radius  *int     `json:"radius" validate:"omitempty,gt=0"`
	Status  string   `json:"status" validate:"omitempty,oneof=active inactive"`
}
