package testutils

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDBConfig holds configuration for test database setup
type TestDBConfig struct {
	ServiceName string        // e.g., "products", "users"
	Models      []interface{} // Models to migrate
}

// SetupTestDB creates a unique test database and returns the connection and cleanup function
func SetupTestDB(t *testing.T, config TestDBConfig) (*gorm.DB, func()) {
	t.Logf("Setting up test database for %s service at %v", config.ServiceName, time.Now())

	// Create unique database name for this test
	testDBName := fmt.Sprintf("test_%s_%s_%d",
		config.ServiceName,
		uuid.New().String()[:8],
		time.Now().UnixNano())

	dsn := getEnv(testDBName)

	// First connect to the default database to create the test database
	defaultDSN := getEnv("postgres") // Connect to default postgres db
	adminDB, err := gorm.Open(postgres.Open(defaultDSN))
	if err != nil {
		t.Fatalf("Failed to connect to admin database: %v", err)
	}

	// Create the test database
	createDBSQL := fmt.Sprintf("CREATE DATABASE %s", testDBName)
	if err := adminDB.Exec(createDBSQL).Error; err != nil {
		// If database already exists, that's okay for cleanup scenarios
		t.Logf("Note: Database creation result: %v", err)
	}

	// Close admin connection
	adminSQLDB, _ := adminDB.DB()
	adminSQLDB.Close()

	// Connect to the test database
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Run migrations for all provided models
	if len(config.Models) > 0 {
		t.Log("Running database migrations...")
		err = db.AutoMigrate(config.Models...)
		if err != nil {
			t.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Get the underlying *sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("Failed to get underlying *sql.DB: %v", err)
	}

	cleanup := func() {
		t.Logf("Cleaning up test database for %s service...", config.ServiceName)

		// Close the connection first
		sqlDB.Close()

		// Connect to admin database to drop test database
		adminDB, err := gorm.Open(postgres.Open(getEnv("postgres")))
		if err != nil {
			t.Logf("Warning: Failed to connect for cleanup: %v", err)
			return
		}
		defer func() {
			if adminSQLDB, err := adminDB.DB(); err == nil {
				adminSQLDB.Close()
			}
		}()

		// Terminate any active connections to the test database
		terminateSQL := fmt.Sprintf(`
			SELECT pg_terminate_backend(pid) 
			FROM pg_stat_activity 
			WHERE datname = '%s' AND pid <> pg_backend_pid()`, testDBName)
		adminDB.Exec(terminateSQL)

		// Drop the test database
		dropDBSQL := fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBName)
		if err := adminDB.Exec(dropDBSQL).Error; err != nil {
			t.Logf("Warning: Failed to drop test database: %v", err)
		} else {
			t.Logf("Cleanup completed successfully for %s service", config.ServiceName)
		}
	}

	return db, cleanup
}

// getEnv constructs the database connection string from environment variables
func getEnv(customDBName ...string) string {
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
		password = "password"
	}

	dbname := "gocart_db"
	if len(customDBName) > 0 && customDBName[0] != "" {
		dbname = customDBName[0]
	}
	if envDBName := os.Getenv("TEST_DB_NAME"); envDBName != "" && len(customDBName) == 0 {
		dbname = envDBName
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
