package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LoginID      string    `gorm:"type:text;uniqueIndex;not null" json:"login_id"`
	Name         string    `gorm:"type:text;not null" json:"name"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}
