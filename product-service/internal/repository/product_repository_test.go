package repository

import (
	"errors"
	"fmt"
	"os"
	"product-service/internal/models"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockProductRepository struct {
	MockCreateProduct  func(product models.Product) (models.Product, error)
	MockGetProductById func(id string) (models.Product, error)
}

func (m *MockProductRepository) CreateProduct(product models.Product) (models.Product, error) {
	return m.MockCreateProduct(product)
}

func (m *MockProductRepository) GetProductById(id string) (models.Product, error) {
	return m.MockGetProductById(id)
}

func TestCreateProduct(t *testing.T) {
	mockRepo := &MockProductRepository{
		MockCreateProduct: func(product models.Product) (models.Product, error) {
			return product, nil
		},
	}

	product := models.Product{
		ProductID:   "test-123",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
	}

	createdProduct, err := mockRepo.CreateProduct(product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if createdProduct.ProductID != product.ProductID {
		t.Errorf("Expected product ID to be %s, but got %s", product.ProductID, createdProduct.ProductID)
	}

	if createdProduct.Name != product.Name {
		t.Errorf("Expected product name to be %s, but got %s", product.Name, createdProduct.Name)
	}

	if createdProduct.Description != product.Description {
		t.Errorf("Expected product description to be %s, but got %s", product.Description, createdProduct.Description)
	}

	if createdProduct.Price != product.Price {
		t.Errorf("Expected product price to be %f, but got %f", product.Price, createdProduct.Price)
	}
}

func TestCreateProductError(t *testing.T) {
	mockRepo := &MockProductRepository{
		MockCreateProduct: func(product models.Product) (models.Product, error) {
			return models.Product{}, errors.New("failed to create product")
		},
	}

	product := models.Product{
		ProductID:   "test-123",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
	}

	_, err := mockRepo.CreateProduct(product)
	if err == nil {
		t.Errorf("Expected error to be returned, but got nil")
	}

	if err.Error() != "failed to create product" {
		t.Errorf("Expected error message to be 'failed to create product', but got '%s'", err.Error())
	}
}

func TestGetProductById(t *testing.T) {
	mockRepo := &MockProductRepository{
		MockGetProductById: func(id string) (models.Product, error) {
			return models.Product{ProductID: id, Name: "Test Product", Description: "Test Description", Price: 100.00}, nil
		},
	}

	product, err := mockRepo.GetProductById("test-123")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if product.ProductID != "test-123" {
		t.Errorf("Expected product ID to be 'test-123', but got '%s'", product.ProductID)
	}

	if product.Name != "Test Product" {
		t.Errorf("Expected product name to be 'Test Product', but got '%s'", product.Name)
	}

	if product.Description != "Test Description" {
		t.Errorf("Expected product description to be 'Test Description', but got '%s'", product.Description)
	}

	if product.Price != 100.00 {
		t.Errorf("Expected product price to be 100.00, but got %f", product.Price)
	}
}

func TestGetProductByIdError(t *testing.T) {
	mockRepo := &MockProductRepository{
		MockGetProductById: func(id string) (models.Product, error) {
			return models.Product{}, errors.New("product not found")
		},
	}

	_, err := mockRepo.GetProductById("test-123")
	if err == nil {
		t.Errorf("Expected error to be returned, but got nil")
	}
}

func TestDeleteProductIntegration(t *testing.T) {
	// Setup DB connection
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewProductRepository(db)

	// Create a test product first
	testProduct := models.Product{
		ProductID:   uuid.New().String(),
		Name:        "Test Product for Deletion",
		Description: "This product should be deleted",
		Price:       99.99,
	}

	// Insert the product
	createdProduct, err := repo.CreateProduct(testProduct)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}

	// Verify product was created
	fetchedProduct, err := repo.GetProductById(createdProduct.ProductID)
	if err != nil {
		t.Fatalf("Failed to fetch created product: %v", err)
	}
	if fetchedProduct.ProductID != createdProduct.ProductID {
		t.Errorf("Expected product ID %s, got %s", createdProduct.ProductID, fetchedProduct.ProductID)
	}

	// Delete the product
	err = repo.DeleteProduct(createdProduct.ProductID)
	if err != nil {
		t.Errorf("Failed to delete product: %v", err)
	}

	// Verify product was deleted
	_, err = repo.GetProductById(createdProduct.ProductID)
	if err == nil {
		t.Error("Expected error when fetching deleted product, got nil")
	}
}

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	dsn := getEnv()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Get the underlying *sql.DB instance and defer its closure
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying *sql.DB: %v", err)
	}
	cleanup := func() {
		sqlDB.Close()
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
