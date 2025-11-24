package repository

import (
	"errors"
	"gocart/internal/user-service/models"
	"testing"
	"time"
)

type MockUserRepository struct {
	MockCreateUser     func(user models.User) (models.User, error)
	MockGetUserByID    func(id string) (models.User, error)
	MockGetUserByEmail func(email string) (models.User, error)
	MockUpdateUser     func(user models.User) (models.User, error)
	MockDeleteUser     func(id string) (models.User, error)
	MockListAllUsers   func() ([]models.User, error)
}

func (m *MockUserRepository) CreateUser(user models.User) (models.User, error) {
	return m.MockCreateUser(user)
}

func (m *MockUserRepository) GetUserById(id string) (models.User, error) {
	return m.MockGetUserByID(id)
}

func (m *MockUserRepository) GetUserByEmail(email string) (models.User, error) {
	return m.MockGetUserByEmail(email)
}

func (m *MockUserRepository) UpdateUser(user models.User) (models.User, error) {
	return m.MockUpdateUser(user)
}

func (m *MockUserRepository) DeleteUser(id string) (models.User, error) {
	return m.MockDeleteUser(id)
}

func (m *MockUserRepository) ListAllUsers() ([]models.User, error) {
	return m.MockListAllUsers()
}

// TestCreateUser tests the business logic for user creation
func TestCreateUser(t *testing.T) {
	t.Run("successful user creation with data transformation", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockCreateUser: func(user models.User) (models.User, error) {
				// Simulate business logic: generate UUID, set timestamps
				if user.UserID == "" {
					user.UserID = "generated-uuid-123"
				}
				user.CreatedAt = time.Now()
				user.UpdatedAt = time.Now()
				return user, nil
			},
		}

		user := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
		}

		createdUser, err := mockRepo.CreateUser(user)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Test business logic: UserID should be generated
		if createdUser.UserID == "" {
			t.Error("Expected UserID to be generated")
		}

		// Test business logic: timestamps should be set
		if createdUser.CreatedAt.IsZero() {
			t.Error("Expected CreatedAt to be set")
		}

		if createdUser.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}

		// Test data transformation: other fields should remain unchanged
		if createdUser.FirstName != user.FirstName {
			t.Errorf("Expected FirstName %s, got %s", user.FirstName, createdUser.FirstName)
		}

		if createdUser.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, createdUser.Email)
		}
	})

	t.Run("database error handling", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockCreateUser: func(user models.User) (models.User, error) {
				return models.User{}, errors.New("database connection failed")
			},
		}

		user := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
		}

		_, err := mockRepo.CreateUser(user)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "database connection failed" {
			t.Errorf("Expected error 'database connection failed', got '%s'", err.Error())
		}
	})

	t.Run("duplicate email validation", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockCreateUser: func(user models.User) (models.User, error) {
				// Simulate unique constraint violation
				if user.Email == "existing@example.com" {
					return models.User{}, errors.New("duplicate email address")
				}
				user.UserID = "new-uuid"
				user.CreatedAt = time.Now()
				return user, nil
			},
		}

		user := models.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "existing@example.com", // Duplicate email
			Phone:     "+1-555-0101",
		}

		_, err := mockRepo.CreateUser(user)
		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}

		if err.Error() != "duplicate email address" {
			t.Errorf("Expected error 'duplicate email address', got '%s'", err.Error())
		}
	})
}

// TestGetUserByID tests user retrieval with various scenarios
func TestGetUserByID(t *testing.T) {
	t.Run("successful user retrieval", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockGetUserByID: func(id string) (models.User, error) {
				return models.User{
					UserID:    id,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Phone:     "+1-555-0101",
					CreatedAt: time.Now().Add(-24 * time.Hour),
					UpdatedAt: time.Now(),
				}, nil
			},
		}

		user, err := mockRepo.GetUserById("test-123")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if user.UserID != "test-123" {
			t.Errorf("Expected UserID %s, got %s", "test-123", user.UserID)
		}

		if user.FirstName != "John" {
			t.Errorf("Expected FirstName %s, got %s", "John", user.FirstName)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockGetUserByID: func(id string) (models.User, error) {
				return models.User{}, errors.New("user not found")
			},
		}

		_, err := mockRepo.GetUserById("non-existent")
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "user not found" {
			t.Errorf("Expected error 'user not found', got '%s'", err.Error())
		}
	})

	t.Run("empty user ID validation", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockGetUserByID: func(id string) (models.User, error) {
				if id == "" {
					return models.User{}, errors.New("user ID is required")
				}
				return models.User{UserID: id}, nil
			},
		}

		_, err := mockRepo.GetUserById("")
		if err == nil {
			t.Error("Expected error for empty UserID, got nil")
		}

		if err.Error() != "user ID is required" {
			t.Errorf("Expected error 'user ID is required', got '%s'", err.Error())
		}
	})
}

// TestUpdateUser tests user update with business logic validation
func TestUpdateUser(t *testing.T) {
	t.Run("successful user update with timestamp update", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockUpdateUser: func(user models.User) (models.User, error) {
				// Simulate business logic: update timestamp
				user.UpdatedAt = time.Now()
				return user, nil
			},
		}

		user := models.User{
			UserID:    "test-123",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
		}

		updatedUser, err := mockRepo.UpdateUser(user)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Test business logic: UpdatedAt should be set
		if updatedUser.UpdatedAt.IsZero() {
			t.Error("Expected UpdatedAt to be set")
		}

		// Test data integrity: other fields should remain unchanged
		if updatedUser.FirstName != user.FirstName {
			t.Errorf("Expected FirstName %s, got %s", user.FirstName, updatedUser.FirstName)
		}

		if updatedUser.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, updatedUser.Email)
		}
	})

	t.Run("update non-existent user", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockUpdateUser: func(user models.User) (models.User, error) {
				// Simulate user not found
				if user.UserID == "non-existent" {
					return models.User{}, errors.New("user not found")
				}
				user.UpdatedAt = time.Now()
				return user, nil
			},
		}

		user := models.User{
			UserID:    "non-existent",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
		}

		_, err := mockRepo.UpdateUser(user)
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}

		if err.Error() != "user not found" {
			t.Errorf("Expected error 'user not found', got '%s'", err.Error())
		}
	})

	t.Run("update with empty user ID", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockUpdateUser: func(user models.User) (models.User, error) {
				if user.UserID == "" {
					return models.User{}, errors.New("user ID is required")
				}
				user.UpdatedAt = time.Now()
				return user, nil
			},
		}

		user := models.User{
			UserID:    "", // Empty UserID
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
		}

		_, err := mockRepo.UpdateUser(user)
		if err == nil {
			t.Error("Expected error for empty UserID, got nil")
		}

		if err.Error() != "user ID is required" {
			t.Errorf("Expected error 'user ID is required', got '%s'", err.Error())
		}
	})

	t.Run("database error during update", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockUpdateUser: func(user models.User) (models.User, error) {
				return models.User{}, errors.New("database connection failed")
			},
		}

		user := models.User{
			UserID:    "test-123",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Phone:     "+1-555-0101",
		}

		_, err := mockRepo.UpdateUser(user)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "database connection failed" {
			t.Errorf("Expected error 'database connection failed', got '%s'", err.Error())
		}
	})
}

// TestDeleteUser tests user deletion with proper return value handling
func TestDeleteUser(t *testing.T) {
	t.Run("successful user deletion with returned user data", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockDeleteUser: func(id string) (models.User, error) {
				// Simulate business logic: return deleted user data
				return models.User{
					UserID:    id,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					Phone:     "+1-555-0101",
				}, nil
			},
		}

		deletedUser, err := mockRepo.DeleteUser("test-123")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Test business logic: should return deleted user data
		if deletedUser.UserID != "test-123" {
			t.Errorf("Expected UserID %s, got %s", "test-123", deletedUser.UserID)
		}

		if deletedUser.FirstName != "John" {
			t.Errorf("Expected FirstName %s, got %s", "John", deletedUser.FirstName)
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockDeleteUser: func(id string) (models.User, error) {
				if id == "non-existent" {
					return models.User{}, errors.New("user not found")
				}
				return models.User{UserID: id}, nil
			},
		}

		_, err := mockRepo.DeleteUser("non-existent")
		if err == nil {
			t.Error("Expected error for non-existent user, got nil")
		}

		if err.Error() != "user not found" {
			t.Errorf("Expected error 'user not found', got '%s'", err.Error())
		}
	})

	t.Run("delete with empty user ID", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockDeleteUser: func(id string) (models.User, error) {
				if id == "" {
					return models.User{}, errors.New("user ID is required")
				}
				return models.User{UserID: id}, nil
			},
		}

		_, err := mockRepo.DeleteUser("")
		if err == nil {
			t.Error("Expected error for empty UserID, got nil")
		}

		if err.Error() != "user ID is required" {
			t.Errorf("Expected error 'user ID is required', got '%s'", err.Error())
		}
	})
}

// TestListAllUsers tests user listing functionality
func TestListAllUsers(t *testing.T) {
	t.Run("successful user listing", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockListAllUsers: func() ([]models.User, error) {
				return []models.User{
					{
						UserID:    "user-1",
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john@example.com",
						Phone:     "+1-555-0101",
					},
					{
						UserID:    "user-2",
						FirstName: "Jane",
						LastName:  "Smith",
						Email:     "jane@example.com",
						Phone:     "+1-555-0102",
					},
				}, nil
			},
		}

		users, err := mockRepo.ListAllUsers()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(users) != 2 {
			t.Errorf("Expected 2 users, got %d", len(users))
		}

		// Test first user
		if users[0].UserID != "user-1" {
			t.Errorf("Expected first user ID %s, got %s", "user-1", users[0].UserID)
		}

		// Test second user
		if users[1].UserID != "user-2" {
			t.Errorf("Expected second user ID %s, got %s", "user-2", users[1].UserID)
		}
	})

	t.Run("empty user list", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockListAllUsers: func() ([]models.User, error) {
				return []models.User{}, nil
			},
		}

		users, err := mockRepo.ListAllUsers()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Expected 0 users, got %d", len(users))
		}
	})

	t.Run("database error during listing", func(t *testing.T) {
		mockRepo := &MockUserRepository{
			MockListAllUsers: func() ([]models.User, error) {
				return nil, errors.New("database connection failed")
			},
		}

		_, err := mockRepo.ListAllUsers()
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "database connection failed" {
			t.Errorf("Expected error 'database connection failed', got '%s'", err.Error())
		}
	})
}
