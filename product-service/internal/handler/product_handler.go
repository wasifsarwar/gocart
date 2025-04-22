package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"product-service/internal/models"
	"product-service/internal/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func ListAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := repository.ListAllProducts()
	if err != nil {
		log.Printf("Error fetching products with error: %v", err) // Log the error with context
		http.Error(w, "Unable to retrieve products. Please try again later.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate and assign a new UUID
	product.ID = uuid.New().String()

	newProduct, err := repository.CreateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // Set status to 201 Created

	// Optionally, include the product header
	w.Header().Set("Location", "/products/"+newProduct.ID)

	// Return the created product
	json.NewEncoder(w).Encode(newProduct)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := repository.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err) // Log the error with context
		http.Error(w, fmt.Sprintf("Unable to retrieve product with id: %v. Please try again later.", id), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// First get the existing product to update
	existingProduct, err := repository.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Update fields
	updatedProduct.ID = existingProduct.ID

	// Call repository function
	result, err := repository.UpdateProduct(updatedProduct)
	if err != nil {
		log.Printf("Error updating product with id: %v and error: %v", id, err)
		http.Error(w, "Unable to update product", http.StatusInternalServerError)
		return
	}

	// Set the response header and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // Set status to 200 Created

	// Return the updated product
	json.NewEncoder(w).Encode(result)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// First get the product to delete
	product, err := repository.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	// Delete the product
	_, err = repository.DeleteProduct(product)
	if err != nil {
		log.Printf("Error deleting product with id: %v and error: %v", id, err)
		http.Error(w, "Unable to delete product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
