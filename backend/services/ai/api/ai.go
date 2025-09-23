package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"ems.dev/backend/services/ai/database"
	"ems.dev/backend/services/ai/types"
	conversationapi "ems.dev/backend/services/conversation/api"
	conversationtypes "ems.dev/backend/services/conversation/types"
	memberapi "ems.dev/backend/services/member/api"
	membertypes "ems.dev/backend/services/member/types"
	orgapi "ems.dev/backend/services/organization/api"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	teamapi "ems.dev/backend/services/team/api"
	teamtypes "ems.dev/backend/services/team/types"
	userapi "ems.dev/backend/services/user/api"
)

// AIServiceInterface defines the interface for AI operations
type AIServiceInterface interface {
	Query(ctx context.Context, req *types.AIQueryRequest, userID string) (*types.AIQueryResponse, error)
	GetQueryHistory(ctx context.Context, organizationID string, userID *string, limit int) ([]*types.AIQueryHistory, error)
	GetQueryStats(ctx context.Context, organizationID string, userID *string, days int) (*types.AIQueryStats, error)
}

// AIService implements the AI service
type AIService struct {
	db               *database.AIDB
	memberAPI        memberapi.MemberAPI
	teamAPI          teamapi.TeamAPI
	conversationAPI  conversationapi.ConversationAPIInterface
	sourceControlAPI sourcecontrolapi.SourceControlAPI
	organizationAPI  orgapi.OrganizationAPI
	userAPI          userapi.UserAPI
	config           *types.AIServiceConfig
}

// NewAIService creates a new AI service instance
func NewAIService(
	db *database.AIDB,
	memberAPI memberapi.MemberAPI,
	teamAPI teamapi.TeamAPI,
	conversationAPI conversationapi.ConversationAPIInterface,
	sourceControlAPI sourcecontrolapi.SourceControlAPI,
	organizationAPI orgapi.OrganizationAPI,
	userAPI userapi.UserAPI,
) *AIService {
	config := &types.AIServiceConfig{
		Provider:          getEnvOrDefault("AI_PROVIDER", "openai"),
		Model:             getEnvOrDefault("AI_MODEL", "gpt-4"),
		MaxTokens:         2000,
		Temperature:       0.7,
		MaxContextSize:    4000,
		EnableHistory:     true,
		EnableSuggestions: true,
	}

	return &AIService{
		db:               db,
		memberAPI:        memberAPI,
		teamAPI:          teamAPI,
		conversationAPI:  conversationAPI,
		sourceControlAPI: sourceControlAPI,
		organizationAPI:  organizationAPI,
		userAPI:          userAPI,
		config:           config,
	}
}

// Query processes an AI query request
func (s *AIService) Query(ctx context.Context, req *types.AIQueryRequest, userID string) (*types.AIQueryResponse, error) {
	// Gather relevant data based on entity type
	entityData, err := s.gatherEntityData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to gather entity data: %w", err)
	}

	// Format data for AI
	contextData := s.formatDataForAI(entityData, req.EntityType)

	// Build prompt
	prompt := s.buildPrompt(req.Query, contextData, req.EntityType, req.Context)

	// Call AI service (placeholder implementation)
	aiResponse, err := s.callAIService(prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI service: %w", err)
	}

	// Generate suggestions
	suggestions := s.generateSuggestions(req.EntityType, req.Context)

	response := &types.AIQueryResponse{
		Answer:      aiResponse.Answer,
		Sources:     aiResponse.Sources,
		Confidence:  aiResponse.Confidence,
		Suggestions: suggestions,
	}

	// Save to history if enabled
	if s.config.EnableHistory {
		history := &types.AIQueryHistory{
			OrganizationID: req.OrganizationID,
			UserID:         userID,
			EntityType:     req.EntityType,
			EntityID:       req.EntityID,
			Query:          req.Query,
			Answer:         response.Answer,
			Context:        req.Context,
			Confidence:     response.Confidence,
			Sources:        response.Sources,
			CreatedAt:      time.Now(),
		}

		if err := s.db.SaveQueryHistory(ctx, history); err != nil {
			log.Printf("Failed to save query history: %v", err)
		}
	}

	return response, nil
}

// GetQueryHistory retrieves query history
func (s *AIService) GetQueryHistory(ctx context.Context, organizationID string, userID *string, limit int) ([]*types.AIQueryHistory, error) {
	return s.db.GetQueryHistory(ctx, organizationID, userID, limit)
}

// GetQueryStats returns query statistics
func (s *AIService) GetQueryStats(ctx context.Context, organizationID string, userID *string, days int) (*types.AIQueryStats, error) {
	return s.db.GetQueryStats(ctx, organizationID, userID, days)
}

// gatherEntityData collects relevant data for the entity
func (s *AIService) gatherEntityData(ctx context.Context, req *types.AIQueryRequest) (*types.EntityData, error) {
	entityData := &types.EntityData{
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		OrganizationID: req.OrganizationID,
		Data:           make(map[string]interface{}),
		LastUpdated:    time.Now(),
	}

	switch req.EntityType {
	case "member":
		return s.gatherMemberData(ctx, req, entityData)
	case "team":
		return s.gatherTeamData(ctx, req, entityData)
	case "organization":
		return s.gatherOrganizationData(ctx, req, entityData)
	default:
		return nil, fmt.Errorf("unsupported entity type: %s", req.EntityType)
	}
}

// gatherMemberData collects member-specific data
func (s *AIService) gatherMemberData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) (*types.EntityData, error) {
	// Get member details
	members, err := s.memberAPI.GetOrganizationMembers(ctx, req.OrganizationID, nil)
	if err != nil {
		return nil, err
	}

	var member *membertypes.OrganizationMember
	for _, m := range members {
		if m.ID == req.EntityID {
			member = &m
			break
		}
	}

	if member == nil {
		return nil, fmt.Errorf("member not found")
	}

	entityData.Data["member"] = member

	// Get conversations if context includes conversations
	if strings.Contains(req.Context, "conversation") || req.Context == "overview" {
		conversations, err := s.conversationAPI.ListConversations(ctx, req.OrganizationID, &conversationtypes.ListConversationsQuery{
			DirectMemberID: &req.EntityID,
		})
		if err == nil {
			entityData.Data["conversations"] = conversations
		}
	}

	// Get source control data if context includes performance or overview
	if strings.Contains(req.Context, "performance") || req.Context == "overview" {
		// This would integrate with your existing source control metrics
		// For now, we'll add a placeholder
		entityData.Data["source_control_metrics"] = "Source control data would be integrated here"
	}

	return entityData, nil
}

// gatherTeamData collects team-specific data
func (s *AIService) gatherTeamData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) (*types.EntityData, error) {
	// Get team details
	team, err := s.teamAPI.GetTeam(ctx, req.EntityID)
	if err != nil {
		return nil, err
	}

	if team == nil {
		return nil, fmt.Errorf("team not found")
	}

	entityData.Data["team"] = team

	// Get team members
	members, err := s.memberAPI.GetOrganizationMembers(ctx, req.OrganizationID, nil)
	if err == nil {
		// Filter members for this team (you'd need to implement team membership logic)
		entityData.Data["team_members"] = members
	}

	return entityData, nil
}

// gatherOrganizationData collects organization-specific data
func (s *AIService) gatherOrganizationData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) (*types.EntityData, error) {
	// Get organization details
	org, err := s.organizationAPI.GetOrganizationByID(ctx, req.OrganizationID)
	if err != nil {
		return nil, err
	}

	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	entityData.Data["organization"] = org

	// Get organization members
	members, err := s.memberAPI.GetOrganizationMembers(ctx, req.OrganizationID, nil)
	if err == nil {
		entityData.Data["members"] = members
	}

	// Get teams
	teams, err := s.teamAPI.ListTeams(ctx, teamtypes.TeamSearchParams{})
	if err == nil {
		entityData.Data["teams"] = teams
	}

	return entityData, nil
}

// formatDataForAI formats the collected data for AI consumption
func (s *AIService) formatDataForAI(entityData *types.EntityData, entityType string) string {
	// Convert data to JSON for AI processing
	jsonData, err := json.MarshalIndent(entityData.Data, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal entity data: %v", err)
		return "Data formatting error"
	}

	// Truncate if too large
	dataStr := string(jsonData)
	if len(dataStr) > s.config.MaxContextSize {
		dataStr = dataStr[:s.config.MaxContextSize] + "... (truncated)"
	}

	return dataStr
}

// buildPrompt constructs the prompt for the AI service
func (s *AIService) buildPrompt(query, contextData, entityType, context string) string {
	var systemPrompt string

	switch entityType {
	case "member":
		systemPrompt = `You are an AI assistant helping managers understand their team members. You have access to member data including conversations, performance metrics, and other relevant information. Provide helpful, actionable insights based on the data.`
	case "team":
		systemPrompt = `You are an AI assistant helping managers understand their teams. You have access to team data including member information, performance metrics, and other relevant information. Provide helpful, actionable insights based on the data.`
	case "organization":
		systemPrompt = `You are an AI assistant helping with organizational insights. You have access to organization data including teams, members, and performance metrics. Provide helpful, actionable insights based on the data.`
	default:
		systemPrompt = `You are an AI assistant providing insights based on organizational data. Provide helpful, actionable insights based on the available data.`
	}

	prompt := fmt.Sprintf(`%s

Context: %s
Entity Type: %s

Available Data:
%s

User Query: %s

Please provide a helpful response based on the available data. If you need more information, please specify what additional data would be helpful.`,
		systemPrompt, context, entityType, contextData, query)

	return prompt
}

// callAIService calls the external AI service (placeholder implementation)
func (s *AIService) callAIService(prompt string) (*types.AIQueryResponse, error) {
	// This is a placeholder implementation
	// In a real implementation, you would call OpenAI, Anthropic, or another AI service

	// For now, return a mock response
	response := &types.AIQueryResponse{
		Answer:     "This is a placeholder AI response. The actual AI service integration would be implemented here.",
		Sources:    []string{"member_data", "conversations", "source_control"},
		Confidence: 0.85,
	}

	return response, nil
}

// generateSuggestions generates follow-up question suggestions
func (s *AIService) generateSuggestions(entityType, context string) []string {
	suggestions := []string{}

	switch entityType {
	case "member":
		switch context {
		case "performance":
			suggestions = []string{
				"How is this member's performance trending?",
				"What are this member's strengths?",
				"What areas could this member improve?",
			}
		case "conversations":
			suggestions = []string{
				"Summarize recent conversations with this member",
				"What topics have been discussed?",
				"Are there any action items from conversations?",
			}
		default:
			suggestions = []string{
				"Give me an overview of this member",
				"How is this member performing?",
				"What conversations have we had recently?",
			}
		}
	case "team":
		suggestions = []string{
			"How is the team performing overall?",
			"What are the team's main challenges?",
			"Who are the top performers in the team?",
		}
	case "organization":
		suggestions = []string{
			"How is the organization performing?",
			"What are the main organizational challenges?",
			"Which teams are performing best?",
		}
	}

	return suggestions
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
