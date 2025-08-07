package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func DefaultConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "admin"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "gocart_db"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}
}

func Connect(config Config) (*gorm.DB, error) {

	// Log the connection attempt (hide password)
	log.Printf("Attempting to connect to database - Host: %s, Port: %s, User: %s, DB: %s, SSL: %s",
		config.Host, config.Port, config.User, config.DBName, config.SSLMode)

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	// Configure logger
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Set global DB connection
	DB = db
	log.Printf("Connected to PostgresSQL database at %s:%s", config.Host, config.Port)

	return db, nil
}

// Migrate runs database migrations
// Pass the models you want to migrate as parameters
func Migrate(models ...interface{}) {
	for _, model := range models {
		if err := DB.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate model: %v", err)
		}
	}
	log.Println("Database migration completed successfully")
}

// MigrateAll automatically migrates all service models
func MigrateAll(db *gorm.DB, models ...interface{}) error {
	if db == nil {
		db = DB
	}

	if db == nil {
		return errors.New("database connection not established")
	}

	log.Println("Starting database migration...")

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate model: %w", err)
		} else {
			log.Printf("Migrated model: %T", model)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
