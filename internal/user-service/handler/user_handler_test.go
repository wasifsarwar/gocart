package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"gocart/internal/user-service/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type MockUserRepository struct {
	MockListAllUsers func() ([]models.User, error)
	MockCreateUser   func(user models.User) (models.User, error)
	MockGetUserById  func(id string) (models.User, error)
	MockUpdateUser   func(user models.User) (models.User, error)
	MockDeleteUser   func(id string) (models.User, error)
}

func (m *MockUserRepository) ListAllUsers() ([]models.User, error) {
	return m.MockListAllUsers()
}

func (m *MockUserRepository) CreateUser(user models.User) (models.User, error) {
	return m.MockCreateUser(user)
}

func (m *MockUserRepository) GetUserById(id string) (models.User, error) {
	return m.MockGetUserById(id)
}

func (m *MockUserRepository) UpdateUser(user models.User) (models.User, error) {
	return m.MockUpdateUser(user)
}

func (m *MockUserRepository) DeleteUser(id string) (models.User, error) {
	return m.MockDeleteUser(id)
}

func TestListUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockUsers      []models.User
		mockError      error
		expectedStatus int
	}{
		{
			name: "Success",
			mockUsers: []models.User{
				{UserID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "123-456-7890"},
				{UserID: "2", FirstName: "Jane", LastName: "Smith", Email: "jane.smith@example.com", Phone: "123-456-7890"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty User List",
			mockUsers:      []models.User{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Database Error",
			mockUsers:      nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				MockListAllUsers: func() ([]models.User, error) {
					return tt.mockUsers, tt.mockError
				},
			}

			handler := NewUserHandler(mockRepo)
			req := httptest.NewRequest(http.MethodGet, "/users", nil)
			w := httptest.NewRecorder()

			handler.ListAllUsers(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Test Content-Type header
			if w.Header().Get("Content-Type") != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
			}

			if tt.mockError == nil {
				var response []models.User
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if len(response) != len(tt.mockUsers) {
					t.Errorf("Expected %d users, got %d", len(tt.mockUsers), len(response))
				}

				// Only test user data if we have users
				if len(tt.mockUsers) > 0 && len(response) > 0 {
					if response[0].FirstName != tt.mockUsers[0].FirstName {
						t.Errorf("Expected first name %s, got %s", tt.mockUsers[0].FirstName, response[0].FirstName)
					}

					if len(tt.mockUsers) > 1 && len(response) > 1 {
						if response[1].FirstName != tt.mockUsers[1].FirstName {
							t.Errorf("Expected first name %s, got %s", tt.mockUsers[1].FirstName, response[1].FirstName)
						}
					}
				}
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		input          models.User
		requestBody    string
		mockUser       models.User
		mockError      error
		expectedStatus int
		setContentType bool
	}{
		{
			name:           "Success",
			input:          models.User{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "123-456-7890"},
			mockUser:       models.User{UserID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "123-456-7890"},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			setContentType: true,
		},
		{
			name:           "Missing FirstName",
			input:          models.User{LastName: "Doe", Email: "john@example.com", Phone: "123-456-7890"},
			mockUser:       models.User{},
			mockError:      errors.New("first name is required"),
			expectedStatus: http.StatusBadRequest,
			setContentType: true,
		},
		{
			name:           "Invalid Email",
			input:          models.User{FirstName: "John", LastName: "Doe", Email: "invalid-email", Phone: "123-456-7890"},
			mockUser:       models.User{},
			mockError:      errors.New("invalid email format"),
			expectedStatus: http.StatusBadRequest,
			setContentType: true,
		},
		{
			name:           "Empty Request Body",
			input:          models.User{},
			requestBody:    "",
			mockUser:       models.User{},
			mockError:      errors.New("request body is required"),
			expectedStatus: http.StatusBadRequest,
			setContentType: true,
		},
		{
			name:           "Malformed JSON",
			input:          models.User{},
			requestBody:    `{"firstName": "John", "lastName": "Doe", "email": "john@example.com"`,
			mockUser:       models.User{},
			mockError:      errors.New("invalid json"),
			expectedStatus: http.StatusBadRequest,
			setContentType: true,
		},
		{
			name:           "Missing Content-Type",
			input:          models.User{FirstName: "John", LastName: "Doe", Email: "john@example.com"},
			mockUser:       models.User{},
			mockError:      errors.New("content-type must be application/json"),
			expectedStatus: http.StatusBadRequest,
			setContentType: false,
		},
		{
			name:           "Database Error",
			input:          models.User{FirstName: "Jane", LastName: "Hunter", Email: "jane.hunter@example.com", Phone: "123-456-7890"},
			mockUser:       models.User{},
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			setContentType: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				MockCreateUser: func(user models.User) (models.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			handler := NewUserHandler(mockRepo)

			var body []byte
			if tt.requestBody != "" {
				body = []byte(tt.requestBody)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
			if tt.setContentType {
				req.Header.Set("Content-Type", "application/json")
			}
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Test Content-Type header for successful responses
			if tt.expectedStatus == http.StatusCreated {
				if w.Header().Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
				}
			}

			if tt.mockError == nil && tt.expectedStatus == http.StatusCreated {
				var response models.User
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.UserID != tt.mockUser.UserID {
					t.Errorf("Expected user ID %s, got %s", tt.mockUser.UserID, response.UserID)
				}

				if response.FirstName != tt.mockUser.FirstName {
					t.Errorf("Expected first name %s, got %s", tt.mockUser.FirstName, response.FirstName)
				}
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockUser       models.User
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			userID:         "1",
			mockUser:       models.User{UserID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "123-456-7890"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Empty User ID",
			userID:         "",
			mockUser:       models.User{},
			mockError:      errors.New("user ID is required"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Not Found",
			userID:         "999",
			mockUser:       models.User{},
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Database Error",
			userID:         "1",
			mockUser:       models.User{},
			mockError:      errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				MockGetUserById: func(id string) (models.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			handler := NewUserHandler(mockRepo)
			req := httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)
			req = mux.SetURLVars(req, map[string]string{"user_id": tt.userID})
			w := httptest.NewRecorder()

			handler.GetUserById(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Test Content-Type header for successful responses
			if tt.expectedStatus == http.StatusOK {
				if w.Header().Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
				}
			}

			if tt.mockError == nil && tt.expectedStatus == http.StatusOK {
				var response models.User
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.UserID != tt.mockUser.UserID {
					t.Errorf("Expected user ID %s, got %s", tt.mockUser.UserID, response.UserID)
				}

				if response.FirstName != tt.mockUser.FirstName {
					t.Errorf("Expected first name %s, got %s", tt.mockUser.FirstName, response.FirstName)
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		input          models.User
		requestBody    string
		mockUser       models.User
		mockError      error
		expectedStatus int
		setContentType bool
	}{
		{
			name:           "Success",
			userID:         "1",
			input:          models.User{FirstName: "Updated", LastName: "User", Email: "updated@example.com", Phone: "123-456-7890"},
			mockUser:       models.User{UserID: "1", FirstName: "Updated", LastName: "User", Email: "updated@example.com", Phone: "123-456-7890"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			setContentType: true,
		},
		{
			name:           "Empty User ID",
			userID:         "",
			input:          models.User{FirstName: "Updated", LastName: "User", Email: "updated@example.com"},
			mockUser:       models.User{},
			mockError:      errors.New("user ID is required"),
			expectedStatus: http.StatusBadRequest,
			setContentType: true,
		},
		{
			name:           "Not Found",
			userID:         "999",
			input:          models.User{FirstName: "Updated", LastName: "User", Email: "updated@example.com", Phone: "123-456-7890"},
			mockUser:       models.User{},
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			setContentType: true,
		},
		{
			name:           "Malformed JSON",
			userID:         "1",
			input:          models.User{},
			requestBody:    `{"firstName": "Updated", "lastName": "User"`,
			mockUser:       models.User{},
			mockError:      errors.New("invalid json"),
			expectedStatus: http.StatusBadRequest,
			setContentType: true,
		},
		{
			name:           "Missing Content-Type",
			userID:         "1",
			input:          models.User{FirstName: "Updated", LastName: "User", Email: "updated@example.com"},
			mockUser:       models.User{},
			mockError:      errors.New("content-type must be application/json"),
			expectedStatus: http.StatusBadRequest,
			setContentType: false,
		},
		{
			name:           "Database Error",
			userID:         "1",
			input:          models.User{FirstName: "Updated", LastName: "User", Email: "updated@example.com"},
			mockUser:       models.User{},
			mockError:      errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
			setContentType: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				MockGetUserById: func(id string) (models.User, error) {
					if tt.mockError != nil && tt.expectedStatus == http.StatusNotFound {
						return models.User{}, tt.mockError
					}
					return tt.mockUser, nil
				},
				MockUpdateUser: func(user models.User) (models.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			handler := NewUserHandler(mockRepo)

			var body []byte
			if tt.requestBody != "" {
				body = []byte(tt.requestBody)
			} else {
				body, _ = json.Marshal(tt.input)
			}

			req := httptest.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(body))
			if tt.setContentType {
				req.Header.Set("Content-Type", "application/json")
			}
			req = mux.SetURLVars(req, map[string]string{"user_id": tt.userID})
			w := httptest.NewRecorder()

			handler.UpdateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Test Content-Type header for successful responses
			if tt.expectedStatus == http.StatusOK {
				if w.Header().Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type application/json, got %s", w.Header().Get("Content-Type"))
				}
			}

			if tt.mockError == nil && tt.expectedStatus == http.StatusOK {
				var response models.User
				err := json.NewDecoder(w.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if response.UserID != tt.mockUser.UserID {
					t.Errorf("Expected user ID %s, got %s", tt.mockUser.UserID, response.UserID)
				}

				if response.FirstName != tt.mockUser.FirstName {
					t.Errorf("Expected first name %s, got %s", tt.mockUser.FirstName, response.FirstName)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockUser       models.User
		mockError      error
		expectedStatus int
	}{
		{
			name:           "Success",
			userID:         "1",
			mockUser:       models.User{UserID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com", Phone: "123-456-7890"},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "Empty User ID",
			userID:         "",
			mockUser:       models.User{},
			mockError:      errors.New("user ID is required"),
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Not Found",
			userID:         "999",
			mockUser:       models.User{},
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Database Error",
			userID:         "1",
			mockUser:       models.User{},
			mockError:      errors.New("database connection failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUserRepository{
				MockDeleteUser: func(id string) (models.User, error) {
					return tt.mockUser, tt.mockError
				},
			}

			handler := NewUserHandler(mockRepo)
			req := httptest.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)
			req = mux.SetURLVars(req, map[string]string{"user_id": tt.userID})
			w := httptest.NewRecorder()

			handler.DeleteUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Only check response body for error cases, not for successful delete
			// Successful delete returns 204 NoContent with no body
			if tt.mockError != nil {
				// For error cases, check that we get an error message
				if w.Body.Len() == 0 {
					t.Error("Expected error message in response body, got empty body")
				}
			} else {
				// For successful delete, expect empty body (204 NoContent)
				if w.Body.Len() != 0 {
					t.Errorf("Expected empty response body for successful delete, got: %s", w.Body.String())
				}
			}
		})
	}
}
