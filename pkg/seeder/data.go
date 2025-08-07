package seeder

import (
	"fmt"
	"log"
	"os"
	"time"

	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	userModels "gocart/internal/user-service/models"
	userRepository "gocart/internal/user-service/repository"

	"gopkg.in/yaml.v3"
)

// YAML data structures
type ProductData struct {
	ProductID   string  `yaml:"product_id"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Price       float64 `yaml:"price"`
	Category    string  `yaml:"category"`
}

type UserData struct {
	UserID    string `yaml:"user_id"`
	FirstName string `yaml:"first_name"`
	LastName  string `yaml:"last_name"`
	Email     string `yaml:"email"`
	Phone     string `yaml:"phone"`
}

type ProductsYAML struct {
	Products []ProductData `yaml:"products"`
}

type UsersYAML struct {
	Users []UserData `yaml:"users"`
}

// SeedData contains all the sample data for seeding
type SeedData struct {
	ProductRepo productRepository.ProductRepository
	UserRepo    userRepository.UserRepository
}

// NewSeeder creates a new seeder instance
func NewSeeder(productRepo productRepository.ProductRepository, userRepo userRepository.UserRepository) *SeedData {
	return &SeedData{
		ProductRepo: productRepo,
		UserRepo:    userRepo,
	}
}

// SeedAll seeds both products and users
func (s *SeedData) SeedAll() error {
	log.Println("Starting database seeding...")

	if err := s.SeedProducts(); err != nil {
		return fmt.Errorf("failed to seed products: %w", err)
	}

	if err := s.SeedUsers(); err != nil {
		return fmt.Errorf("failed to seed users: %w", err)
	}

	log.Println("Database seeding completed successfully!")
	return nil
}

// loadProductsFromYAML loads product data from YAML file
func (s *SeedData) loadProductsFromYAML() ([]productModels.Product, error) {
	data, err := os.ReadFile("pkg/seeder/data/products.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read products.yaml: %w", err)
	}

	var productsYAML ProductsYAML
	if err := yaml.Unmarshal(data, &productsYAML); err != nil {
		return nil, fmt.Errorf("failed to unmarshal products.yaml: %w", err)
	}

	// Convert YAML data to product models
	var products []productModels.Product
	for _, productData := range productsYAML.Products {
		product := productModels.Product{
			ProductID:   productData.ProductID,
			Name:        productData.Name,
			Description: productData.Description,
			Price:       productData.Price,
			Category:    productData.Category,
		}
		products = append(products, product)
	}

	return products, nil
}

// SeedProducts creates sample products from YAML file
func (s *SeedData) SeedProducts() error {
	log.Println("Seeding sample products from YAML...")

	// Check if products already exist
	existingProducts, err := s.ProductRepo.ListAllProducts()
	if err != nil {
		return fmt.Errorf("failed to check existing products: %w", err)
	}

	if len(existingProducts) > 0 {
		log.Printf("Found %d existing products, skipping product seeding", len(existingProducts))
		return nil
	}

	// Load products from YAML file
	sampleProducts, err := s.loadProductsFromYAML()
	if err != nil {
		return fmt.Errorf("failed to load products from YAML: %w", err)
	}

	// Create all products
	for i, product := range sampleProducts {
		createdProduct, err := s.ProductRepo.CreateProduct(product)
		if err != nil {
			log.Printf("Failed to create product %s: %v", product.Name, err)
			continue
		}
		log.Printf("Created product %d: %s (ID: %s) - $%.2f",
			i+1, createdProduct.Name, createdProduct.ProductID, createdProduct.Price)
	}

	log.Printf("Successfully seeded %d products from YAML", len(sampleProducts))
	return nil
}

// loadUsersFromYAML loads user data from YAML file
func (s *SeedData) loadUsersFromYAML() ([]userModels.User, error) {
	data, err := os.ReadFile("pkg/seeder/data/users.yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to read users.yaml: %w", err)
	}

	var usersYAML UsersYAML
	if err := yaml.Unmarshal(data, &usersYAML); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users.yaml: %w", err)
	}

	// Convert YAML data to user models
	now := time.Now()
	var users []userModels.User
	for _, userData := range usersYAML.Users {
		user := userModels.User{
			UserID:    userData.UserID,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Email:     userData.Email,
			Phone:     userData.Phone,
			CreatedAt: now,
			UpdatedAt: now,
		}
		users = append(users, user)
	}

	return users, nil
}

// SeedUsers creates sample users from YAML file
func (s *SeedData) SeedUsers() error {
	log.Println("Seeding sample users from YAML...")

	// Check if users already exist
	existingUsers, err := s.UserRepo.ListAllUsers()
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}

	if len(existingUsers) > 0 {
		log.Printf("Found %d existing users, skipping user seeding", len(existingUsers))
		return nil
	}

	// Load users from YAML file
	sampleUsers, err := s.loadUsersFromYAML()
	if err != nil {
		return fmt.Errorf("failed to load users from YAML: %w", err)
	}

	// Create all users
	for i, user := range sampleUsers {
		createdUser, err := s.UserRepo.CreateUser(user)
		if err != nil {
			log.Printf("Failed to create user %s %s: %v", user.FirstName, user.LastName, err)
			continue
		}
		log.Printf("Created user %d: %s %s (ID: %s) - %s",
			i+1, createdUser.FirstName, createdUser.LastName, createdUser.UserID, createdUser.Email)
	}

	log.Printf("Successfully seeded %d users from YAML", len(sampleUsers))
	return nil
}

// GetSampleProductsByCategory returns products grouped by category for display
func (s *SeedData) GetSampleProductsByCategory() (map[string][]productModels.Product, error) {
	products, err := s.ProductRepo.ListAllProducts()
	if err != nil {
		return nil, err
	}

	categorized := make(map[string][]productModels.Product)
	for _, product := range products {
		categorized[product.Category] = append(categorized[product.Category], product)
	}

	return categorized, nil
}

// PrintSeedingSummary prints a nice summary of seeded data
func (s *SeedData) PrintSeedingSummary() {
	log.Println("\n=== SEEDING SUMMARY ===")

	// Product summary
	products, err := s.ProductRepo.ListAllProducts()
	if err == nil {
		log.Printf("Total Products: %d", len(products))

		// Group by category
		categories := make(map[string]int)
		for _, product := range products {
			categories[product.Category]++
		}

		for category, count := range categories {
			log.Printf("   • %s: %d products", category, count)
		}
	}

	// User summary
	users, err := s.UserRepo.ListAllUsers()
	if err == nil {
		log.Printf("Total Users: %d", len(users))
	}
}
