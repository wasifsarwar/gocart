package main

import (
	db "gocart/shared/db"
	"user-service/handler"
	"user-service/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/google/uuid"
)

func main() {
	db.Connect(db.DefaultConfig())
	db.Migrate(&models.User{})

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

