package repository

import (
	"gocart/internal/user-service/models"
	"gocart/pkg/testutils"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	config := testutils.TestDBConfig{
		ServiceName: "users_repo",
		Models:      []interface{}{&models.User{}},
	}
	return testutils.SetupTestDB(t, config)
}

func TestCreateAndGetUserIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestCreateAndGetUserIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting create and get user integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	testUser := models.User{
		UserID:    uuid.New().String(),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "123-456-7890",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	logger.Printf("Creating test user with ID: %s", testUser.UserID)
	createdUser, err := repo.CreateUser(testUser)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	logger.Printf("Successfully created test user: %+v", createdUser)

	logger.Printf("Fetching user by ID: %s", createdUser.UserID)
	fetchedUser, err := repo.GetUserById(createdUser.UserID)
	if err != nil {
		t.Fatalf("Failed to fetch user: %v", err)
	}

	if fetchedUser.UserID != createdUser.UserID {
		t.Errorf("Expected user ID %s, got %s", createdUser.UserID, fetchedUser.UserID)
	}
	if fetchedUser.Email != testUser.Email {
		t.Errorf("Expected email %s, got %s", testUser.Email, fetchedUser.Email)
	}

	logger.Println("Create and get user test completed successfully")
}

func TestListAllUsersIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestListAllUsersIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting list all users integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	// Create test users
	users := []models.User{
		{
			UserID:    uuid.New().String(),
			FirstName: "Alice",
			LastName:  "Smith",
			Email:     "alice@example.com",
			Phone:     "111-111-1111",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    uuid.New().String(),
			FirstName: "Bob",
			LastName:  "Jones",
			Email:     "bob@example.com",
			Phone:     "222-222-2222",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range users {
		_, err := repo.CreateUser(user)
		if err != nil {
			t.Fatalf("Failed to create user %s: %v", user.Email, err)
		}
		logger.Printf("Created user: %s", user.Email)
	}

	allUsers, err := repo.ListAllUsers()
	if err != nil {
		t.Fatalf("Failed to list all users: %v", err)
	}

	if len(allUsers) != len(users) {
		t.Errorf("Expected %d users, got %d", len(users), len(allUsers))
	}

	logger.Printf("List all users test completed successfully with %d users", len(allUsers))
}
