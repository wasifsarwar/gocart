package handler

import (
	"encoding/json"
	"fmt"
	"gocart/internal/user-service/models"
	"gocart/internal/user-service/repository"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(user.FirstName) == "" {
		http.Error(w, "First name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user.LastName) == "" {
		http.Error(w, "Last name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user.Phone) == "" {
		http.Error(w, "Phone is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(user.Password) == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.PasswordHash = string(hashedPassword)
	user.Password = "" // Don't save plain text password

	// Generate a new UUID for the user
	user.UserID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save to database
	createdUser, err := h.repo.CreateUser(user)
	if err != nil {
		// Check for specific error types
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate") || strings.Contains(lower, "unique") || strings.Contains(lower, "already exists") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
		} else {
			log.Printf("Error creating user: %v", err)
			http.Error(w, "Unable to create user", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByEmail(credentials.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// In a real app, we would generate a JWT token here
	// For now, we'll just return the user object (excluding password hash)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserById(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			http.Error(w, fmt.Sprintf("User with id %v not found", userID), http.StatusNotFound)
		} else {
			log.Printf("Error fetching user with id %v: %v", userID, err)
			http.Error(w, fmt.Sprintf("Unable to retrieve user with id %v", userID), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if strings.TrimSpace(updatedUser.FirstName) == "" {
		http.Error(w, "First name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(updatedUser.LastName) == "" {
		http.Error(w, "Last name is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(updatedUser.Email) == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(updatedUser.Phone) == "" {
		http.Error(w, "Phone is required", http.StatusBadRequest)
		return
	}

	// Get existing user
	existingUser, err := h.repo.GetUserById(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			http.Error(w, fmt.Sprintf("User with id %v not found", userID), http.StatusNotFound)
		} else {
			log.Printf("Error fetching user with id %v: %v", userID, err)
			http.Error(w, fmt.Sprintf("Unable to retrieve user with id %v", userID), http.StatusInternalServerError)
		}
		return
	}

	// Update fields
	updatedUser.UserID = existingUser.UserID
	updatedUser.CreatedAt = existingUser.CreatedAt
	updatedUser.UpdatedAt = time.Now()

	result, err := h.repo.UpdateUser(updatedUser)
	if err != nil {
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "duplicate") || strings.Contains(lower, "unique") || strings.Contains(lower, "already exists") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
		} else {
			log.Printf("Error updating user with id %v: %v", userID, err)
			http.Error(w, "Unable to update user", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	_, err := h.repo.DeleteUser(userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows") {
			http.Error(w, fmt.Sprintf("User with id %v not found", userID), http.StatusNotFound)
		} else {
			log.Printf("Error deleting user with id %v: %v", userID, err)
			http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		http.Error(w, "Unable to retrieve users", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
