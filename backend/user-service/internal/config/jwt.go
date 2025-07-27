package config

import (
	"os"
	"strconv"
	"time"
)

// JWTConfig содержит конфигурацию для JWT токенов
type JWTConfig struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// LoadJWTConfig загружает конфигурацию JWT из переменных окружения или использует значения по умолчанию
func LoadJWTConfig() JWTConfig {
	secretKey := getEnv("JWT_SECRET_KEY", "OMGMYKEY")

	accessExpiryMinutes := getEnvAsInt("JWT_ACCESS_EXPIRY_MINUTES", 15)
	refreshExpiryDays := getEnvAsInt("JWT_REFRESH_EXPIRY_DAYS", 7)

	return JWTConfig{
		SecretKey:     secretKey,
		AccessExpiry:  time.Duration(accessExpiryMinutes) * time.Minute,
		RefreshExpiry: time.Duration(refreshExpiryDays) * 24 * time.Hour,
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
