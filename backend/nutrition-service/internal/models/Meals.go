package models

import "time"

type Meals struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	MealTime    time.Time `json:"meal_time"`
	Description string    `gorm:"not null" json:"description"`
	Calories    int       `gorm:"not null" json:"calories"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
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
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	VolumeML  int       `gorm:"not null" json:"volume_ml"`
	LoggedAt  time.Time `json:"logged_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (WaterLogs) TableName() string {
	return "water_logs"
}
