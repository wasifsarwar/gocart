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

type ProductHandler struct {
	repo repository.ProductRepository
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
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
	var product models.Product
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
	id := vars["product_id"]

	product, err := h.repo.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		http.Error(w, fmt.Sprintf("Unable to retrieve product with id: %v.", id), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["product_id"]

	var updatedProduct models.Product
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
	id := vars["product_id"]

	product, err := h.repo.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	_, err = h.repo.DeleteProduct(product)
	if err != nil {
		log.Printf("Error deleting product with id: %v and error: %v", id, err)
		http.Error(w, "Unable to delete product", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
