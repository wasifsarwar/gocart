package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"gocart/internal/product-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockProductRepository struct {
	MockListAllProducts func() ([]models.Product, error)
	MockCreateProduct   func(product models.Product) (models.Product, error)
	MockGetProductById  func(id string) (models.Product, error)
	MockUpdateProduct   func(product models.Product) (models.Product, error)
	MockDeleteProduct   func(id string) error
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

func (m *MockProductRepository) UpdateProduct(product models.Product) (models.Product, error) {
	return m.MockUpdateProduct(product)
}

func (m *MockProductRepository) DeleteProduct(id string) error {
	return m.MockDeleteProduct(id)
}

func TestListProducts(t *testing.T) {
	tests := []struct {
		name           string
		mockProducts   []models.Product
		mockError      error
		expectedStatus int
	}{
		{
			name: "Success",
			mockProducts: []models.Product{
				{ProductID: "1", Name: "Test Product 1", Price: 99.99},
				{ProductID: "2", Name: "Test Product 2", Price: 199.99},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Database Error",
			mockProducts:   nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockProductRepository{
				MockListAllProducts: func() ([]models.Product, error) {
					return tt.mockProducts, tt.mockError
				},
			}

			handler := NewProductHandler(mockRepo)
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			w := httptest.NewRecorder()

			handler.ListProducts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.mockError == nil {
				var response []models.Product
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if len(response) != len(tt.mockProducts) {
					t.Errorf("Expected %d products, got %d", len(tt.mockProducts), len(response))
				}
			}
		})
	}
}

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		input          models.Product
		mockProduct    models.Product
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			input:          models.Product{Name: "New Product", Price: 99.99},
			mockProduct:    models.Product{ProductID: "1", Name: "New Product", Price: 99.99},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Database Error",
			input:          models.Product{Name: "New Product", Price: 99.99},
			mockProduct:    models.Product{},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockProductRepository{
				MockCreateProduct: func(product models.Product) (models.Product, error) {
					return tt.mockProduct, tt.mockError
				},
			}

			handler := NewProductHandler(mockRepo)
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.mockError == nil {
				var response models.Product
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.ProductID != tt.mockProduct.ProductID {
					t.Errorf("Expected product ID %s, got %s", tt.mockProduct.ProductID, response.ProductID)
				}
			}
		})
	}
}

func TestGetProductById(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockProduct    models.Product
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			productID:      "1",
			mockProduct:    models.Product{ProductID: "1", Name: "Test Product", Price: 99.99},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			productID:      "999",
			mockProduct:    models.Product{},
			mockError:      errors.New("product not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockProductRepository{
				MockGetProductById: func(id string) (models.Product, error) {
					return tt.mockProduct, tt.mockError
				},
			}

			handler := NewProductHandler(mockRepo)
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			w := httptest.NewRecorder()

			handler.GetProductById(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.mockError == nil {
				var response models.Product
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.ProductID != tt.mockProduct.ProductID {
					t.Errorf("Expected product ID %s, got %s", tt.mockProduct.ProductID, response.ProductID)
				}
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		input          models.Product
		mockProduct    models.Product
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			productID:      "1",
			input:          models.Product{ProductID: "1", Name: "Updated Product", Price: 199.99},
			mockProduct:    models.Product{ProductID: "1", Name: "Updated Product", Price: 199.99},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			productID:      "999",
			input:          models.Product{ProductID: "999", Name: "Updated Product", Price: 199.99},
			mockProduct:    models.Product{},
			mockError:      errors.New("product not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockProductRepository{
				MockGetProductById: func(id string) (models.Product, error) {
					if tt.mockError != nil {
						return models.Product{}, tt.mockError
					}
					return tt.mockProduct, nil
				},
				MockUpdateProduct: func(product models.Product) (models.Product, error) {
					return tt.mockProduct, tt.mockError
				},
			}

			handler := NewProductHandler(mockRepo)
			body, _ := json.Marshal(tt.input)
			req := httptest.NewRequest(http.MethodPut, "/products/"+tt.productID, bytes.NewBuffer(body))
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			w := httptest.NewRecorder()

			handler.UpdateProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.mockError == nil {
				var response models.Product
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.ProductID != tt.mockProduct.ProductID {
					t.Errorf("Expected product ID %s, got %s", tt.mockProduct.ProductID, response.ProductID)
				}
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			productID:      "1",
			mockError:      nil,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "Not Found",
			productID:      "999",
			mockError:      errors.New("product not found"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockProductRepository{
				MockDeleteProduct: func(id string) error {
					return tt.mockError
				},
			}

			handler := NewProductHandler(mockRepo)
			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.productID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.productID})
			w := httptest.NewRecorder()

			handler.DeleteProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
