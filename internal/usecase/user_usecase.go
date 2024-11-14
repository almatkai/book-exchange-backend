// internal/usecase/user_usecase.go
package usecase

import (
	"errors"
	"github.com/google/uuid"
	"time"

	"github.com/almatkai/book-exchange-backend/internal/entity"
	"github.com/almatkai/book-exchange-backend/internal/repository"
	"github.com/almatkai/book-exchange-backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserUseCase interface defines methods for user business logic
type UserUseCase interface {
	RegisterUser(username, email, password string) (*entity.User, error)
	LoginUser(username, password string) (string, error) // Returns JWT token on success
}

// userUseCase implements UserUseCase interface
type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo}
}

// RegisterUser creates a new user with hashed password
func (u *userUseCase) RegisterUser(username, email, password string) (*entity.User, error) {
	// Check if user already exists
	existingUser, _ := u.userRepo.GetUserByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user entity
	user := &entity.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save user in database
	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser authenticates user and returns JWT if successful
func (u *userUseCase) LoginUser(username, password string) (string, error) {
	user, err := u.userRepo.GetUserByUsername(username)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("incorrect password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String())
	if err != nil {
		return "", err
	}
	return token, nil
}
