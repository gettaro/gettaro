package types

import (
	"time"
)

// Title represents a job title within an organization
type Title struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// UserTitle represents a user's title assignment within an organization
type UserTitle struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string    `json:"userId"`
	TitleID        string    `json:"titleId"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
