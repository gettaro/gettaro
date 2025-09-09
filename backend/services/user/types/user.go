package types

import (
	"time"
)

// DB types
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
}

// API types
type UserSearchParams struct {
	ID    *string `json:"id"`
	Email *string `json:"email"`
}
