package handler

import (
	"bytes"
	"encoding/json"
	"gocart/internal/user-service/models"
	"gocart/internal/user-service/repository"
	"gocart/pkg/testutils"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	config := testutils.TestDBConfig{
		ServiceName: "users_handler",
		Models:      []interface{}{&models.User{}},
	}
	return testutils.SetupTestDB(t, config)
}

func TestCreateAndGetUserIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestCreateAndGetUserIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting create and get user integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)
	testUser := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "123-456-7890",
		Password:  "password123",
	}
	body, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("Failed to marshal test user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
	var createdUserResponse models.User
	if err := json.Unmarshal(w.Body.Bytes(), &createdUserResponse); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if createdUserResponse.UserID == "" {
		t.Error("Expected UserID to be set, got empty string")
	}
	if createdUserResponse.FirstName != testUser.FirstName {
		t.Errorf("Expected first name %s, got %s", testUser.FirstName, createdUserResponse.FirstName)
	}

	getUserID := createdUserResponse.UserID

	getReq := httptest.NewRequest(http.MethodGet, "/users/"+getUserID, nil)
	getReq = mux.SetURLVars(getReq, map[string]string{"user_id": getUserID})
	getW := httptest.NewRecorder()

	handler.GetUserById(getW, getReq)

	if getW.Code != http.StatusOK {
		t.Errorf("Expected status code %d for GET request, got %d. Response body: %s", http.StatusOK, getW.Code, getW.Body.String())
	}

	var getResponse models.User
	if err := json.Unmarshal(getW.Body.Bytes(), &getResponse); err != nil {
		t.Fatalf("Failed to unmarshal get response: %v. Body: %s", err, getW.Body.String())
	}

	if getResponse.UserID != getUserID {
		t.Errorf("Expected fetched user ID %s, got %s", getUserID, getResponse.UserID)
	}
	if getResponse.FirstName != testUser.FirstName {
		t.Errorf("Expected first name %s, got %s", testUser.FirstName, getResponse.FirstName)
	}
}

func TestListUsersIntegration(t *testing.T) {
	logger := log.New(os.Stdout, "[TestListUsersIntegration]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting list users integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)

	user1 := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "123-456-7890",
	}
	user2 := models.User{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane.smith@example.com",
		Phone:     "123-456-7890",
	}
	user3 := models.User{
		FirstName: "Jim",
		LastName:  "Beam",
		Email:     "jim.beam@example.com",
		Phone:     "123-456-7890",
	}

	userRepo.CreateUser(user1)
	userRepo.CreateUser(user2)
	userRepo.CreateUser(user3)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler.ListAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response []models.User
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 users, got %d", len(response))
	}

	if response[0].FirstName != user1.FirstName {
		t.Errorf("Expected first name %s, got %s", user1.FirstName, response[0].FirstName)
	}

	if response[1].FirstName != user2.FirstName {
		t.Errorf("Expected first name %s, got %s", user2.FirstName, response[1].FirstName)
	}

	if response[2].FirstName != user3.FirstName {
		t.Errorf("Expected first name %s, got %s", user3.FirstName, response[2].FirstName)
	}
}

func TestCreateUserValidationErrors(t *testing.T) {
	logger := log.New(os.Stdout, "[TestCreateUserValidationErrors]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting create user validation errors integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)

	tests := []struct {
		name           string
		user           models.User
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing FirstName",
			user: models.User{
				LastName: "Doe",
				Email:    "john.doe@example.com",
				Phone:    "123-456-7890",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "First name is required",
		},
		{
			name: "Missing LastName",
			user: models.User{
				FirstName: "John",
				Email:     "john.doe@example.com",
				Phone:     "123-456-7890",
				Password:  "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Last name is required",
		},
		{
			name: "Missing Email",
			user: models.User{
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "123-456-7890",
				Password:  "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Email is required",
		},
		{
			name: "Missing Phone",
			user: models.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
				Password:  "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Phone is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.user)
			if err != nil {
				t.Fatalf("Failed to marshal test user: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, w.Code)
			}

			if !strings.Contains(w.Body.String(), tt.expectedError) {
				t.Errorf("Expected error message containing '%s', got '%s'", tt.expectedError, w.Body.String())
			}
		})
	}
}

func TestGetUserNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "[TestGetUserNotFound]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting get user not found integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)

	// Try to get a user that doesn't exist
	req := httptest.NewRequest(http.MethodGet, "/users/non-existent-id", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "non-existent-id"})
	w := httptest.NewRecorder()

	handler.GetUserById(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), "not found") {
		t.Errorf("Expected error message containing 'not found', got '%s'", w.Body.String())
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	logger := log.New(os.Stdout, "[TestCreateUserDuplicateEmail]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting create user duplicate email integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)

	// Create first user
	user1 := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Phone:     "123-456-7890",
		Password:  "password123",
	}
	body1, _ := json.Marshal(user1)
	req1 := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	handler.CreateUser(w1, req1)

	if w1.Code != http.StatusCreated {
		t.Fatalf("Failed to create first user, got status %d", w1.Code)
	}

	// Try to create second user with same email (should fail)
	user2 := models.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "john.doe@example.com", // Same email
		Phone:     "123-456-7890",
		Password:  "password123",
	}
	body2, _ := json.Marshal(user2)
	req2 := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	handler.CreateUser(w2, req2)

	// Should get 409 Conflict for duplicate email
	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status code %d for duplicate email, got %d", http.StatusConflict, w2.Code)
	}

	if !strings.Contains(w2.Body.String(), "already exists") {
		t.Errorf("Expected error message containing 'already exists', got '%s'", w2.Body.String())
	}
}

func TestUpdateUserNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "[TestUpdateUserNotFound]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting update user not found integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)

	// Try to update a user that doesn't exist
	updateUser := models.User{
		FirstName: "Updated",
		LastName:  "User",
		Email:     "updated@example.com",
		Phone:     "123-456-7890",
	}
	body, _ := json.Marshal(updateUser)
	req := httptest.NewRequest(http.MethodPut, "/users/non-existent-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req = mux.SetURLVars(req, map[string]string{"user_id": "non-existent-id"})
	w := httptest.NewRecorder()

	handler.UpdateUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), "not found") {
		t.Errorf("Expected error message containing 'not found', got '%s'", w.Body.String())
	}
}

func TestDeleteUserNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "[TestDeleteUserNotFound]", log.Ltime|log.Lmicroseconds)
	logger.Println("Starting delete user not found integration test")

	db, cleanup := setupTestDB(t)
	defer cleanup()

	userRepo := repository.NewUserRepository(db)
	handler := NewUserHandler(userRepo)

	// Try to delete a user that doesn't exist
	req := httptest.NewRequest(http.MethodDelete, "/users/non-existent-id", nil)
	req = mux.SetURLVars(req, map[string]string{"user_id": "non-existent-id"})
	w := httptest.NewRecorder()

	handler.DeleteUser(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), "not found") {
		t.Errorf("Expected error message containing 'not found', got '%s'", w.Body.String())
	}
}
