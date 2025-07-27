package models

import "time"

type Activity struct {
	ID             uint `gorm:"primaryKey"`
	UserID         uint `gorm:"not null"`
	ActivityTypeID uint `gorm:"not null"`
	Duration       int
	Intensity      string `gorm:"type:text;not null;check:intensity IN ('low','medium','high')"`
	Calories       int
	StartedAt      time.Time `gorm:"default:now()"`
	EndedAt        time.Time `gorm:"default:now()"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	ActivityType ActivityType `gorm:"foreignKey:ActivityTypeID"`
}

func (Activity) TableName() string {
	return "activities"
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

type ActivityType struct {
	ID                uint   `gorm:"primaryKey"`
	Name              string `gorm:"unique;not null"` // running, swimming и т.д.
	CaloriesPerMinute int    `gorm:"not null"`
}

func (ActivityType) TableName() string {
	return "activity_type"
}
