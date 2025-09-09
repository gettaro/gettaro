package types

import (
	"time"
)

// Title represents a job title within an organization
type Title struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organization_id"`
	IsManager      bool      `json:"is_manager" gorm:"default:false"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// MemberTitle represents a member's title assignment within an organization
type MemberTitle struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MemberID       string    `json:"member_id"`
	TitleID        string    `json:"title_id"`
	OrganizationID string    `json:"organization_id"`
	CreatedAt      time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updated_at"`
}
