package user

import "ems.dev/backend/services/user/types"

type GetUserResponse struct {
	User *types.User `json:"user"`
}
