package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleManager  Role = "manager"
	RoleEmployee Role = "employee"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Code      string    `json:"code" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" gorm:"type:varchar(150);uniqueIndex;not null"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-" gorm:"type:text;not null"`
	Role      Role      `json:"role" gorm:"type:varchar(20);not null"`
	BranchID  *uint     `json:"branch_id" gorm:"index"`
	Branch    *Branch   `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	Status    string    `json:"status" gorm:"type:varchar(20);default:'active'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Code     string `json:"code" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Phone    string `json:"phone"`
	Password string `json:"password" validate:"required,min=6"`
	Role     Role   `json:"role" validate:"required,oneof=admin manager employee"`
	BranchID *uint  `json:"branch_id"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     Role   `json:"role" validate:"omitempty,oneof=admin manager employee"`
	BranchID *uint  `json:"branch_id"`
	Status   string `json:"status" validate:"omitempty,oneof=active inactive"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
