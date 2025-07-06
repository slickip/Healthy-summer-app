package models

import "time"

type Activity struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null"`
	Type      string `gorm:"not null"`
	Duration  int
	Intensity string `gorm:"type:text;not null;check:intensity IN ('low','medium','high')"`
	Calories  int
	StartedAt time.Time `gorm:"default:now()"`
	EndedAt   time.Time `gorm:"default:now()"`
}

func (Activity) TableName() string {
	return "ativities"
}

type Steps struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	StepCount int       `gorm:"not null"`
	Date      time.Time `gorm:"type:date"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Steps) TableName() string {
	return "steps"
}
