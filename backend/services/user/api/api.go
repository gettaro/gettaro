package api

import (
	"ems.dev/backend/services/user/database"
)

type Api struct {
	db *database.UserDB
}

func NewApi(userDb *database.UserDB) *Api {
	return &Api{
		db: userDb,
	}
}
