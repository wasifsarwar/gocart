package seeder

import (
	"fmt"
	"log"
	"time"

	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	userModels "gocart/internal/user-service/models"
	userRepository "gocart/internal/user-service/repository"

	"github.com/google/uuid"
)

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

// SeedProducts creates sample products
func (s *SeedData) SeedProducts() error {
	log.Println("Seeding sample products...")

	// Check if products already exist
	existingProducts, err := s.ProductRepo.ListAllProducts()
	if err != nil {
		return fmt.Errorf("failed to check existing products: %w", err)
	}

	if len(existingProducts) > 0 {
		log.Printf("Found %d existing products, skipping product seeding", len(existingProducts))
		return nil
	}

	sampleProducts := []productModels.Product{
		// Electronics
		{
			ProductID:   uuid.New().String(),
			Name:        "iPhone 15 Pro",
			Description: "Latest Apple smartphone with titanium design, A17 Pro chip, and advanced camera system",
			Price:       999.99,
			Category:    "Electronics",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "MacBook Pro 16-inch",
			Description: "Powerful laptop with M3 Max chip, 32GB RAM, and stunning Liquid Retina XDR display",
			Price:       2499.99,
			Category:    "Electronics",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "AirPods Pro (2nd gen)",
			Description: "Premium wireless earbuds with active noise cancellation and spatial audio",
			Price:       249.99,
			Category:    "Electronics",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Samsung 4K Smart TV 65\"",
			Description: "Ultra HD Smart TV with HDR, built-in streaming apps, and voice control",
			Price:       899.99,
			Category:    "Electronics",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Sony WH-1000XM5 Headphones",
			Description: "Industry-leading noise canceling wireless headphones with 30-hour battery",
			Price:       399.99,
			Category:    "Electronics",
		},

		// Clothing
		{
			ProductID:   uuid.New().String(),
			Name:        "Levi's 501 Original Jeans",
			Description: "Classic straight-leg jeans in premium denim, available in multiple washes",
			Price:       89.99,
			Category:    "Clothing",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Nike Air Max 270",
			Description: "Comfortable running shoes with Max Air heel unit and breathable mesh upper",
			Price:       149.99,
			Category:    "Clothing",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Patagonia Down Jacket",
			Description: "Lightweight, packable down jacket perfect for outdoor adventures",
			Price:       229.99,
			Category:    "Clothing",
		},

		// Home & Garden
		{
			ProductID:   uuid.New().String(),
			Name:        "Dyson V15 Detect Vacuum",
			Description: "Cordless vacuum with laser dust detection and powerful suction",
			Price:       749.99,
			Category:    "Home & Garden",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Instant Pot Duo 7-in-1",
			Description: "Multi-use pressure cooker, slow cooker, rice cooker, and more",
			Price:       99.99,
			Category:    "Home & Garden",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Philips Hue Smart Bulbs (4-pack)",
			Description: "Color-changing smart LED bulbs controllable via smartphone app",
			Price:       199.99,
			Category:    "Home & Garden",
		},

		// Books
		{
			ProductID:   uuid.New().String(),
			Name:        "The Psychology of Money",
			Description: "Timeless lessons on wealth, greed, and happiness by Morgan Housel",
			Price:       16.99,
			Category:    "Books",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "Clean Code",
			Description: "A handbook of agile software craftsmanship by Robert C. Martin",
			Price:       42.99,
			Category:    "Books",
		},

		// Sports & Outdoors
		{
			ProductID:   uuid.New().String(),
			Name:        "Yeti Rambler 30 oz Tumbler",
			Description: "Stainless steel insulated tumbler that keeps drinks hot or cold for hours",
			Price:       39.99,
			Category:    "Sports & Outdoors",
		},
		{
			ProductID:   uuid.New().String(),
			Name:        "REI Co-op Trail 40 Backpack",
			Description: "Versatile hiking backpack with adjustable suspension and multiple pockets",
			Price:       139.99,
			Category:    "Sports & Outdoors",
		},
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

	log.Printf("Successfully seeded %d products", len(sampleProducts))
	return nil
}

// SeedUsers creates sample users
func (s *SeedData) SeedUsers() error {
	log.Println("Seeding sample users...")

	// Check if users already exist
	existingUsers, err := s.UserRepo.ListAllUsers()
	if err != nil {
		return fmt.Errorf("failed to check existing users: %w", err)
	}

	if len(existingUsers) > 0 {
		log.Printf("Found %d existing users, skipping user seeding", len(existingUsers))
		return nil
	}

	now := time.Now()
	sampleUsers := []userModels.User{
		{
			UserID:    uuid.New().String(),
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			Phone:     "+1-555-0102",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Michael",
			LastName:  "Johnson",
			Email:     "michael.johnson@example.com",
			Phone:     "+1-555-0103",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Emily",
			LastName:  "Brown",
			Email:     "emily.brown@example.com",
			Phone:     "+1-555-0104",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "David",
			LastName:  "Wilson",
			Email:     "david.wilson@example.com",
			Phone:     "+1-555-0105",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Sarah",
			LastName:  "Davis",
			Email:     "sarah.davis@example.com",
			Phone:     "+1-555-0106",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Robert",
			LastName:  "Miller",
			Email:     "robert.miller@example.com",
			Phone:     "+1-555-0107",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Lisa",
			LastName:  "Garcia",
			Email:     "lisa.garcia@example.com",
			Phone:     "+1-555-0108",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Christopher",
			LastName:  "Martinez",
			Email:     "chris.martinez@example.com",
			Phone:     "+1-555-0109",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Amanda",
			LastName:  "Anderson",
			Email:     "amanda.anderson@example.com",
			Phone:     "+1-555-0110",
			CreatedAt: now,
			UpdatedAt: now,
		},
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

	log.Printf("Successfully seeded %d users", len(sampleUsers))
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
