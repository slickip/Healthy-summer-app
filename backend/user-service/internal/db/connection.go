package db

import (
	"fmt"
	"log"
	"os"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/models"

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.FriendRequest{},
		&models.Achievement{},
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
