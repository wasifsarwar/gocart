//go:generate go tool oapi-codegen -config cfg.yaml ../../api.yaml
package server

import (
	"net/http"

	"product-service/gen"
	"product-service/internal/handler"

	"github.com/gorilla/mux"
)

type Server struct {
	productHandler *handler.ProductHandler
}

func NewServer(productHandler *handler.ProductHandler) *Server {
	return &Server{
		productHandler: productHandler,
	}
}

func (s *Server) RegisterRoutes(r *mux.Router) {
	// Create a new StrictServer
	strictHandler := gen.NewStrictHandler(s, nil)

	// Register the routes
	gen.RegisterHandlers(r, strictHandler)
}

// Implement the generated interface methods
func (s *Server) ListProducts(w http.ResponseWriter, r *http.Request) {
	s.productHandler.ListProducts(w, r)
}

func (s *Server) CreateProduct(w http.ResponseWriter, r *http.Request) {
	s.productHandler.CreateProduct(w, r)
}

func (s *Server) GetProductById(w http.ResponseWriter, r *http.Request, id string) {
	s.productHandler.GetProductById(w, r, id)
}

func (s *Server) UpdateProduct(w http.ResponseWriter, r *http.Request, id string) {
	s.productHandler.UpdateProduct(w, r, id)
}

func (s *Server) DeleteProduct(w http.ResponseWriter, r *http.Request, id string) {
	s.productHandler.DeleteProduct(w, r, id)
}
