package main

import (
	"log"
	"net/http"
	"product-service/internal/db"
	"product-service/internal/handler"
	"product-service/internal/models"
	"product-service/internal/repository"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	db.Connect()
	db.Migrate()

	// Create router
	r := mux.NewRouter()

	// Test data insertion
	inputProduct := models.Product{
		ID:          "test-123", // Set an explicit ID
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
		Category:    "Test Category",
	}

	// Use repository for data operations
	createdProduct, err := repository.CreateProduct(inputProduct)
	if err != nil {
		log.Printf("Error creating product: %v", err)
	} else {
		log.Printf("Successfully created product with ID: %s", createdProduct.ID)

		// Retrieve the product to verify
		product, err := repository.GetProductById(createdProduct.ID)
		if err != nil {
			log.Printf("Error getting product: %v", err)
		} else {
			log.Printf("Retrieved product: %+v", product)
		}
	}

	// Define routes
	r.HandleFunc("/products", handler.ListAllProducts).Methods("GET")
	r.HandleFunc("/products", handler.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", handler.GetProductById).Methods("GET")
	r.HandleFunc("/products/{id}", handler.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", handler.DeleteProduct).Methods("DELETE")

	// Start the server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
