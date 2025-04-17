package main

import (
	"log"

	"ems.dev/backend/database"
	"ems.dev/backend/http/server"
	orgapi "ems.dev/backend/services/organization/api"
	orgdb "ems.dev/backend/services/organization/database"
	teamapi "ems.dev/backend/services/team/api"
	teamdb "ems.dev/backend/services/team/database"
	userapi "ems.dev/backend/services/user/api"
	userdb "ems.dev/backend/services/user/database"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	database.InitDB()

	// Initialize services
	userDb := userdb.NewUserDB()
	userApi := userapi.NewApi(userDb)
	orgDb := orgdb.NewOrganizationDB()
	orgApi := orgapi.NewApi(orgDb)
	teamDb := teamdb.NewTeamDB()
	teamApi := teamapi.NewApi(teamDb)

	// Initialize and run server
	srv := server.New(database.DB, userApi, orgApi, teamApi)
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
