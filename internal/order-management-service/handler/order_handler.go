package handler

import (
	"encoding/json"
	"fmt"
	"gocart/internal/order-management-service/models"
	"gocart/internal/order-management-service/repository"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderRepo repository.OrderRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo: orderRepo,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	order, err := h.orderRepo.CreateOrder(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/orders/"+order.OrderID)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	order, err := h.orderRepo.GetOrderById(orderId)
	if err != nil {
		log.Printf("Error fetching order with id: %v and error: %v", orderId, err)
		if err.Error() == "order not found" {
			http.Error(w, fmt.Sprintf("Order with id %v not found.", orderId), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Unable to retrieve order with id: %v.", orderId), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	var updatedOrder models.Order
	if err := json.NewDecoder(r.Body).Decode(&updatedOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	existingOrder, err := h.orderRepo.GetOrderById(orderId)
	if err != nil {
		log.Printf("Error fetching order with id: %v and error: %v", orderId, err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	updatedOrder.OrderID = existingOrder.OrderID

	result, err := h.orderRepo.UpdateOrder(updatedOrder)
	if err != nil {
		log.Printf("Error updating order with id: %v and error: %v", orderId, err)
		http.Error(w, "Unable to update order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
	json.NewEncoder(w).Encode(map[string]string{"message": "Order updated successfully"})
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	orderId := mux.Vars(r)["id"]
	if err := h.orderRepo.DeleteOrder(orderId); err != nil {
		log.Printf("Error deleting order with id: %v and error: %v", orderId, err)
		http.Error(w, "Unable to delete order", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.orderRepo.ListAllOrders()
	if err != nil {
		log.Printf("Error listing all orders and error: %v", err)
		http.Error(w, "Unable to list orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) ListOrdersByUserId(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user_id"]
	orders, err := h.orderRepo.ListOrdersByUserId(userId)
	if err != nil {
		log.Printf("Error listing orders by user id: %v and error: %v", userId, err)
		http.Error(w, "Unable to list orders", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
