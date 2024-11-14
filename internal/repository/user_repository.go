package repository

import (
	"errors"

	"github.com/almatkai/book-exchange-backend/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository interface defines methods for user operations
type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByID(id uuid.UUID) (*entity.User, error)
}

// userRepository implements UserRepository interface
type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
