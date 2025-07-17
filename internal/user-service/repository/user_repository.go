package repository

import (
	"gocart/internal/user-service/models"
	"time"

	db "gocart/pkg/db"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserById(userID string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	DeleteUser(userID string) (models.User, error)
	ListAllUsers() ([]models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	user.CreatedAt = time.Now()
	if err := db.DB.Create(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserById(userID string) (models.User, error) {
	var user models.User
	if err := db.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(userID string) (models.User, error) {
	var user models.User
	if err := db.DB.Where("user_id = ?", userID).Delete(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	user.UpdatedAt = time.Now()
	if err := db.DB.Model(&models.User{}).Where("user_id = ?", user.UserID).Updates(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) ListAllUsers() ([]models.User, error) {
	var users []models.User
	if err := db.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
