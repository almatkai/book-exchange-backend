// internal/delivery/router/handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/almatkai/book-exchange-backend/internal/entity"
	"github.com/almatkai/book-exchange-backend/internal/usecase"

	"unicode"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with username, email, and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body entity.UserCredentials true "User Credentials"
// @Success 201 {object} entity.User
// @Failure 400 {string} string "Invalid input"
// @Failure 409 {string} string "Username or email already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds entity.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Input validation
	if creds.Username == "" || creds.Email == "" || creds.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Password strength validation
	if !isValidPassword(creds.Password) {
		http.Error(w, "Password does not meet strength requirements", http.StatusBadRequest)
		return
	}

	// Register user
	user, err := h.userUseCase.RegisterUser(creds.Username, creds.Email, creds.Password)
	if err != nil {
		switch err {
		case usecase.ErrUsernameExists:
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		case usecase.ErrEmailExists:
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		default:
			http.Error(w, "Failed to register user", http.StatusInternalServerError)
			return
		}
	}

	// Exclude password from response
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user and return a JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param user body entity.UserCredentials true "User Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Invalid username or password"
// @Failure 500 {string} string "Internal server error"
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds entity.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Input validation
	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Login user
	token, err := h.userUseCase.LoginUser(creds.Username, creds.Password)
	if err != nil {
		if err == usecase.ErrInvalidCredentials {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Failed to login", http.StatusInternalServerError)
		return
	}

	// Return token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// Utility functions for password validation

func isValidPassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
