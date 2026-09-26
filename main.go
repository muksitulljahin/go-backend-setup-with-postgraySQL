package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/api/routes"
	"github.com/muksitulljahin/go-backend-setup-with-postgraySQL/config"
)
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {
	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database Connection (PostgreSQL with GORM)
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Printf("❌ Database connection failed: %v\nServer will not start until database is ready.", err)
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
		localIP := getLocalIP()
		log.Println("--------------------------------------------------")
		log.Println("🚀 Server is running successfully!")
		log.Printf("🏡 Local:   http://localhost:%s/\n", cfg.Port)
		if localIP != "" && localIP != "127.0.0.1" {
			log.Printf("🌐 Network: http://%s:%s/\n", localIP, cfg.Port)
		}
		log.Printf("📚 Local Swagger: http://localhost:%s/swagger/index.html\n", cfg.Port)
		log.Printf("📚 Network Swagger: http://%s:%s/swagger/index.html\n", localIP, cfg.Port)
		log.Println("--------------------------------------------------")

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
