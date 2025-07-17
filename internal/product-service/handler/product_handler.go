package handler

import (
	"encoding/json"
	"fmt"
	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductHandler struct {
	repo productRepository.ProductRepository
}

func NewProductHandler(repo productRepository.ProductRepository) *ProductHandler {
	return &ProductHandler{
		repo: repo,
	}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.ListAllProducts()
	if err != nil {
		log.Printf("Error fetching products with error: %v", err)
		http.Error(w, "Unable to retrieve products. Please try again later.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product productModels.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product.ProductID = uuid.New().String()

	newProduct, err := h.repo.CreateProduct(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/products/"+newProduct.ProductID)
	json.NewEncoder(w).Encode(newProduct)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	product, err := h.repo.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		if err.Error() == "product not found" {
			http.Error(w, fmt.Sprintf("Product with id %v not found.", id), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Unable to retrieve product with id: %v.", id), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedProduct productModels.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	existingProduct, err := h.repo.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	updatedProduct.ProductID = existingProduct.ProductID

	result, err := h.repo.UpdateProduct(updatedProduct)
	if err != nil {
		log.Printf("Error updating product with id: %v and error: %v", id, err)
		http.Error(w, "Unable to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.repo.DeleteProduct(id)
	if err != nil {
		log.Printf("Error deleting product with id: %v and error: %v", id, err)
		http.Error(w, "Unable to delete product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
