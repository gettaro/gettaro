package auth

import "context"

// UserClaims represents the claims we expect from Auth0
type UserClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Provider string `json:"provider"`
}

// CustomClaims represents the custom claims we expect from Auth0
type CustomClaims struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// Validate implements the validator.CustomClaims interface
func (c CustomClaims) Validate(ctx context.Context) error {
	return nil
}
