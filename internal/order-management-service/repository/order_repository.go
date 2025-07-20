package repository

import (
	"errors"
	"fmt"
	"gocart/internal/order-management-service/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order models.Order) (models.Order, error)
	GetOrderById(id string) (models.Order, error)
	UpdateOrder(order models.Order) (models.Order, error)
	DeleteOrder(id string) error
	ListAllOrders() ([]models.Order, error)
	ListOrdersByUserId(userId string) ([]models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func (r *orderRepository) CreateOrder(order models.Order) (models.Order, error) {

	// Step 1: Basic validation (expand as needed)
	if len(order.Items) == 0 {
		return models.Order{}, errors.New("order must have at least one item")
	}

	for _, item := range order.Items {
		if item.ProductID == "" {
			return models.Order{}, errors.New("invalid order item: missing product_id")
		}
		if item.Quantity <= 0 {
			return models.Order{}, errors.New("invalid order item: quantity must be greater than 0")
		}
		if item.Price <= 0 {
			return models.Order{}, errors.New("invalid order item: price must be greater than 0")
		}
	}

	// Step 2: Create order and items in a transaction
	order.OrderID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.TotalAmount = calculateTotal(order.Items)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		for _, item := range order.Items {
			item.OrderItemID = uuid.New().String()
			item.OrderID = order.OrderID
			item.CreatedAt = time.Now()
			item.UpdatedAt = time.Now()
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) GetOrderById(id string) (models.Order, error) {
	var order models.Order

	//Preload associated items
	err := r.db.Preload("Items").Where("order_id = ?", id).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Order{}, errors.New("order not found")
		}
		return models.Order{}, err
	}
	return order, nil
}

func (r *orderRepository) UpdateOrder(order models.Order) (models.Order, error) {

	existingOrder, err := r.GetOrderById(order.OrderID)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to get order %s: %w", order.OrderID, err)
	}

	existingOrder.Status = order.Status
	existingOrder.UpdatedAt = time.Now()
	existingOrder.Items = order.Items
	existingOrder.TotalAmount = calculateTotal(existingOrder.Items)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Order{}).Where("order_id = ?", existingOrder.OrderID).Updates(&order).Error; err != nil {
			return fmt.Errorf("failed to update order %s: %w", existingOrder.OrderID, err)
		}
		if err := tx.Delete(&models.OrderItem{}, "order_id = ?", existingOrder.OrderID).Error; err != nil {
			return fmt.Errorf("failed to delete order items for order %s: %w", existingOrder.OrderID, err)
		}
		for _, item := range existingOrder.Items {
			item.OrderID = existingOrder.OrderID
			item.CreatedAt = existingOrder.CreatedAt
			item.UpdatedAt = time.Now()
			if err := tx.Create(&item).Error; err != nil {
				return fmt.Errorf("failed to create order item %s: %w", item.OrderItemID, err)
			}
		}
		return nil
	})

	if err != nil {
		return models.Order{}, fmt.Errorf("failed to update order %s: %w", existingOrder.OrderID, err)
	}

	return r.GetOrderById(existingOrder.OrderID)
}

func (r *orderRepository) DeleteOrder(id string) error {
	result := r.db.Delete(&models.Order{}, "order_id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order %s not found", id)
	}
	return nil
}

func (r *orderRepository) ListAllOrders() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("Items").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) ListOrdersByUserId(userId string) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("Items").Where("user_id = ?", userId).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to list orders by user %s: %w", userId, err)
	}
	return orders, nil
}

func calculateTotal(items []models.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
