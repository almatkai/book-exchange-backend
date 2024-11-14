// internal/entity/message.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ExchangeID uuid.UUID `gorm:"type:uuid;not null" json:"exchange_id"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	SentAt     time.Time `gorm:"autoCreateTime" json:"sent_at"`

	// Relationships
	Exchange Exchange `gorm:"foreignKey:ExchangeID" json:"exchange,omitempty"`
	Sender   User     `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}
