package types

import (
	"time"
)

type Organization struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug" gorm:"uniqueIndex"`
	IsOwner   bool      `json:"is_owner" gorm:"-"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrganizationRequest struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
