// internal/entity/exchange.go
package entity

import (
	"time"

	"github.com/google/uuid"
	_ "gorm.io/gorm"
)

type Exchange struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	PostID       uuid.UUID `gorm:"type:uuid;not null" json:"post_id"`
	RequesterID  uuid.UUID `gorm:"type:uuid;not null" json:"requester_id"`
	OwnerID      uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	Status       string    `gorm:"type:varchar(10);check:status IN ('pending','accepted','rejected','completed');not null" json:"status"`
	Location     string    `gorm:"type:varchar(255);not null" json:"location"`
	ExchangeDate time.Time `gorm:"type:timestamp;not null" json:"exchange_date"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Post      Post      `gorm:"foreignKey:PostID" json:"post,omitempty"`
	Requester User      `gorm:"foreignKey:RequesterID" json:"requester,omitempty"`
	Owner     User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Messages  []Message `gorm:"foreignKey:ExchangeID" json:"messages,omitempty"`
	Ratings   []Rating  `gorm:"foreignKey:ExchangeID" json:"ratings,omitempty"`
}
