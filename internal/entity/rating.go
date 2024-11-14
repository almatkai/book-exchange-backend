// internal/entity/rating.go
package entity

import (
	"time"

	"github.com/google/uuid"
)

type Rating struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ExchangeID uuid.UUID `gorm:"type:uuid;not null" json:"exchange_id"`
	RaterID    uuid.UUID `gorm:"type:uuid;not null" json:"rater_id"`
	RateeID    uuid.UUID `gorm:"type:uuid;not null" json:"ratee_id"`
	Rating     int       `gorm:"type:integer;check:rating BETWEEN 1 AND 5" json:"rating"`
	Comment    string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Exchange Exchange `gorm:"foreignKey:ExchangeID" json:"exchange,omitempty"`
	Rater    User     `gorm:"foreignKey:RaterID" json:"rater,omitempty"`
	Ratee    User     `gorm:"foreignKey:RateeID" json:"ratee,omitempty"`
}
