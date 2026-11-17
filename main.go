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

	"github.com/go-backend-setup-with-postgray/api/routes"
	"github.com/go-backend-setup-with-postgray/config"
)

// @title GO Lang Backend With GIn
// @version 1.0
// @description Project Backend Set
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token (e.g. "Bearer eyJhbGciOi...")

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection (PostgreSQL with GORM)
	db, err := config.ConnectDatabase(cfg)

	// 3. Setup Routes
	r := routes.SetupRouter(db, cfg)

	// 4. Configure HTTP Server
	port := cfg.Port

	serverAddr := fmt.Sprintf(":%s", port)
	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	localIP := config.GetLocalIP()

	go func() {
		log.Println("--------------------------------------------------")
		if err == nil {
			log.Println("🚀 Server is running successfully and Database is connected!")
			log.Printf("Local:   http://localhost:%s/\n", port)
			if localIP != "" {
				log.Printf("Network: http://%s:%s/\n", localIP, port)
			}
			log.Printf("Local Swagger:   http://localhost:%s/swagger/index.html\n", port)
			if localIP != "" {
				log.Printf("Network Swagger: http://%s:%s/swagger/index.html\n", localIP, port)
			}
		} else {
			log.Printf("Oops: Database connection failed: %v\n", err)
			log.Println("Server will not be accessible properly until database is ready. Please check your .env configuration.")
		}
		log.Println("--------------------------------------------------")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed to start: %v\n", err)
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
