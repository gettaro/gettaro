package types

import (
	"time"

	"gorm.io/datatypes"
)

// DB types
type User struct {
	ID            string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email         string         `json:"email" gorm:"uniqueIndex"`
	Name          string         `json:"name"`
	IsActive      bool           `json:"isActive" gorm:"default:true"`
	Status        string         `json:"status"`
	TitleID       *string        `json:"titleId"`
	Organizations []Organization `json:"organizations" gorm:"many2many:user_organizations;"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"default:now()"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}

type Title struct {
	ID   string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string `json:"name"`
}

type AuthProvider struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"userId"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"providerId"`
	CreatedAt  time.Time `json:"createdAt" gorm:"default:now()"`
}

type Role struct {
	ID          string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string         `json:"name"`
	Permissions datatypes.JSON `json:"permissions"`
}

// Organization represents the organization model
type Organization struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug" gorm:"uniqueIndex"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updatedAt"`
	IsOwner   bool      `json:"isOwner" gorm:"-"`
}

// API types
type UserSearchParams struct {
	ID    *string `json:"id"`
	Email *string `json:"email"`
}
