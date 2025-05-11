package main

import (
	"fmt"
	db "gocart/pkg/db"
	"gocart/product-service/handler"
	"gocart/product-service/models"
	"gocart/product-service/repository"
	"gocart/product-service/server"
	"log"
	"time"
)

func main() {
	// Initialize database connection
	db.Connect(db.DefaultConfig())
	db.Migrate(&models.Product{})

	// Initialize dependencies
	productRepo := repository.NewProductRepository(db.DB)
	productHandler := handler.NewProductHandler(productRepo)
	srv := server.NewServer(productHandler)

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
	inputProduct := models.Product{
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

	// Start the server
	if err := srv.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}
