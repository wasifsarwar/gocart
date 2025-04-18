package repository

import (
	"errors"
	"product-service/internal/models"
	"testing"
)

type MockProductRepository struct {
	MockCreateProduct func(product models.Product) (models.Product, error)
}

func (m *MockProductRepository) CreateProduct(product models.Product) (models.Product, error) {
	return m.MockCreateProduct(product)
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
