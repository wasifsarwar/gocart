package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	productHandler "gocart/internal/product-service/handler"
	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	productServer "gocart/internal/product-service/server"
	userHandler "gocart/internal/user-service/handler"
	userModels "gocart/internal/user-service/models"
	userRepository "gocart/internal/user-service/repository"
	userServer "gocart/internal/user-service/server"
	db "gocart/pkg/db"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)

	// Initialize database connection
	db.Connect(db.DefaultConfig())
	db.Migrate(&productModels.Product{})
	db.Migrate(&userModels.User{})

	// Initialize dependencies
	productRepo := productRepository.NewProductRepository(db.DB)
	productHandler := productHandler.NewProductHandler(productRepo)
	userRepo := userRepository.NewUserRepository(db.DB)
	userHandler := userHandler.NewUserHandler(userRepo)

	// Initialize servers
	productSrv := productServer.NewServer(productHandler)
	userSrv := userServer.NewServer(userHandler)

	//Test grab current list of products
	allProducts, err := productRepo.ListAllProducts()
	if err != nil {
		log.Printf("Error listing existing products %v", err)
	} else {
		log.Printf("Here's a list of products")
		for _, product := range allProducts {
			log.Printf("Product %s (ID: %s) - Price: %.2f, Category: %s",
				product.Name,
				product.ProductID,
				product.Price,
				product.Category)
		}

	}

	// Test data insertion
	inputProduct := productModels.Product{
		ProductID:   fmt.Sprintf("test-%d", time.Now().Unix()),
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    "Test Category",
	}

	// Use repository for data operations
	createdProduct, err := productRepo.CreateProduct(inputProduct)
	if err != nil {
		log.Printf("Error creating product: %v", err)
	} else {
		log.Printf("Successfully created product with ID: %s", createdProduct.ProductID)

		// Retrieve the product to verify
		product, err := productRepo.GetProductById(createdProduct.ProductID)
		if err != nil {
			log.Printf("Error getting product: %v", err)
		} else {
			log.Printf("Retrieved product: %+v", product)
		}
	}

	// Test User data insertion
	inputUser := userModels.User{
		UserID:    fmt.Sprintf("test-%d", time.Now().Unix()),
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		Phone:     "+1234567890",
	}

	createdUser, err := userRepo.CreateUser(inputUser)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	} else {

		log.Printf("Successfully created user with ID: %s", createdUser.UserID)

		// Retrieve the user to verify
		user, err := userRepo.GetUserById(createdUser.UserID)
		if err != nil {
			log.Printf("Error getting user: %v", err)
		} else {
			log.Printf("Retrieved user: %+v", user)
		}
	}

	// Start servers concurrently on different ports
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := productSrv.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Product server failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := userSrv.Start(":8081"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("User server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Shutdown both (basic; extend Start with ctx if needed)
	wg.Wait()
	log.Println("Servers stopped")
}
