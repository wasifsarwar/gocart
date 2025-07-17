package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	productHandler "gocart/internal/product-service/handler"
	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	productServer "gocart/internal/product-service/server"
	userHandler "gocart/internal/user-service/handler"
	userModels "gocart/internal/user-service/models"
	userRepository "gocart/internal/user-service/repository"
	userServer "gocart/internal/user-service/server"
	db "gocart/pkg/db"
	"gocart/pkg/seeder"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting GoCart E-commerce API...")

	// Initialize database connection
	db.Connect(db.DefaultConfig())
	db.Migrate(&productModels.Product{})
	db.Migrate(&userModels.User{})

	// Initialize repositories
	productRepo := productRepository.NewProductRepository(db.DB)
	userRepo := userRepository.NewUserRepository(db.DB)

	// Initialize handlers
	productHandler := productHandler.NewProductHandler(productRepo)
	userHandler := userHandler.NewUserHandler(userRepo)

	// Initialize servers
	productSrv := productServer.NewServer(productHandler)
	userSrv := userServer.NewServer(userHandler)

	// Seed database with sample data
	seederInstance := seeder.NewSeeder(productRepo, userRepo)
	if err := seederInstance.SeedAll(); err != nil {
		log.Printf("⚠️  Warning: Failed to seed database: %v", err)
	} else {
		// Print seeding summary
		seederInstance.PrintSeedingSummary()
	}

	// Start servers concurrently on different ports
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Printf("🛍️  Product Service starting on http://localhost:8080")
		log.Printf("📖 Product API docs: http://localhost:8080/products")
		if err := productSrv.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Product server failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		log.Printf("👥 User Service starting on http://localhost:8081")
		log.Printf("📖 User API docs: http://localhost:8081/users")
		if err := userSrv.Start(":8081"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("User server failed: %v", err)
		}
	}()

	log.Println("🚀 All services started successfully!")
	log.Println("📚 Full API Documentation: https://wasifsarwar.github.io/gocart/")

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down servers...")

	wg.Wait()
	log.Println("✅ Servers stopped gracefully")
}
