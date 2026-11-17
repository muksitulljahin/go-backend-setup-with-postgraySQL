package main

import (
	"log"

	imagecollection "github.com/go-backend-setup-with-postgray/api/v1/imageCollection"
	"github.com/go-backend-setup-with-postgray/config"
)

func main() {
	cfg := config.LoadConfig()
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}

	log.Println("🚀 Running database migration...")
	err = db.AutoMigrate(
		&imagecollection.ImageCollection{},
	)
	if err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	log.Println("✅ Database migration completed successfully!")
}
