package handler

import (
	"bytes"
	"encoding/json"
	"gocart/internal/product-service/models"
	"gocart/internal/product-service/repository"
	"gocart/pkg/testutils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	config := testutils.TestDBConfig{
		ServiceName: "products_handler",
		Models:      []interface{}{&models.Product{}},
	}
	return testutils.SetupTestDB(t, config)
}

func TestCreateAndGetProductIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestCreateAndGetProductIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting create a product and retrieve integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	productRepo := repository.NewProductRepository(db)
	handler := NewProductHandler(productRepo)
	testProduct := models.Product{
		Category:    "Test Category",
		Name:        "TestProduct",
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

	getProductID := createdProductResponse.ProductID

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

func TestListProductsIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestListProductsIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting list products integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	productRepo := repository.NewProductRepository(db)
	handler := NewProductHandler(productRepo)

	product1 := models.Product{
		Name:        "Product 1",
		Description: "Description 1",
		Price:       100.00,
		Category:    "Category 1",
	}

	product2 := models.Product{
		Name:        "Product 2",
		Description: "Description 2",
		Price:       200.00,
		Category:    "Category 2",
	}

	product3 := models.Product{
		Name:        "Product 3",
		Description: "Description 3",
		Price:       300.00,
		Category:    "Category 3",
	}

	productRepo.CreateProduct(product1)
	productRepo.CreateProduct(product2)
	productRepo.CreateProduct(product3)

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

	if len(response) != 3 {
		t.Errorf("Expected 3 products, got %d", len(response))
	}
}
