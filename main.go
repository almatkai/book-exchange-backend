package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/almatkai/book-exchange-backend/internal/config"
	"github.com/almatkai/book-exchange-backend/internal/delivery/router"
	"github.com/almatkai/book-exchange-backend/internal/delivery/router/handlers" // Alias to `handlers`
	"github.com/almatkai/book-exchange-backend/internal/repository"
	"github.com/almatkai/book-exchange-backend/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Repositories and Use Cases
	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Handlers
	userHandler := handlers.NewUserHandler(userUseCase)

	// Initialize Router
	newRouter := router.NewRouter(userHandler)

	// Start Server with dynamic port from config
	port := cfg.ServerPort
	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), newRouter))
}

//book-exchange-backend/
//├── cmd/
//│   └── server/
//│       └── main.go
//├── internal/
//│   ├── entity/
//│   │   └── user.go
//│   │   └── book.go
//│   │   └── post.go
//│   │   └── exchange.go
//│   │   └── message.go
//│   │   └── rating.go
//│   ├── server/
//│   │   └── database.go
//│   ├── usecase/
//│   │   └── user_usecase.go
//│   │   └── book_usecase.go
//│   │   └── post_usecase.go
//│   │   └── exchange_usecase.go
//│   │   └── message_usecase.go
//│   │   └── rating_usecase.go
//│   ├── repository/
//│   │   └── user_repository.go
//│   │   └── book_repository.go
//│   │   └── post_repository.go
//│   │   └── exchange_repository.go
//│   │   └── message_repository.go
//│   │   └── rating_repository.go
//│   ├── delivery/
//│   │   └── http/
//│   │       └── handlers/
//│   │           └── user_handler.go
//│   │           └── book_handler.go
//│   │           └── post_handler.go
//│   │           └── exchange_handler.go
//│   │           └── message_handler.go
//│   │           └── rating_handler.go
//│   │       └── router.go
//│   ├── config/
//│   │   └── config.go
//│   └── middleware/
//│       └── auth.go
//├── pkg/
//│   └── utils/
//│       └── jwt.go
//├── go.mod
//├── go.sum
//└── Dockerfile