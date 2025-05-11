package repository

import (
	"gocart/product-service/models"

	"gorm.io/gorm"
)

type ProductRepository interface {
	ListAllProducts() ([]models.Product, error)
	CreateProduct(product models.Product) (models.Product, error)
	GetProductById(id string) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
	DeleteProduct(id string) error
}

/**
Error Handling:
	-GORM methods return a *gorm.DB object, which includes an Error field.
	- This field will be populated with any error that occurs during the execution of the database operation.
	- By checking err := db.DB.Method().Error, you can determine if the operation was successful or if an error occurred.
Chaining:
	- GORM allows method chaining, meaning you can call multiple methods in a single line.
	- The .Error field provides a way to check for errors after executing a chain of methods.

*/

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) ListAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) CreateProduct(product models.Product) (models.Product, error) {
	if err := r.db.Create(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *productRepository) GetProductById(productId string) (models.Product, error) {
	var product models.Product
	if err := r.db.Where("product_id = ?", productId).First(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *productRepository) UpdateProduct(product models.Product) (models.Product, error) {
	if err := r.db.Model(&models.Product{}).Where("product_id = ?", product.ProductID).Updates(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func (r *productRepository) DeleteProduct(id string) error {
	return r.db.Delete(&models.Product{}, "product_id = ?", id).Error
}
