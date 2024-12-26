// internal/usecase/user_usecase.go
package usecase

import (
	"errors"
	"fmt"

	"github.com/almatkai/book-exchange-backend/internal/entity"
	"github.com/almatkai/book-exchange-backend/internal/repository"
	"github.com/almatkai/book-exchange-backend/pkg/utils" // Import the utils package
)

var (
	ErrUsernameExists     = repository.ErrUsernameExists
	ErrEmailExists        = repository.ErrEmailExists
	ErrInvalidCredentials = repository.ErrInvalidCredentials
)

type UserUseCase interface {
	RegisterUser(username, email, password string) (*entity.User, error)
	LoginUser(username, password string) (string, error) // Returns JWT token on success
}

type userUseCase struct {
	userRepo repository.UserRepository
	jwtSvc   utils.JWTService
}

func NewUserUseCase(userRepo repository.UserRepository, jwtSvc utils.JWTService) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
		jwtSvc:   jwtSvc,
	}
}

func (uc *userUseCase) RegisterUser(username, email, password string) (*entity.User, error) {
	// Create new user entity
	user := &entity.User{
		Username: username,
		Email:    email,
		Password: password, // Plain password; will be hashed in repository
		// Initialize other fields as needed
		IsActive: true,
	}

	// Create user in repository (this will handle hashing and uniqueness)
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *userUseCase) LoginUser(username, password string) (string, error) {
	user, err := uc.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	// Compare password
	if err := utils.CheckPasswordHash(password, user.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	// Update last login
	if err := uc.userRepo.UpdateLastLogin(user.UserID); err != nil {
		// Log the error but don't prevent login
		fmt.Printf("Error updating last login: %v\n", err)
	}

	// Generate JWT token
	token, err := uc.jwtSvc.GenerateToken(fmt.Sprintf("%d", user.UserID))
	if err != nil {
		return "", err
	}

	return token, nil
}
