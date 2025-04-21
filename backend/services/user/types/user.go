package types

import (
	"time"
)

// DB types
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"isActive" gorm:"default:true"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// API types
type UserSearchParams struct {
	ID    *string `json:"id"`
	Email *string `json:"email"`
}
