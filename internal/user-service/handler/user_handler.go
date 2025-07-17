package handler

import (
	"encoding/json"
	"gocart/internal/user-service/models"
	"gocart/internal/user-service/repository"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	// Generate a new UUID for the user
	user.UserID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Save to database
	createdUser, err := h.repo.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (h *UserHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	user, err := h.repo.GetUserById(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get existing user
	existingUser, err := h.repo.GetUserById(userID)
	if err != nil {
		log.Printf("Error fetching user with id: %v and error: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields (implement updateUser logic as needed)
	updatedUser.UserID = existingUser.UserID
	updatedUser.CreatedAt = existingUser.CreatedAt
	updatedUser.UpdatedAt = time.Now()

	result, err := h.repo.UpdateUser(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	_, err := h.repo.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
