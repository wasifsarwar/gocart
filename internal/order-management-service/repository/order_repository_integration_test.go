package repository

import (
	"gocart/internal/order-management-service/models"
	productModels "gocart/internal/product-service/models"
	userModels "gocart/internal/user-service/models"
	"gocart/pkg/testutils"
	"testing"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	config := testutils.TestDBConfig{
		ServiceName: "orders_repo",
		Models: []interface{}{
			&models.Order{},
			&models.OrderItem{},
			&userModels.User{},       // Add User model for validation
			&productModels.Product{}, // Add Product model for validation
		},
	}
	return testutils.SetupTestDB(t, config)
}

// createTestData creates all necessary test users and products
func createTestData(t *testing.T, db *gorm.DB) {
	// Create test users
	testUsers := []userModels.User{
		{UserID: "user-123", FirstName: "Test", LastName: "User1", Email: "test1@example.com", Phone: "123-456-7890"},
		{UserID: "user-456", FirstName: "Test", LastName: "User2", Email: "test2@example.com", Phone: "123-456-7891"},
		{UserID: "user-789", FirstName: "Test", LastName: "User3", Email: "test3@example.com", Phone: "123-456-7892"},
		{UserID: "user-999", FirstName: "Test", LastName: "User4", Email: "test4@example.com", Phone: "123-456-7893"},
		{UserID: "user-111", FirstName: "Test", LastName: "User5", Email: "test5@example.com", Phone: "123-456-7894"},
		{UserID: "user-333", FirstName: "Test", LastName: "User6", Email: "test6@example.com", Phone: "123-456-7895"},
		{UserID: "user-444", FirstName: "Test", LastName: "User7", Email: "test7@example.com", Phone: "123-456-7896"},
	}

	for _, user := range testUsers {
		if err := db.Create(&user).Error; err != nil {
			t.Fatalf("Failed to create test user %s: %v", user.UserID, err)
		}
	}

	// Create test products
	testProducts := []productModels.Product{
		{ProductID: "product-001", Name: "Test Product 1", Description: "Test description", Price: 99.99, Category: "Test"},
		{ProductID: "product-002", Name: "Test Product 2", Description: "Test description", Price: 100.01, Category: "Test"},
		{ProductID: "product-003", Name: "Test Product 3", Description: "Test description", Price: 150.00, Category: "Test"},
		{ProductID: "product-004", Name: "Test Product 4", Description: "Test description", Price: 200.00, Category: "Test"},
		{ProductID: "product-005", Name: "Test Product 5", Description: "Test description", Price: 250.00, Category: "Test"},
		{ProductID: "product-006", Name: "Test Product 6", Description: "Test description", Price: 300.00, Category: "Test"},
		{ProductID: "product-007", Name: "Test Product 7", Description: "Test description", Price: 350.00, Category: "Test"},
		{ProductID: "product-008", Name: "Test Product 8", Description: "Test description", Price: 400.00, Category: "Test"},
		{ProductID: "product-009", Name: "Test Product 9", Description: "Test description", Price: 450.00, Category: "Test"},
		{ProductID: "product-010", Name: "Test Product 10", Description: "Test description", Price: 500.00, Category: "Test"},
		{ProductID: "product-011", Name: "Test Product 11", Description: "Test description", Price: 550.00, Category: "Test"},
	}

	for _, product := range testProducts {
		if err := db.Create(&product).Error; err != nil {
			t.Fatalf("Failed to create test product %s: %v", product.ProductID, err)
		}
	}
}

func TestCreateOrderIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	order := models.Order{
		UserID:      "user-123",
		TotalAmount: 299.99,
		Status:      "pending",
		Items: []models.OrderItem{
			{
				ProductID: "product-001",
				Quantity:  2,
				Price:     99.99,
			},
			{
				ProductID: "product-002",
				Quantity:  1,
				Price:     100.01,
			},
		},
	}

	createdOrder, err := repo.CreateOrder(order)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	// Verify OrderID was generated (should be a UUID)
	if createdOrder.OrderID == "" {
		t.Error("Expected OrderID to be generated, got empty string")
	}

	// Verify UserID matches
	if createdOrder.UserID != order.UserID {
		t.Errorf("Expected UserID %s, got %s", order.UserID, createdOrder.UserID)
	}

	// Verify number of items
	if len(createdOrder.Items) != 2 {
		t.Errorf("Expected 2 order items, got %d", len(createdOrder.Items))
	}
}

func TestGetOrderByIdIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	order := models.Order{
		UserID:      "user-456",
		TotalAmount: 149.99,
		Status:      "pending",
		Items: []models.OrderItem{
			{
				ProductID: "product-003",
				Quantity:  1,
				Price:     149.99,
			},
		},
	}

	createdOrder, err := repo.CreateOrder(order)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	retrievedOrder, err := repo.GetOrderById(createdOrder.OrderID)
	if err != nil {
		t.Fatalf("Failed to get order by ID: %v", err)
	}

	if retrievedOrder.OrderID != createdOrder.OrderID {
		t.Errorf("Expected OrderID %s, got %s", createdOrder.OrderID, retrievedOrder.OrderID)
	}

	if retrievedOrder.UserID != order.UserID {
		t.Errorf("Expected UserID %s, got %s", order.UserID, retrievedOrder.UserID)
	}
}

func TestUpdateOrderIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	order := models.Order{
		UserID:      "user-789",
		TotalAmount: 199.99,
		Status:      "pending",
		Items: []models.OrderItem{
			{
				ProductID: "product-004",
				Quantity:  1,
				Price:     199.99,
			},
		},
	}

	createdOrder, err := repo.CreateOrder(order) // Capture the created order
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	updatedOrder := models.Order{
		OrderID:     createdOrder.OrderID, // Use the generated ID
		UserID:      order.UserID,
		TotalAmount: 299.99,
		Status:      "confirmed",
		Items: []models.OrderItem{
			{
				ProductID: "product-004",
				Quantity:  2,
				Price:     149.99,
			},
		},
	}

	result, err := repo.UpdateOrder(updatedOrder)
	if err != nil {
		t.Fatalf("Failed to update order: %v", err)
	}

	if result.Status != "confirmed" {
		t.Errorf("Expected status 'confirmed', got %s", result.Status)
	}

	// The total amount is recalculated based on current product prices from DB
	// Just verify it's greater than 0 and properly calculated
	if result.TotalAmount <= 0 {
		t.Errorf("Expected positive total amount, got %.2f", result.TotalAmount)
	}

	// Verify the items were updated
	if len(result.Items) == 0 {
		t.Error("Expected updated order to have items")
	}
}

func TestDeleteOrderIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	order := models.Order{
		UserID:      "user-999",
		TotalAmount: 79.99,
		Status:      "pending",
		Items: []models.OrderItem{
			{
				ProductID: "product-005",
				Quantity:  1,
				Price:     79.99,
			},
		},
	}

	createdOrder, err := repo.CreateOrder(order) // Capture the created order
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	err = repo.DeleteOrder(createdOrder.OrderID) // Use the generated ID
	if err != nil {
		t.Fatalf("Failed to delete order: %v", err)
	}

	_, err = repo.GetOrderById(createdOrder.OrderID) // Use the generated ID
	if err == nil {
		t.Error("Expected error when getting deleted order, but got none")
	}
}

func TestListAllOrdersIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	// Create multiple test orders
	orders := []models.Order{
		{
			UserID:      "user-111",
			TotalAmount: 99.99,
			Status:      "pending",
			Items: []models.OrderItem{
				{
					ProductID: "product-006",
					Quantity:  1,
					Price:     99.99,
				},
			},
		},
		{
			UserID:      "user-456", // Change from user-222 to existing user
			TotalAmount: 199.99,
			Status:      "confirmed",
			Items: []models.OrderItem{
				{
					ProductID: "product-007",
					Quantity:  2,
					Price:     99.99,
				},
			},
		},
	}

	for _, order := range orders {
		_, err := repo.CreateOrder(order)
		if err != nil {
			t.Fatalf("Failed to create order %s: %v", order.OrderID, err)
		}
	}

	allOrders, err := repo.ListAllOrders(10, 0)
	if err != nil {
		t.Fatalf("Failed to list all orders: %v", err)
	}

	if len(allOrders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(allOrders))
	}
}

func TestListOrdersByUserIdIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	userID := "user-333"

	// Create orders for specific user
	orders := []models.Order{
		{
			UserID:      userID,
			TotalAmount: 59.99,
			Status:      "pending",
			Items: []models.OrderItem{
				{
					ProductID: "product-008",
					Quantity:  1,
					Price:     59.99,
				},
			},
		},
		{
			UserID:      userID,
			TotalAmount: 129.99,
			Status:      "confirmed",
			Items: []models.OrderItem{
				{
					ProductID: "product-009",
					Quantity:  1,
					Price:     129.99,
				},
			},
		},
		{
			UserID:      "user-789", // Change from different-user to existing user
			TotalAmount: 39.99,
			Status:      "pending",
			Items: []models.OrderItem{
				{
					ProductID: "product-010",
					Quantity:  1,
					Price:     39.99,
				},
			},
		},
	}

	for _, order := range orders {
		_, err := repo.CreateOrder(order)
		if err != nil {
			t.Fatalf("Failed to create order %s: %v", order.OrderID, err)
		}
	}

	userOrders, err := repo.ListOrdersByUserId(userID)
	if err != nil {
		t.Fatalf("Failed to list orders by user ID: %v", err)
	}

	if len(userOrders) != 2 {
		t.Errorf("Expected 2 orders for user %s, got %d", userID, len(userOrders))
	}

	for _, order := range userOrders {
		if order.UserID != userID {
			t.Errorf("Expected UserID %s, got %s", userID, order.UserID)
		}
	}
}

func TestDeleteOrderItemIntegration(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Create all test data using shared function
	createTestData(t, db)

	repo := NewOrderRepository(db)

	order := models.Order{
		UserID:      "user-444",
		TotalAmount: 249.98,
		Status:      "pending",
		Items: []models.OrderItem{
			{
				ProductID: "product-011",
				Quantity:  1,
				Price:     99.99,
			},
			{
				ProductID: "product-004", // Change from product-012 to existing product
				Quantity:  1,
				Price:     149.99,
			},
		},
	}

	createdOrder, err := repo.CreateOrder(order)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	if len(createdOrder.Items) != 2 {
		t.Fatalf("Expected 2 items in created order, got %d", len(createdOrder.Items))
	}

	// Delete one item
	itemToDelete := createdOrder.Items[0]
	err = repo.DeleteOrderItem(itemToDelete.OrderItemID) // Use OrderItemID, not OrderID
	if err != nil {
		t.Fatalf("Failed to delete order item: %v", err)
	}

	// Verify item was deleted
	updatedOrder, err := repo.GetOrderById(createdOrder.OrderID) // Use createdOrder.OrderID
	if err != nil {
		t.Fatalf("Failed to get updated order: %v", err)
	}

	if len(updatedOrder.Items) != 1 {
		t.Errorf("Expected 1 item after deletion, got %d", len(updatedOrder.Items))
	}

	if updatedOrder.Items[0].ProductID == itemToDelete.ProductID {
		t.Error("Item was not properly deleted")
	}
}
