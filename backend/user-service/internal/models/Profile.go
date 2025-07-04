package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
