package models

import "time"

type Activity struct {
	id           uint `gorm:"primaryKey"`
	user_id      uint `gorm:"not null"`
	duration_min uint
	intensity    string    `gorm:"type:text;not null;check:intensity IN ('low','medium','high')"`
	StartedAt    time.Time `gorm:"default:now()"`
	EndedAt      time.Time `gorm:"default:now()"`
}

func (Activity) TableName() string {
	return "ativities"
}

type Steps struct {
	id        uint      `gorm:"primaryKey"`
	user_id   uint      `gorm:"not null"`
	step_cout int       `gorm:"not null"`
	date      time.Time `gorm:"type:date"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Steps) TableName() string {
	return "steps"
}
