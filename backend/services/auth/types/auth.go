package types

import "time"

type AuthProvider struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"userId"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"providerId"`
	CreatedAt  time.Time `json:"createdAt" gorm:"default:now()"`
}
