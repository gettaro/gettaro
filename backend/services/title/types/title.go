package types

import (
	"time"
)

// Title represents a job title within an organization
type Title struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name           string    `json:"name"`
	OrganizationID string    `json:"organizationId"`
	IsManager      bool      `json:"isManager" gorm:"default:false"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// MemberTitle represents a member's title assignment within an organization
type MemberTitle struct {
	ID             string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	MemberID       string    `json:"memberId"`
	TitleID        string    `json:"titleId"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt" gorm:"default:now()"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
