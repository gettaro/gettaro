package api

import (
	"ems.dev/backend/services/team/database"
)

type Api struct {
	db *database.TeamDB
}

func NewApi(teamDb *database.TeamDB) *Api {
	return &Api{
		db: teamDb,
	}
}
