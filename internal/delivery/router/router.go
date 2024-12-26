// internal/delivery/router/router.go
package router

import (
	"net/http"

	"github.com/almatkai/book-exchange-backend/internal/delivery/router/handlers"
	"github.com/almatkai/book-exchange-backend/internal/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(userHandler *handlers.UserHandler, jwtKey []byte) *mux.Router {
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Book Exchange API"))
	}).Methods(http.MethodGet)

	// Protected routes
	protected := router.PathPrefix("/protected").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(next, jwtKey)
	})
	protected.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is protected route"))
	}).Methods(http.MethodGet)

	return router
}

//curl -X POST http://localhost:8000/register -H "Content-Type: application/json" -d "{\"username\": \"almatKAI\", \"email\": \"almatkai@example.com\", \"password\": \"password123\"}"
//
//curl -X POST http://localhost:8000/register -H "Content-Type: application/json" -d "{\"username\": \"almatKAI\", \"email\": \"almatkai@example.com\", \"password\": \"password123\"}"
//
//curl -X POST http://localhost:8000/login -H "Content-Type: application/json" -d "{\"username\": \"almatKAI\", \"email\": \"almatkai@example.com\", \"password\": \"password123\"}"

//curl -X POST http://localhost:8080/register \
//-H "Content-Type: application/json" \
//-d '{
//"username": "john_doe",
//"email": "john@example.com",
//"password": "SecureP@ssw0rd!"
//}'
