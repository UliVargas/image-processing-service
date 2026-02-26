package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;size=24" json:"id"`
	Name      string         `gorm:"not null;" json:"name"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=2,max=100"`
	Email *string `json:"email" validate:"omitempty,email"`
}

type UpdatePasswordUserRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=6,max=32"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=32"`
}
