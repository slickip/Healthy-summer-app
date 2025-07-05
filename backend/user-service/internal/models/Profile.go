package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex:users_email_key;not null"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
