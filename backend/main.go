package main

import (
	"log"
	"os"

	"ems.dev/backend/database"
	"ems.dev/backend/http/server"
	auth0client "ems.dev/backend/libraries/auth0"
	authapi "ems.dev/backend/services/auth/api"
	authdb "ems.dev/backend/services/auth/database"
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

	// Initialize Auth0 client
	auth0Client := auth0client.NewClient(
		os.Getenv("AUTH0_AUTHORITY"),
		os.Getenv("AUTH0_CLIENT_ID"),
	)

	// Initialize services
	userDb := userdb.NewUserDB(database.DB)
	userApi := userapi.NewApi(userDb)
	orgDb := orgdb.NewOrganizationDB(database.DB)
	orgApi := orgapi.NewApi(orgDb, userApi)
	teamDb := teamdb.NewTeamDB(database.DB)
	teamApi := teamapi.NewApi(teamDb, orgApi)
	authDb := authdb.New(database.DB)
	authApi := authapi.NewApi(auth0Client, authDb)

	// Initialize and run server
	srv := server.New(database.DB, userApi, orgApi, teamApi, authApi)
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
