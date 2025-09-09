package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"ems.dev/backend/database"
	"ems.dev/backend/http/server"
	"ems.dev/backend/jobs/scheduler"
	"ems.dev/backend/jobs/sourcecontrol"
	scprovider "ems.dev/backend/jobs/sourcecontrol/providers"
	githubprovider "ems.dev/backend/jobs/sourcecontrol/providers/github"
	auth0client "ems.dev/backend/libraries/auth0"
	"ems.dev/backend/libraries/github"
	authapi "ems.dev/backend/services/auth/api"
	authdb "ems.dev/backend/services/auth/database"
	conversationtemplateapi "ems.dev/backend/services/conversationtemplate/api"
	conversationtemplatedb "ems.dev/backend/services/conversationtemplate/database"
	directsapi "ems.dev/backend/services/directs/api"
	directsdb "ems.dev/backend/services/directs/database"
	integrationapi "ems.dev/backend/services/integration/api"
	integrationdb "ems.dev/backend/services/integration/database"
	memberapi "ems.dev/backend/services/member/api"
	memberdb "ems.dev/backend/services/member/database"
	orgapi "ems.dev/backend/services/organization/api"
	orgdb "ems.dev/backend/services/organization/database"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	sourcecontroldb "ems.dev/backend/services/sourcecontrol/database"
	teamapi "ems.dev/backend/services/team/api"
	teamdb "ems.dev/backend/services/team/database"
	titleapi "ems.dev/backend/services/title/api"
	titleDb "ems.dev/backend/services/title/database"
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
	authDb := authdb.New(database.DB)
	authApi := authapi.NewApi(auth0Client, authDb)
	integrationDb := integrationdb.NewIntegrationDB(database.DB)
	integrationApi := integrationapi.NewApi(integrationDb, []byte("QI$Pi!<Jc@L<%bwI"))
	sourcecontrolDb := sourcecontroldb.NewSourceControlDB(database.DB)
	sourcecontrolApi := sourcecontrolapi.NewAPI(sourcecontrolDb)
	titleDb := titleDb.NewTitleDB(database.DB)
	titleApi := titleapi.NewApi(titleDb)
	orgDb := orgdb.NewOrganizationDB(database.DB)
	orgApi := orgapi.NewApi(orgDb, userApi, titleApi, sourcecontrolApi)
	directsDb := directsdb.NewDirectReportsDB(database.DB)
	directsApi := directsapi.NewDirectReportsAPI(directsDb)
	memberDb := memberdb.NewMemberDB(database.DB)
	memberApi := memberapi.NewApi(memberDb, userApi, sourcecontrolApi, titleApi, directsApi)
	teamDb := teamdb.NewTeamDB(database.DB)
	teamApi := teamapi.NewApi(teamDb, orgApi)
	conversationTemplateDb := conversationtemplatedb.NewConversationTemplateDatabase(database.DB)
	conversationTemplateApi := conversationtemplateapi.NewConversationTemplateAPI(conversationTemplateDb)

	// Initialize and start sync job scheduler
	// Check if jobs are enabled
	if os.Getenv("JOBS_ENABLED") == "true" {
		log.Println("Jobs are enabled")
		syncInterval := getSyncInterval()
		githubProvider := githubprovider.NewProvider(github.NewClient(), integrationApi, sourcecontrolApi)
		scProviderFactory := scprovider.NewFactory([]scprovider.SourceControlProvider{githubProvider})
		syncJob := sourcecontrol.NewSyncJob(integrationApi, orgApi, scProviderFactory)
		//go syncJob.Run(context.Background())
		scheduler := scheduler.NewScheduler(syncJob, syncInterval)
		go scheduler.Start(context.Background())
	}

	// Initialize and run server
	srv := server.New(database.DB, userApi, orgApi, teamApi, titleApi, authApi, integrationApi, sourcecontrolApi, memberApi, directsApi, conversationTemplateApi)
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getSyncInterval() time.Duration {
	intervalStr := os.Getenv("SYNC_INTERVAL_HOURS")
	if intervalStr == "" {
		intervalStr = "4" // Default to 4 hours
	}

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Printf("Invalid SYNC_INTERVAL_HOURS value, using default 4 hours")
		interval = 4
	}

	return time.Duration(interval) * time.Hour
}
