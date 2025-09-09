package types

import "time"

type AuthProvider struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"user_id"`
	Provider   string    `json:"provider"`
	ProviderID string    `json:"provider_id"`
	CreatedAt  time.Time `json:"created_at" gorm:"default:now()"`
}
