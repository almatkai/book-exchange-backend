package repository

import (
	"errors"
	"github.com/almatkai/book-exchange-backend/pkg/utils"
	"gorm.io/gorm"

	"github.com/almatkai/book-exchange-backend/internal/entity"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
)

// UserRepository defines methods for user data persistence
type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	UpdateLastLogin(userID int) error
}

// GormUserRepository is a GORM implementation of UserRepository
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GormUserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{db: db}
}

// Create inserts a new user into the database
func (repo *GormUserRepository) Create(user *entity.User) error {
	// Check if username already exists
	var existingUser entity.User
	if err := repo.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return ErrUsernameExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Check if email already exists
	if err := repo.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return ErrEmailExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Create the user
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

// FindByUsername retrieves a user by username
func (repo *GormUserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := repo.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (repo *GormUserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := repo.db.Where("email = ? AND is_active = ?", email, true).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateLastLogin updates the last_login timestamp for a user
func (repo *GormUserRepository) UpdateLastLogin(userID int) error {
	return repo.db.Model(&entity.User{}).Where("user_id = ?", userID).Update("last_login", gorm.Expr("NOW()")).Error
}
