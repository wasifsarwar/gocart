package main

import (
	"log"
	"os"

	orderModels "gocart/internal/order-management-service/models"
	productModels "gocart/internal/product-service/models"
	productRepository "gocart/internal/product-service/repository"
	userModels "gocart/internal/user-service/models"
	userRepository "gocart/internal/user-service/repository"
	"gocart/pkg/db"
	"gocart/pkg/seeder"
)

func main() {
	log.Println("🚜 Starting GoCart database seeder...")

	// Allow overriding connection details via env vars like DB_HOST, DB_PORT, etc.
	cfg := db.DefaultConfig()
	if _, err := db.Connect(cfg); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Ensure schemas exist before inserting seed data.
	if err := db.MigrateAll(
		db.DB,
		&productModels.Product{},
		&userModels.User{},
		&orderModels.Order{},
		&orderModels.OrderItem{},
	); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	productRepo := productRepository.NewProductRepository(db.DB)
	userRepo := userRepository.NewUserRepository(db.DB)

	seed := seeder.NewSeeder(productRepo, userRepo)
	if err := seed.SeedAll(); err != nil {
		log.Fatalf("failed to seed database: %v", err)
	}

	seed.PrintSeedingSummary()
	log.Println("✅ Database seeding completed successfully")

	// Exit explicitly to make it easy to use in scripts/CI pipelines.
	os.Exit(0)
}
