// internal/entity/user.go
package entity

import (
	"time"
)

type User struct {
	UserID    int        `gorm:"primaryKey;column:user_id" json:"id"`
	Username  string     `gorm:"unique;not null;size:50" json:"username"`
	Email     string     `gorm:"unique;not null;size:255" json:"email"`
	Password  string     `gorm:"column:password_hash;not null;size:255" json:"-"`
	FullName  string     `gorm:"size:100" json:"full_name,omitempty"`
	Location  string     `gorm:"size:255" json:"location,omitempty"`
	Bio       string     `gorm:"type:text" json:"bio,omitempty"`
	Rating    float32    `gorm:"type:decimal(3,2);default:0.00;check:rating >= 0 AND rating <= 5" json:"rating,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}
