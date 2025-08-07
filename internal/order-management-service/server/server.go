package server

import (
	"gocart/internal/order-management-service/handler"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	handler *handler.OrderHandler
	router  *mux.Router
}

func NewServer(handler *handler.OrderHandler) *Server {
	s := &Server{
		handler: handler,
		router:  mux.NewRouter(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.HandleFunc("/orders", s.handler.CreateOrder).Methods("POST")
	s.router.HandleFunc("/orders/{id}", s.handler.GetOrderById).Methods("GET")
	s.router.HandleFunc("/orders/{id}", s.handler.UpdateOrder).Methods("PUT")
	s.router.HandleFunc("/orders/{id}", s.handler.DeleteOrder).Methods("DELETE")
	s.router.HandleFunc("/orders/{id}/items", s.handler.DeleteOrderItem).Methods("DELETE")
	s.router.HandleFunc("/orders", s.handler.ListAllOrders).Methods("GET")
	s.router.HandleFunc("/orders/user/{user_id}", s.handler.ListOrdersByUserId).Methods("GET")
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	s.handler.CreateOrder(w, r)
}

func (s *Server) GetOrderById(w http.ResponseWriter, r *http.Request) {
	s.handler.GetOrderById(w, r)
}

func (s *Server) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	s.handler.UpdateOrder(w, r)
}

func (s *Server) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	s.handler.DeleteOrder(w, r)
}

func (s *Server) DeleteOrderItem(w http.ResponseWriter, r *http.Request) {
	s.handler.DeleteOrderItem(w, r)
}

func (s *Server) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	s.handler.ListAllOrders(w, r)
}

func (s *Server) ListOrdersByUserId(w http.ResponseWriter, r *http.Request) {
	s.handler.ListOrdersByUserId(w, r)
}
