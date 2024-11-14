// internal/delivery/router/router.go
package router

import (
	"github.com/almatkai/book-exchange-backend/internal/delivery/router/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(userHandler *handlers.UserHandler) *mux.Router {
	router := mux.NewRouter()

	// User routes
	router.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Book Exchange API"))
	}).Methods(http.MethodGet)

	return router
}

//curl -X POST http://localhost:8000/register -H "Content-Type: application/json" -d "{\"username\": \"almatKAI\", \"email\": \"almatkai@example.com\", \"password\": \"password123\"}"
//
//curl -X POST http://localhost:8000/register -H "Content-Type: application/json" -d "{\"username\": \"almatKAI\", \"email\": \"almatkai@example.com\", \"password\": \"password123\"}"
//
//curl -X POST http://localhost:8000/login -H "Content-Type: application/json" -d "{\"username\": \"almatKAI\", \"email\": \"almatkai@example.com\", \"password\": \"password123\"}"
