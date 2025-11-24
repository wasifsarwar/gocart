package server

import (
	"gocart/internal/user-service/handler"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	handler *handler.UserHandler
	router  *mux.Router
}

func NewServer(handler *handler.UserHandler) *Server {
	s := &Server{
		handler: handler,
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/users", s.handler.ListAllUsers).Methods("GET")
	s.router.HandleFunc("/users/register", s.handler.CreateUser).Methods("POST")
	s.router.HandleFunc("/users/login", s.handler.Login).Methods("POST")
	s.router.HandleFunc("/users/{id}", s.handler.GetUserById).Methods("GET")
	s.router.HandleFunc("/users/{id}", s.handler.UpdateUser).Methods("PUT")
	s.router.HandleFunc("/users/{id}", s.handler.DeleteUser).Methods("DELETE")
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	s.handler.CreateUser(w, r)
}

func (s *Server) GetUserById(w http.ResponseWriter, r *http.Request) {
	s.handler.GetUserById(w, r)
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	s.handler.UpdateUser(w, r)
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	s.handler.DeleteUser(w, r)
}

func (s *Server) ListAllUsers(w http.ResponseWriter, r *http.Request) {
	s.handler.ListAllUsers(w, r)
}
