package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/config"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/internal/modules/user"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/internal/routes"
)

// @title           Go Backend with PostgreSQL & GORM API
// @version         1.0
// @description     Production-ready RESTful API server built with Go, Gin, PostgreSQL, and GORM.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.email   support@example.com

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection (PostgreSQL with GORM)
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Printf("⚠️ Warning: PostgreSQL database connection failed: %v\n", err)
		log.Println("Server will still start, but database-dependent endpoints will fail until PostgreSQL is configured properly in .env")
	} else {
		// Auto Migration for registered models
		log.Println("Running AutoMigration for models...")
		if err := db.AutoMigrate(&user.User{}); err != nil {
			log.Printf("⚠️ Auto migration warning: %v\n", err)
		} else {
			log.Println("✅ Database migration completed successfully")
		}
	}

	// 3. Setup Routes
	r := routes.SetupRouter(db)

	// 4. Configure HTTP Server
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start Server in a Goroutine
	go func() {
		log.Printf("🚀 Server is running on port %s (http://localhost:%s/)\n", cfg.Port, cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("❌ Server failed to start: %v\n", err)
		}
	}()

	// 6. Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited cleanly")
}
