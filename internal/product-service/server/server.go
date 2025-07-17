//go:generate go tool oapi-codegen -config cfg.yaml ../../api.yaml
package server

import (
	"gocart/internal/product-service/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	handler *handler.ProductHandler
	router  *mux.Router
}

func NewServer(handler *handler.ProductHandler) *Server {
	s := &Server{
		handler: handler,
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/products", s.handler.ListProducts).Methods("GET")
	s.router.HandleFunc("/products", s.handler.CreateProduct).Methods("POST")
	s.router.HandleFunc("/products/{id}", s.handler.GetProductById).Methods("GET")
	s.router.HandleFunc("/products/{id}", s.handler.UpdateProduct).Methods("PUT")
	s.router.HandleFunc("/products/{id}", s.handler.DeleteProduct).Methods("DELETE")
}

func (s *Server) Start(port string) error {
	log.Printf("Starting server on %s", port)
	return http.ListenAndServe(port, s.router)
}

// Implement the generated interface methods
func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {
	s.handler.ListProducts(w, r)
}

func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	s.handler.CreateProduct(w, r)
}

func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request, id string) {
	s.handler.GetProductById(w, r)
}

func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request, id string) {
	s.handler.UpdateProduct(w, r)
}

func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	s.handler.DeleteProduct(w, r)
}
