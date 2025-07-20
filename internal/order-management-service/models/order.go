package models

import "time"

type Order struct {
	OrderID     string      `gorm:"primaryKey;type:uuid" json:"order_id"`
	UserID      string      `gorm:"not null" json:"user_id"`
	Status      string      `gorm:"not null" json:"status"`
	TotalAmount float64     `gorm:"not null" json:"total_amount"`
	Items       []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"` // 1:N relationship, cascade delete
	CreatedAt   time.Time   `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"not null" json:"updated_at"`
}

type OrderItem struct {
	OrderItemID string    `gorm:"primaryKey;type:uuid" json:"order_item_id"`
	OrderID     string    `gorm:"index" json:"order_id"`
	ProductID   string    `gorm:"not null" json:"product_id"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}
