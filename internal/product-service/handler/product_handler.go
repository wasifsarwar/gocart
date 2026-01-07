package handler

import (
	"encoding/json"
	"fmt"
	"io"
	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	w.WriteHeader(http.StatusNoContent)
}

// UploadProductImage accepts multipart/form-data with field name "image".
// It stores the file under ./uploads/products and updates product.image_url to /uploads/products/<filename>.
func (h *ProductHandler) UploadProductImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ensure product exists
	product, err := h.repo.GetProductById(id)
	if err != nil {
		log.Printf("Error fetching product with id: %v and error: %v", id, err)
		http.Error(w, fmt.Sprintf("Product with id %v not found.", id), http.StatusNotFound)
		return
	}
	oldImageURL := product.ImageURL

	// Limit upload size (5MB)
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(6 << 20); err != nil {
		http.Error(w, "Invalid multipart form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Missing image file (field name: image)", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Infer extension (fallback to original filename ext)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
	default:
		http.Error(w, "Unsupported image type. Use jpg, png, webp, or gif.", http.StatusBadRequest)
		return
	}

	dir := filepath.Join("uploads", "products")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	filename := fmt.Sprintf("%s-%s%s", id, uuid.New().String(), ext)
	dstPath := filepath.Join(dir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Track whether upload completed successfully to clean up orphaned files
	uploadSuccess := false
	defer func() {
		if !uploadSuccess {
			// Remove the file if any operation after creation failed
			if removeErr := os.Remove(dstPath); removeErr != nil {
				log.Printf("Failed to clean up orphaned file %s: %v", dstPath, removeErr)
			}
		}
	}()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to write image", http.StatusInternalServerError)
		return
	}

	newImageURL := "/uploads/products/" + filename
	product.ImageURL = newImageURL
	updated, err := h.repo.UpdateProduct(product)
	if err != nil {
		http.Error(w, "Failed to update product image", http.StatusInternalServerError)
		return
	}

	// Mark upload as successful so defer doesn't remove the file
	uploadSuccess = true

	// Best-effort cleanup: delete the previous uploaded image (only if it was a local upload).
	// Don't delete external seed URLs (picsum) or non-local paths.
	if oldImageURL != "" && oldImageURL != newImageURL && strings.HasPrefix(oldImageURL, "/uploads/products/") {
		oldFilename := filepath.Base(oldImageURL)
		oldPath := filepath.Join("uploads", "products", oldFilename)
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to remove old product image %s: %v", oldPath, err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}
