package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID      uuid.UUID `json:"id"`
	LoginID string    `json:"login_id"`
	Name    string    `json:"name"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
