package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gocart/product-service/models"
	"gocart/product-service/repository"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {

	t.Logf("Setting up test database connection at %v", time.Now())

	dsn := getEnv()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	t.Log("Running database migrations...")
	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Get the underlying *sql.DB instance and defer its closure
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying *sql.DB: %v", err)
	}
	cleanup := func() {
		t.Log("Cleaning up test database...")
		db.Migrator().DropTable(&models.Product{})
		sqlDB.Close()
		t.Log("Cleanup completed")

	}
	return db, cleanup
}

func getEnv() string {
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	user := os.Getenv("TEST_DB_USER")
	if user == "" {
		user = "admin"
	}

	password := os.Getenv("TEST_DB_PASSWORD")
	if password == "" {
		password = "admin"
	}

	dbname := os.Getenv("TEST_DB_NAME")
	if dbname == "" {
		dbname = "gocart_db"
	}

	port := os.Getenv("TEST_DB_PORT")
	if port == "" {
		port = "5432"
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	return dsn
}

func TestCreateAndGetProductIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestCreateAndGetProductIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting create a product and retrieve integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	productRepo := repository.NewProductRepository(db)
	handler := NewProductHandler(productRepo)
	productUUID := uuid.New().String()
	testProduct := models.Product{
		Name:        fmt.Sprintf("TestProduct-%s", productUUID),
		ProductID:   productUUID,
		Description: "Test product description",
		Price:       90.99,
	}

	body, err := json.Marshal(testProduct)
	if err != nil {
		t.Fatalf("Failed to marshal test product: %v", err)

	}

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateProduct(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	var createdProductResponse models.Product
	if err := json.Unmarshal(w.Body.Bytes(), &createdProductResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdProductResponse.ProductID == "" {
		t.Error("Expected ProductID to be set, got empty string")
	}
	if createdProductResponse.Name != testProduct.Name {
		t.Errorf("Expected product name %s, got %s", testProduct.Name, createdProductResponse.Name)
	}

	getProductID := testProduct.ProductID

	getReq := httptest.NewRequest(http.MethodGet, "/products/"+getProductID, nil)
	getReq = mux.SetURLVars(getReq, map[string]string{"id": getProductID}) // Set Mux vars
	getW := httptest.NewRecorder()

	handler.GetProductById(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected status code %d for GET request, got %d. Response body: %s", http.StatusOK, getW.Code, getW.Body.String())
	}

	var getResponse models.Product
	if err := json.Unmarshal(getW.Body.Bytes(), &getResponse); err != nil {
		t.Fatalf("Failed to unmarshal get response: %v. Body: %s", err, getW.Body.String())
	}

	if getResponse.ProductID != getProductID {
		t.Errorf("Expected fetched product ID %s, got %s", getProductID, getResponse.ProductID)
	}
	if getResponse.Name != testProduct.Name {
		t.Errorf("Expected product name %s, got %s", testProduct.Name, getResponse.Name)
	}
}

/**
d

func TestProductHandlerIntegration(t *testing.T) {
	setupTestDB(t)

	productRepo := repository.NewProductRepository(db.DB)
	handler := NewProductHandler(productRepo)

	// Test Create Product
	t.Run("Create Product", func(t *testing.T) {
		product := models.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			Category:    "Test Category",
		}

		body, err := json.Marshal(product)
		if err != nil {
			t.Fatalf("Failed to marshal product: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.CreateProduct(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var response models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.ProductID == "" {
			t.Error("Expected ProductID to be set, got empty string")
		}
		if response.Name != product.Name {
			t.Errorf("Expected product name %s, got %s", product.Name, response.Name)
		}
	})

	// Test Get Product
	t.Run("Get Product", func(t *testing.T) {
		product := models.Product{
			Name:        "Get Test Product",
			Description: "Test Description",
			Price:       99.99,
			Category:    "Test Category",
		}
		createdProduct, err := productRepo.CreateProduct(product)
		if err != nil {
			t.Fatalf("Failed to create test product: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/products/"+createdProduct.ProductID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": createdProduct.ProductID})
		w := httptest.NewRecorder()

		handler.GetProductById(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.ProductID != createdProduct.ProductID {
			t.Errorf("Expected ProductID %s, got %s", createdProduct.ProductID, response.ProductID)
		}
		if response.Name != createdProduct.Name {
			t.Errorf("Expected Name %s, got %s", createdProduct.Name, response.Name)
		}
	})

	// Test List Products
	t.Run("List Products", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		handler.ListProducts(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response []models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if len(response) == 0 {
			t.Error("Expected non-empty product list, got empty list")
		}
	})

	// Test Update Product
	t.Run("Update Product", func(t *testing.T) {
		product := models.Product{
			Name:        "Update Test Product",
			Description: "Test Description",
			Price:       99.99,
			Category:    "Test Category",
		}
		createdProduct, err := productRepo.CreateProduct(product)
		if err != nil {
			t.Fatalf("Failed to create test product: %v", err)
		}

		updatedProduct := createdProduct
		updatedProduct.Name = "Updated Name"
		updatedProduct.Price = 199.99

		body, err := json.Marshal(updatedProduct)
		if err != nil {
			t.Fatalf("Failed to marshal updated product: %v", err)
		}

		req := httptest.NewRequest(http.MethodPut, "/products/"+createdProduct.ProductID, bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": createdProduct.ProductID})
		w := httptest.NewRecorder()

		handler.UpdateProduct(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response models.Product
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response.Name != updatedProduct.Name {
			t.Errorf("Expected Name %s, got %s", updatedProduct.Name, response.Name)
		}
		if response.Price != updatedProduct.Price {
			t.Errorf("Expected Price %.2f, got %.2f", updatedProduct.Price, response.Price)
		}
	})

	// Test Delete Product
	t.Run("Delete Product", func(t *testing.T) {
		product := models.Product{
			Name:        "Delete Test Product",
			Description: "Test Description",
			Price:       99.99,
			Category:    "Test Category",
		}
		createdProduct, err := productRepo.CreateProduct(product)
		if err != nil {
			t.Fatalf("Failed to create test product: %v", err)
		}

		req := httptest.NewRequest(http.MethodDelete, "/products/"+createdProduct.ProductID, nil)
		req = mux.SetURLVars(req, map[string]string{"id": createdProduct.ProductID})
		w := httptest.NewRecorder()

		handler.DeleteProduct(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Verify product is deleted
		_, err = productRepo.GetProductById(createdProduct.ProductID)
		if err == nil {
			t.Error("Expected error when getting deleted product, got nil")
		}
	})
}
*/
