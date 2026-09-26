package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port            string
	AppEnv          string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	DBTimeZone      string
	SwaggerUser     string
	SwaggerPassword string
}

var Config *AppConfig

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found or could not be loaded, using environment variables")
	}

	Config = &AppConfig{
		Port:            getEnv("PORT", "8080"),
		AppEnv:          getEnv("APP_ENV", "development"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "postgres_db"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		DBTimeZone:      getEnv("DB_TIMEZONE", "UTC"),
		SwaggerUser:     getEnv("SWAGGER_USER", "admin"),
		SwaggerPassword: getEnv("SWAGGER_PASSWORD", "admin123"),
	}

	return Config
}
