package main

import (
	"fmt"
	db "gocart/shared/db"
	"log"
	"net/http"
	"user-service/handler"
	"user-service/models"
	"user-service/repository"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	db.Connect(db.DefaultConfig())
	db.Migrate(&models.User{})

	// TestCreateUser()
	r := mux.NewRouter()

	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{user_id}", handler.GetUserById).Methods("GET")
	r.HandleFunc("/users/{user_id}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{user_id}", handler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", handler.ListAllUsers).Methods("GET")

	log.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", r)
}

func TestCreateUser() {
	user := models.User{
		UserID:    fmt.Sprintf("Test-User-%d", uuid.New().String()),
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	}

	createdUser, err := repository.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
	}

	log.Printf("User created: %v", createdUser)
}
