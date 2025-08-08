package repository

import (
	"errors"
	"gocart/internal/user-service/models"
	"strings"
	"time"

	"github.com/google/uuid"
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
	user.UserID = uuid.New().String()
	if err := r.db.Create(&user).Error; err != nil {
		// Check for duplicate email error
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "duplicate entry") {
			return models.User{}, errors.New("user with this email already exists")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserById(userID string) (models.User, error) {
	var user models.User
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) DeleteUser(userID string) (models.User, error) {
	var user models.User
	// First get the user to return it
	if err := r.db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}
	// Then delete it
	if err := r.db.Where("user_id = ?", userID).Delete(&models.User{}).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	user.UpdatedAt = time.Now()
	if err := r.db.Model(&models.User{}).Where("user_id = ?", user.UserID).Updates(&user).Error; err != nil {
		// Check for duplicate email error
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") ||
			strings.Contains(err.Error(), "duplicate entry") {
			return models.User{}, errors.New("user with this email already exists")
		}
		return models.User{}, err
	}
	return user, nil
}

func (r *userRepository) ListAllUsers() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
