package handler

import (
	"encoding/json"
	"gocart/user-service/models"
	"gocart/user-service/repository"
	"log"
	"net/http"
	"time"

	shared "gocart/pkg/db"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
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
	createdUser, err := repository.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	user, err := repository.GetUserById(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var updatedUser models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// First get the existing user to update
	existingUser, err := repository.GetUserById(userID)
	if err != nil {
		log.Printf("Error fetching user with id: %v and error: %v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update fields
	updatedUser = updateUser(existingUser, updatedUser)

	user, err := func() (models.User, error) {
		var user models.User = updatedUser
		user.UpdatedAt = time.Now()
		if err := shared.DB.Model(&models.User{}).Where("user_id = ?", user.UserID).Updates(&user).Error; err != nil {
			return models.User{}, err
		}
		return user, nil
	}()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(user)
}

func updateUser(existingUser models.User, updatedUser models.User) models.User {
	updatedUser.UserID = existingUser.UserID
	updatedUser.CreatedAt = existingUser.CreatedAt
	updatedUser.UpdatedAt = time.Now()
	updatedUser.FirstName = existingUser.FirstName
	updatedUser.LastName = existingUser.LastName
	updatedUser.Email = existingUser.Email
	updatedUser.Phone = existingUser.Phone
	return updatedUser
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	_, err := repository.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}

func ListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := repository.ListAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
