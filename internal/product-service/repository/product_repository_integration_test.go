package repository

import (
	"fmt"
	"gocart/internal/product-service/models"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestListAllProducts(t *testing.T) {
	logger := log.New(os.Stdout, "[TestListAllProducts]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting list all products integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewProductRepository(db)

	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("product-%v", i)
		product := models.Product{
			ProductID:   uuid.New().String(),
			Name:        name,
			Description: fmt.Sprintf("product description for %v", name),
			Price:       float64(i) * 29.99,
		}
		logger.Printf("Creating test product: %s", product.Name)
		_, err := repo.CreateProduct(product)
		if err != nil {
			t.Fatalf("Failed to create product %s: %v", product.Name, err)
		}
	}

	logger.Printf("Listing all products currently in db")
	products, err := repo.ListAllProducts()
	if err != nil {
		t.Fatalf("Failed to list all products: %v", err)
	}

	for _, product := range products {
		logger.Printf("Retrieved product %v with product id: %v", product.Name, product.ProductID)
	}
	logger.Printf("List all products test complicated successfully")

}

func TestUpdateProductIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestUpdateProductIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting update product integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewProductRepository(db)

	testProduct := models.Product{
		ProductID:   uuid.New().String(),
		Name:        "Test Product for Create Initially",
		Description: "This product should be created first",
		Price:       99.99,
	}

	logger.Printf("Creating test product with ID: %s", testProduct.ProductID)
	createdProduct, err := repo.CreateProduct(testProduct)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
	logger.Printf("Successfully created test product: %+v", createdProduct)

	updatedProduct := models.Product{
		ProductID:   testProduct.ProductID,
		Name:        "Updated Product Name",
		Description: "This product has been updated",
		Price:       109.99,
	}

	logger.Printf("Updating test product with ID: %s", updatedProduct.ProductID)
	updated, err := repo.UpdateProduct(updatedProduct)
	if err != nil {
		t.Fatalf("Failed to update test product: %v", err)
	}
	logger.Printf("Successfully updated test product: %+v", updated)

}

func TestCreateAndDeleteProductIntegration(t *testing.T) {

	logger := log.New(os.Stdout, "[TestCreateAndDeleteProductIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting delete product integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewProductRepository(db)

	testProduct := models.Product{
		ProductID:   uuid.New().String(),
		Name:        "Test Product for Deletion",
		Description: "This product should be deleted",
		Price:       99.99,
	}

	logger.Printf("Creating test product with ID: %s", testProduct.ProductID)
	createdProduct, err := repo.CreateProduct(testProduct)
	if err != nil {
		t.Fatalf("Failed to create test product: %v", err)
	}
	logger.Printf("Successfully created test product: %+v", createdProduct)

	logger.Printf("Verifying product creation...")
	fetchedProduct, err := repo.GetProductById(createdProduct.ProductID)
	if err != nil {
		t.Fatalf("Failed to fetch created product: %v", err)
	}
	if fetchedProduct.ProductID != createdProduct.ProductID {
		t.Errorf("Expected product ID %s, got %s", createdProduct.ProductID, fetchedProduct.ProductID)
	}
	logger.Printf("Successfully fetched product: %+v", fetchedProduct)

	logger.Printf("Attempting to delete product with ID: %s", createdProduct.ProductID)
	err = repo.DeleteProduct(createdProduct.ProductID)
	if err != nil {
		t.Errorf("Failed to delete product: %v", err)
	}
	logger.Println("Successfully deleted product")

	logger.Printf("Verifying product deletion...")
	_, err = repo.GetProductById(createdProduct.ProductID)
	if err == nil {
		t.Error("Expected error when fetching deleted product, got nil")
	} else {
		logger.Printf("Verified product deletion: %v", err)
	}

	logger.Println("Create and delete a product test completed successfully")
}
