package session

import (
	"image-processing-service/internal/modules/user"
	"time"
)

type Session struct {
	ID        string    `gorm:"primaryKey;size=24" json:"id"`
	TokenHash string    `gorm:"not null;" json:"token_hash"`
	AccessJti string    `gorm:"not null;" json:"accessJti"`
	UserID    string    `gorm:"not null;" json:"user_id"`
	User      user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ExpiresAt time.Time `gorm:"not null;" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateSessionRequest struct {
	TokenHash string    `json:"token_hash"`
	AccessJti string    `json:"accessJti"`
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"created_at"`
}

type UpdateSessionRequest struct {
	SessionID    string    `json:"session_id"`
	NewTokenHash string    `json:"new_token_hash"`
	NewAccessJti string    `json:"new_access_jti"`
	ExpiresAt    time.Time `json:"expires_at"`
}
