package models

import "time"

type Meals struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null"`
	MealTime    time.Time
	Description string    `gorm:"not null"`
	Calories    int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (Meals) TableName() string {
	return "meals"
}

type Foods struct {
	ID               uint    `gorm:"primaryKey" json:"id"`
	Name             string  `gorm:"type:text;not null" json:"name"`
	CalloriesPer100g float32 `gorm:"not null" json:"callories_per_100g"`
	Proteins         float32 `gorm:"not null" json:"proteins"`
	Fats             float32 `gorm:"not null" json:"fats"`
	Carbs            float32 `gorm:"not null" json:"carbs"`
}

func (Foods) TableName() string {
	return "foods"
}

type WaterLogs struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	VolumeML  int  `gorm:"not null"`
	LoggedAt  time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (WaterLogs) TableName() string {
	return "water_logs"
}
