package db

import (
	"log"

	"github.com/slickip/Healthy-summer-app/backend/activity-service/internal/models" // импорт модели
	"gorm.io/gorm"
)

func SeedActivityType(db *gorm.DB) {
	var count int64
	if err := db.Model(&models.ActivityType{}).Count(&count).Error; err != nil {
		log.Printf("Error checking activity_type table: %v", err)
		return
	}

	if count > 0 {
		log.Println("activity_type already seeded — skipping")
		return
	}

	types := []models.ActivityType{
		{ID: 1, Name: "running", CaloriesPerMinute: 10},
		{ID: 2, Name: "cycling", CaloriesPerMinute: 8},
		{ID: 3, Name: "swimming", CaloriesPerMinute: 12},
		{ID: 4, Name: "yoga", CaloriesPerMinute: 4},
	}

	if err := db.Create(&types).Error; err != nil {
		log.Printf("Error seeding activity_type: %v", err)
		return
	}

	log.Println("Successfully seeded activity_type table.")
}
