package repository

import (
	"errors"
	"gocart/product-service/internal/models"
	"testing"
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
		ID:          "test-123",
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.00,
	}

	createdProduct, err := mockRepo.CreateProduct(product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	if createdProduct.ID != product.ID {
		t.Errorf("Expected product ID to be %s, but got %s", product.ID, createdProduct.ID)
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
		ID:          "test-123",
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
			return models.Product{ID: id, Name: "Test Product", Description: "Test Description", Price: 100.00}, nil
		},
	}

	product, err := mockRepo.GetProductById("test-123")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if product.ID != "test-123" {
		t.Errorf("Expected product ID to be 'test-123', but got '%s'", product.ID)
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
