package types

import (
	"time"

	"gorm.io/datatypes"
)

type User struct {
	ID             string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email          string `gorm:"uniqueIndex"`
	Name           string
	IsActive       bool `gorm:"default:true"`
	Status         string
	TitleID        *string
	OrganizationID *string
	CreatedAt      time.Time `gorm:"default:now()"`
	UpdatedAt      time.Time
}

type Title struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string
}

type AuthProvider struct {
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string
	Provider   string
	ProviderID string
	CreatedAt  time.Time `gorm:"default:now()"`
}

type Role struct {
	ID          string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name        string
	Permissions datatypes.JSON
}
