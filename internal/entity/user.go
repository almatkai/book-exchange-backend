// internal/entity/user.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Username     string    `gorm:"type:varchar(50);unique;not null" json:"username"`
	Email        string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Books           []Book     `gorm:"foreignKey:UserID" json:"books,omitempty"`
	Posts           []Post     `gorm:"foreignKey:UserID" json:"posts,omitempty"`
	Exchanges       []Exchange `gorm:"foreignKey:OwnerID" json:"exchanges,omitempty"`
	Messages        []Message  `gorm:"foreignKey:SenderID" json:"messages,omitempty"`
	RatingsGiven    []Rating   `gorm:"foreignKey:RaterID" json:"ratings_given,omitempty"`
	RatingsReceived []Rating   `gorm:"foreignKey:RateeID" json:"ratings_received,omitempty"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
