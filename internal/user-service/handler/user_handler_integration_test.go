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
