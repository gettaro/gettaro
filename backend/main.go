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
	apiai "ems.dev/backend/services/ai/api"
	aidb "ems.dev/backend/services/ai/database"
	"ems.dev/backend/services/ai/providers"
	anthropicprovider "ems.dev/backend/services/ai/providers/anthropic"
	groqprovider "ems.dev/backend/services/ai/providers/groq"
	aitypes "ems.dev/backend/services/ai/types"
	authapi "ems.dev/backend/services/auth/api"
	authdb "ems.dev/backend/services/auth/database"
	conversationapi "ems.dev/backend/services/conversation/api"
	conversationdb "ems.dev/backend/services/conversation/database"
	conversationtemplateapi "ems.dev/backend/services/conversationtemplate/api"
	conversationtemplatedb "ems.dev/backend/services/conversationtemplate/database"
	directsapi "ems.dev/backend/services/directs/api"
	directsdb "ems.dev/backend/services/directs/database"
	integrationapi "ems.dev/backend/services/integration/api"
	integrationdb "ems.dev/backend/services/integration/database"
	memberapi "ems.dev/backend/services/member/api"
	memberdb "ems.dev/backend/services/member/database"
	metricsapi "ems.dev/backend/services/metrics/api"
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
	metricsApi := metricsapi.NewApi(memberApi, teamApi, sourcecontrolApi)
	conversationTemplateDb := conversationtemplatedb.NewConversationTemplateDatabase(database.DB)
	conversationTemplateApi := conversationtemplateapi.NewConversationTemplateAPI(conversationTemplateDb)
	conversationDb := conversationdb.NewConversationDB(database.DB)
	conversationApi := conversationapi.NewConversationAPI(conversationDb, conversationTemplateApi)
	aiDb := aidb.NewAIDB(database.DB)

	// Create AI service configuration
	aiConfig := &aitypes.AIServiceConfig{
		Provider:          getEnvOrDefault("AI_PROVIDER", "anthropic"),
		Model:             getEnvOrDefault("AI_MODEL", "claude-3-5-sonnet-20241022"),
		MaxTokens:         100000,
		Temperature:       0.7,
		MaxContextSize:    500000,
		EnableHistory:     true,
		EnableSuggestions: true,
	}

	// Initialize AI provider based on configuration
	aiProvider := initializeAIProvider(aiConfig)
	aiApi := apiai.NewAIService(aiDb, memberApi, teamApi, conversationApi, sourcecontrolApi, orgApi, userApi, aiConfig, aiProvider)

	// Initialize and start sync job scheduler
	// Check if jobs are enabled
	if os.Getenv("JOBS_ENABLED") == "true" {
		log.Println("Jobs are enabled")
		syncInterval := getSyncInterval()
		githubProvider := githubprovider.NewProvider(github.NewClient(), integrationApi, sourcecontrolApi, teamApi)
		scProviderFactory := scprovider.NewFactory([]scprovider.SourceControlProvider{githubProvider})
		syncJob := sourcecontrol.NewSyncJob(integrationApi, orgApi, scProviderFactory)
		//go syncJob.Run(context.Background())
		scheduler := scheduler.NewScheduler(syncJob, syncInterval)
		go scheduler.Start(context.Background())
	}

	// Initialize and run server
	srv := server.New(database.DB, userApi, orgApi, teamApi, titleApi, authApi, integrationApi, sourcecontrolApi, memberApi, metricsApi, directsApi, conversationTemplateApi, conversationApi, aiApi)
	if err := srv.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initializeAIProvider initializes the appropriate AI provider based on configuration
func initializeAIProvider(config *aitypes.AIServiceConfig) providers.AIProviderInterface {
	provider := getEnvOrDefault("AI_PROVIDER", "anthropic")

	switch provider {
	case "groq":
		groqAPIKey := getEnvOrDefault("GROQ_API_KEY", "")
		if groqAPIKey == "" {
			log.Fatal("GROQ_API_KEY environment variable is required when using Groq provider")
		}
		log.Printf("Initializing Groq AI provider with model: %s", config.Model)
		return groqprovider.NewProvider(groqAPIKey)
	case "anthropic":
		fallthrough
	default:
		anthropicAPIKey := getEnvOrDefault("ANTHROPIC_API_KEY", "")
		if anthropicAPIKey == "" {
			log.Fatal("ANTHROPIC_API_KEY environment variable is required when using Anthropic provider")
		}
		log.Printf("Initializing Anthropic AI provider with model: %s", config.Model)
		return anthropicprovider.NewProvider(anthropicAPIKey)
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

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
