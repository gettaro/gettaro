package main

import (
	"log"

	"ems.dev/backend/database"
	"ems.dev/backend/http/server"
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

	// Initialize and run server
	srv := server.New(database.DB, userApi)
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
