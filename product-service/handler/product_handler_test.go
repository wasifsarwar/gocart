package handler

import "gocart/product-service/models"

type MockProductRepository struct {
	MockListAllProducts func() ([]models.Product, error)
	MockCreateProduct   func(product models.Product) (models.Product, error)
	MockGetProductById  func(id string) (models.Product, error)
	MockUpdateProduct   func(product models.Product) (models.Product, error)
	MockDeleteProduct   func(product models.Product)
}

func (m *MockProductRepository) ListAllProducts() ([]models.Product, error) {
	return m.MockListAllProducts()
}

func (m *MockProductRepository) CreateProduct(product models.Product) (models.Product, error) {
	return m.MockCreateProduct(product)
}

func (m *MockProductRepository) GetProductById(id string) (models.Product, error) {
	return m.MockGetProductById(id)
}
