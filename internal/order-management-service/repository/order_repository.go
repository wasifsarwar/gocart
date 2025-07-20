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
	DeleteOrderItem(orderItemID string) error
	ListAllOrders(limit, offset int) ([]models.Order, error)
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

	// validate user exists
	if err := r.validateUserExists(order.UserID); err != nil {
		return models.Order{}, fmt.Errorf("user validation failed: %w", err)
	}

	// Step 1: Basic validation (expand as needed)
	if len(order.Items) == 0 {
		return models.Order{}, errors.New("order must have at least one item")
	}

	for i, item := range order.Items {
		if item.ProductID == "" {
			return models.Order{}, errors.New("invalid order item: missing product_id")
		}
		if item.Quantity <= 0 {
			return models.Order{}, errors.New("invalid order item: quantity must be greater than 0")
		}
		if item.Price <= 0 {
			return models.Order{}, errors.New("invalid order item: price must be greater than 0")
		}

		currentPrice, err := r.validateProductAndFetchPrice(item.ProductID)
		if err != nil {
			return models.Order{}, fmt.Errorf("product validation failed for item %d: %w", i+1, err)
		}
		order.Items[i].Price = currentPrice
	}

	// Step 2: Create order and items in a transaction
	order.OrderID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.TotalAmount = calculateTotal(order.Items)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Create order without items first
		orderWithoutItems := models.Order{
			OrderID:     order.OrderID,
			UserID:      order.UserID,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			CreatedAt:   order.CreatedAt,
			UpdatedAt:   order.UpdatedAt,
		}
		if err := tx.Create(&orderWithoutItems).Error; err != nil {
			return err
		}

		// Create items separately
		for i := range order.Items {
			order.Items[i].OrderItemID = uuid.New().String()
			order.Items[i].OrderID = orderWithoutItems.OrderID
			order.Items[i].CreatedAt = time.Now()
			order.Items[i].UpdatedAt = time.Now()
			if err := tx.Create(&order.Items[i]).Error; err != nil {
				return fmt.Errorf("failed to create order item %s: %w", order.Items[i].OrderItemID, err)
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
	// atomic update of order and items
	err = r.db.Transaction(func(tx *gorm.DB) error {

		// update order-level fields only if provided (partial update)
		orderUpdates := make(map[string]interface{})
		if order.Status != "" {
			orderUpdates["status"] = order.Status
		}
		orderUpdates["updated_at"] = time.Now()

		// update parent order
		if len(orderUpdates) > 0 {
			if err := tx.Model(&existingOrder).Updates(orderUpdates).Error; err != nil {
				return fmt.Errorf("failed to update order %s: %w", existingOrder.OrderID, err)
			}
		}

		// process order items: handle deletions, updates and inserts
		for i := range order.Items {

			order.Items[i].UpdatedAt = time.Now()
			order.Items[i].OrderID = existingOrder.OrderID

			// check if item is marked for deletion
			if order.Items[i].Delete {
				if order.Items[i].OrderItemID == "" {
					return fmt.Errorf("invalid order item: missing order_item_id")
				}
				// delete item with proper WHERE clause
				if err := tx.Where("order_item_id = ?", order.Items[i].OrderItemID).Delete(&models.OrderItem{}).Error; err != nil {
					return fmt.Errorf("failed to delete order item %s: %w", order.Items[i].OrderItemID, err)
				}

				continue // Skip further processing for deleted items
			}

			if order.Items[i].OrderItemID != "" {
				// update existing item
				itemUpdates := make(map[string]interface{})
				if order.Items[i].ProductID != "" {
					// Validate product exists
					currentPrice, err := r.validateProductAndFetchPrice(order.Items[i].ProductID)
					if err != nil {
						return fmt.Errorf("product validation failed for item update: %w", err)
					}
					order.Items[i].Price = currentPrice
					itemUpdates["product_id"] = order.Items[i].ProductID
				}
				if order.Items[i].Quantity > 0 {
					itemUpdates["quantity"] = order.Items[i].Quantity
				}
				if order.Items[i].Price > 0 {
					itemUpdates["price"] = order.Items[i].Price
				}
				itemUpdates["updated_at"] = time.Now()
				if len(itemUpdates) > 0 {
					if err := tx.Model(&models.OrderItem{}).Where("order_item_id = ?", order.Items[i].OrderItemID).Updates(itemUpdates).Error; err != nil {
						return fmt.Errorf("failed to update order item %s: %w", order.Items[i].OrderItemID, err)
					}
				}
			} else {
				// Validate required fields for new items
				if order.Items[i].ProductID == "" {
					return fmt.Errorf("new item missing product_id")
				}
				if order.Items[i].Quantity <= 0 {
					return fmt.Errorf("new item invalid quantity")
				}
				// validate product exists, and get current price for new items
				currentPrice, err := r.validateProductAndFetchPrice(order.Items[i].ProductID)
				if err != nil {
					return fmt.Errorf("product validation failed for new item %d: %w", i+1, err)
				}
				order.Items[i].Price = currentPrice

				// create new item
				order.Items[i].OrderItemID = uuid.New().String()
				order.Items[i].OrderID = existingOrder.OrderID
				order.Items[i].CreatedAt = time.Now()
				order.Items[i].UpdatedAt = time.Now()
				if err := tx.Create(&order.Items[i]).Error; err != nil {
					return fmt.Errorf("failed to create order item %s: %w", order.Items[i].OrderItemID, err)
				}
			}
		}

		//recalculate total amount based on all current items
		var allItems []models.OrderItem
		if err := tx.Where("order_id = ?", order.OrderID).Find(&allItems).Error; err != nil {
			return fmt.Errorf("failed to fetch order items for order %s: %w", order.OrderID, err)
		}
		newTotal := calculateTotal(allItems)
		if err := tx.Model(&existingOrder).Update("total_amount", newTotal).Error; err != nil {
			return fmt.Errorf("failed to update order total amount for order %s: %w", order.OrderID, err)
		}
		return nil

	})

	if err != nil {
		return models.Order{}, fmt.Errorf("failed to update order %s: %w", existingOrder.OrderID, err)
	}

	// Return updated order with all items
	return r.GetOrderById(existingOrder.OrderID)
}

func (r *orderRepository) DeleteOrderItem(orderItemID string) error {
	result := r.db.Where("order_item_id = ?", orderItemID).Delete(&models.OrderItem{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete order item %s: %w", orderItemID, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.New("order item not found")
	}
	return nil
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

func (r *orderRepository) ListAllOrders(limit, offset int) ([]models.Order, error) {
	var orders []models.Order

	if limit <= 0 {
		limit = 10 // default page size
	}

	if offset < 0 {
		offset = 0
	}

	if limit > 100 {
		limit = 100 // maximum page size
	}

	if err := r.db.Preload("Items").Order("created_at DESC").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
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

func (r *orderRepository) validateUserExists(userId string) error {
	var count int64
	err := r.db.Table("users").Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to validate user exists: %w", err)
	}
	if count == 0 {
		return fmt.Errorf("user %s not found", userId)
	}
	return nil
}

func (r *orderRepository) validateProductAndFetchPrice(productID string) (float64, error) {
	var product struct {
		Price float64 `gorm:"column:price"`
	}
	err := r.db.Table("products").
		Select("price").
		Where("product_id = ?", productID).
		First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("product %s not found", productID)
		}
		return 0, fmt.Errorf("failed to validate product exists: %w", err)
	}
	return product.Price, nil
}

func calculateTotal(items []models.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
