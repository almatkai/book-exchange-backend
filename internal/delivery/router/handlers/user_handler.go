// internal/delivery/router/handlers/user_handler.go
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/almatkai/book-exchange-backend/internal/entity"
	"github.com/almatkai/book-exchange-backend/internal/usecase"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

// Register endpoint
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var creds entity.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	user, err := h.userUseCase.RegisterUser(creds.Username, creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login endpoint
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds entity.UserCredentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := h.userUseCase.LoginUser(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
