package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port              string
	AppEnv            string
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBTimeZone        string
	SwaggerUser       string
	SwaggerPassword   string
	JWTSecret         string
	JWTExpiryHours    string
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	R2PublicURL       string
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
		Port:              getEnv("PORT", "8080"),
		AppEnv:            getEnv("APP_ENV", "development"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USER"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_NAME"),
		DBSSLMode:         os.Getenv("DB_SSLMODE"),
		DBTimeZone:        os.Getenv("DB_TIMEZONE"),
		SwaggerUser:       os.Getenv("SWAGGER_USER"),
		SwaggerPassword:   os.Getenv("SWAGGER_PASSWORD"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		JWTExpiryHours:    os.Getenv("JWT_EXPIRATION_HOURS"),
		R2AccountID:       os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKeyID:     os.Getenv("R2_ACCESS_KEY_ID"),
		R2SecretAccessKey: os.Getenv("R2_SECRET_ACCESS_KEY"),
		R2BucketName:      os.Getenv("R2_BUCKET_NAME"),
		R2PublicURL:       os.Getenv("R2_PUBLIC_URL"),
	}

	return Config
}
