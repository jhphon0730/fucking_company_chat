package model

import (
	"g_kk_ch/pkg/apperror"
	"g_kk_ch/pkg/auth"
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

func (u *User) CheckPassword(password string) error {
	if err := auth.CompareHashAndPassword(u.PasswordHash, password); err != nil {
		return apperror.ErrUserInvalidPassword
	}
	return nil
}
