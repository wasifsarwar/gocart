package repository

import (
	"product-service/internal/db"
	"product-service/internal/models"
)

type ProductRepository interface {
	CreateProduct(product models.Product) (models.Product, error)
	GetProductById(id string) (models.Product, error)
	UpdateProduct(product models.Product) (models.Product, error)
	DeleteProduct(product models.Product) (models.Product, error)
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

func ListAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := db.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func CreateProduct(product models.Product) (models.Product, error) {
	if err := db.DB.Create(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func GetProductById(id string) (models.Product, error) {
	var product models.Product
	if err := db.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func UpdateProduct(product models.Product) (models.Product, error) {
	if err := db.DB.Model(&models.Product{}).Where("id = ?", product.ID).Updates(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}

func DeleteProduct(product models.Product) (models.Product, error) {
	if err := db.DB.Where("id = ?", product.ID).Delete(&product).Error; err != nil {
		return models.Product{}, err
	}
	return product, nil
}
