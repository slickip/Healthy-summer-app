package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New creates DB connection and runs AutoMigrate
func New() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "postgres"),
		getEnv("DB_USER", "healthyuser"),
		getEnv("DB_PASSWORD", "healthypass"),
		getEnv("DB_NAME", "healthydb"),
		getEnv("DB_PORT", "5432"),
	)

	var db *gorm.DB
	var err error

	// 5 попыток подключения
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database not ready yet (attempt %d/5): %v", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("failed to connect database after retries: %v", err)
	}

	if err := db.AutoMigrate(
		&models.Meals{},
		&models.Foods{},
		&models.WaterLogs{},
	); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	log.Println("Database connected and migrated successfully")
	return db
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
