package api

import (
	"ems.dev/backend/services/organization/database"
)

type Api struct {
	db *database.OrganizationDB
}

func NewApi(orgDb *database.OrganizationDB) *Api {
	return &Api{
		db: orgDb,
	}
}
