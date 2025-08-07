package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderHandler "gocart/internal/order-management-service/handler"
	orderModels "gocart/internal/order-management-service/models"
	orderRepository "gocart/internal/order-management-service/repository"
	orderServer "gocart/internal/order-management-service/server"
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

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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
	db.Migrate(&orderModels.Order{})
	db.Migrate(&orderModels.OrderItem{})

	// Initialize repositories
	productRepo := productRepository.NewProductRepository(db.DB)
	userRepo := userRepository.NewUserRepository(db.DB)
	orderRepo := orderRepository.NewOrderRepository(db.DB)

	// Initialize handlers
	productHandler := productHandler.NewProductHandler(productRepo)
	userHandler := userHandler.NewUserHandler(userRepo)
	orderHandler := orderHandler.NewOrderHandler(orderRepo)

	// Initialize servers
	productSrv := productServer.NewServer(productHandler)
	userSrv := userServer.NewServer(userHandler)
	orderSrv := orderServer.NewServer(orderHandler)

	// Seed database with sample data
	seederInstance := seeder.NewSeeder(productRepo, userRepo)
	if err := seederInstance.SeedAll(); err != nil {
		log.Printf("⚠️  Warning: Failed to seed database: %v", err)
	} else {
		// Print seeding summary
		seederInstance.PrintSeedingSummary()
	}

	mainRouter := mux.NewRouter()

	//Mount each service's router to the main router
	mainRouter.PathPrefix("/products").Handler(productSrv.GetRouter())
	mainRouter.PathPrefix("/users").Handler(userSrv.GetRouter())
	mainRouter.PathPrefix("/orders").Handler(orderSrv.GetRouter())

	// Health check endpoint
	mainRouter.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Add CORS middleware to main router
	corsRouter := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(mainRouter)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsRouter,
	}

	log.Printf("🚀 Starting unified GoCart API on http://localhost:%s", port)
	log.Printf("🛍️  Product API: http://localhost:%s/products", port)
	log.Printf("👥 User API: http://localhost:%s/users", port)
	log.Printf("�� Order API: http://localhost:%s/orders", port)

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	log.Println("📚 Full API Documentation: https://wasifsarwar.github.io/gocart/")

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down servers...")

	// Graceful shutdown
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
	log.Println("✅ Servers stopped gracefully")
}
