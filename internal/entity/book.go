// internal/entity/book.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Author      string    `gorm:"type:varchar(255);not null" json:"author"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	IsAvailable bool      `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User  User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Posts []Post `gorm:"foreignKey:BookID" json:"posts,omitempty"`
}
