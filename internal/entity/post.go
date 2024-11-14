// internal/entity/post.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	BookID         uuid.UUID  `gorm:"type:uuid;not null" json:"book_id"`
	ExchangeType   string     `gorm:"type:varchar(10);check:exchange_type IN ('permanent','temporary');not null" json:"exchange_type"`
	AvailableUntil *time.Time `gorm:"type:timestamp" json:"available_until,omitempty"`
	Location       string     `gorm:"type:varchar(255);not null" json:"location"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book      Book       `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Exchanges []Exchange `gorm:"foreignKey:PostID" json:"exchanges,omitempty"`
}
