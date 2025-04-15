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

type AuthProvider struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"userId"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"providerId"`
	CreatedAt  time.Time `json:"createdAt" gorm:"default:now()"`
}

// API types
type UserSearchParams struct {
	ID    *string `json:"id"`
	Email *string `json:"email"`
}
