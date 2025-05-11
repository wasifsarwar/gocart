package models

type Product struct {
	ProductID   string  `gorm:"primaryKey" json:"product_id"`
	Name        string  `gorm:"not null" json:"name"`
	Description string  `json:"description"`
	Price       float64 `gorm:"not null" json:"price"`
	Category    string  `json:"category"`
}
