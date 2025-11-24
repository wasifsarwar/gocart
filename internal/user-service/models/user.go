package models

import "time"

type User struct {
	UserID       string    `gorm:"primaryKey" json:"user_id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	Password     string    `gorm:"-" json:"password,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
