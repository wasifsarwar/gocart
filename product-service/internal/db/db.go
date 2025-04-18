package db

import (
	"fmt"
	"log"
	"os"
	"product-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect establishes a connection to the PostgreSQL database
func Connect() {
	var err error

	// Get environment variables with defaults if not set
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "admin")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "gocart_db")
	port := getEnv("DB_PORT", "5432")

	// Construct the DSN string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	// Connect to the database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Printf("Connected to PostgreSQL database at %s:%s", host, port)
}

func Migrate() {
	err := DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	log.Println("Database migration completed successfully")
}

// GetFirstProduct retrieves the first product from the database
func GetFirstProduct() (models.Product, error) {
	var product models.Product
	err := DB.First(&product).Error
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

// CreateProduct adds a new product to the database
func CreateProduct(product models.Product) (models.Product, error) {
	err := DB.Create(&product).Error
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
