package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port        string
	AppEnv      string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	DBTimeZone  string
}

var Config *AppConfig

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found or could not be loaded, using environment variables")
	}

	Config = &AppConfig{
		Port:       os.Getenv("PORT"),
		AppEnv:     os.Getenv("APP_ENV"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBTimeZone: os.Getenv("DB_TIMEZONE"),
	}

	return Config
}
