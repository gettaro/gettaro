package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"ems.dev/backend/services/ai/database"
	"ems.dev/backend/services/ai/providers"
	"ems.dev/backend/services/ai/types"
	conversationapi "ems.dev/backend/services/conversation/api"
	conversationtypes "ems.dev/backend/services/conversation/types"
	memberapi "ems.dev/backend/services/member/api"
	membertypes "ems.dev/backend/services/member/types"
	orgapi "ems.dev/backend/services/organization/api"
	sourcecontrolapi "ems.dev/backend/services/sourcecontrol/api"
	sourcecontroltypes "ems.dev/backend/services/sourcecontrol/types"
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
	provider         providers.AIProviderInterface
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
	config *types.AIServiceConfig,
	provider providers.AIProviderInterface,
) *AIService {
	return &AIService{
		db:               db,
		memberAPI:        memberAPI,
		teamAPI:          teamAPI,
		conversationAPI:  conversationAPI,
		sourceControlAPI: sourceControlAPI,
		organizationAPI:  organizationAPI,
		userAPI:          userAPI,
		config:           config,
		provider:         provider,
	}
}

// Query processes an AI query request
func (s *AIService) Query(ctx context.Context, req *types.AIQueryRequest, userID string) (*types.AIQueryResponse, error) {
	// Generate retrieval plan using LLM
	retrievalPlan, err := s.generateRetrievalPlan(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to generate retrieval plan: %w", err)
	}

	// Execute the retrieval plan
	entityData, err := s.executeRetrievalPlan(ctx, req, retrievalPlan)
	if err != nil {
		return nil, fmt.Errorf("failed to execute retrieval plan: %w", err)
	}

	// Format data for AI
	contextData := s.formatDataForAI(entityData, req.EntityType)

	// Build prompt
	prompt := s.buildPrompt(req.Query, contextData, req.EntityType, req.Context)

	// Call AI provider
	aiResponse, err := s.provider.Query(ctx, prompt, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI provider: %w", err)
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
			Sources:        types.StringArray(response.Sources),
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

// generateRetrievalPlan uses LLM to generate a structured retrieval plan
func (s *AIService) generateRetrievalPlan(ctx context.Context, req *types.AIQueryRequest) (*types.RetrievalPlan, error) {
	// Build prompt for retrieval plan generation
	prompt := s.buildRetrievalPlanPrompt(req)

	// Call AI provider to generate plan
	aiResponse, err := s.provider.Query(ctx, prompt, s.config)
	if err != nil {
		return nil, fmt.Errorf("failed to call AI provider for retrieval plan: %w", err)
	}

	// Parse the JSON response
	var plan types.RetrievalPlan
	if err := json.Unmarshal([]byte(aiResponse.Answer), &plan); err != nil {
		// Fallback to default plan if parsing fails
		log.Printf("Failed to parse retrieval plan, using default: %v", err)
		return s.getDefaultRetrievalPlan(req), nil
	}

	// Validate the plan
	if err := s.validateRetrievalPlan(&plan, req); err != nil {
		log.Printf("Retrieval plan validation failed, using default: %v", err)
		return s.getDefaultRetrievalPlan(req), nil
	}

	return &plan, nil
}

// executeRetrievalPlan executes the generated retrieval plan
func (s *AIService) executeRetrievalPlan(ctx context.Context, req *types.AIQueryRequest, plan *types.RetrievalPlan) (*types.EntityData, error) {
	entityData := &types.EntityData{
		EntityType:     req.EntityType,
		EntityID:       req.EntityID,
		OrganizationID: req.OrganizationID,
		Data:           make(map[string]interface{}),
		LastUpdated:    time.Now(),
	}

	// Always include basic entity data
	if err := s.gatherBasicEntityData(ctx, req, entityData); err != nil {
		return nil, fmt.Errorf("failed to gather basic entity data: %w", err)
	}

	// Execute plan based on data sources
	for _, source := range plan.DataSources {
		switch source {
		case "source_control":
			if err := s.gatherSourceControlData(ctx, req, entityData, plan); err != nil {
				log.Printf("Failed to gather source control data: %v", err)
			}
		case "conversations":
			if err := s.gatherConversationData(ctx, req, entityData, plan); err != nil {
				log.Printf("Failed to gather conversation data: %v", err)
			}
		case "member_data":
			if err := s.gatherMemberData(ctx, req, entityData, plan); err != nil {
				log.Printf("Failed to gather member data: %v", err)
			}
		case "team_data":
			if err := s.gatherTeamData(ctx, req, entityData, plan); err != nil {
				log.Printf("Failed to gather team data: %v", err)
			}
		case "organization_data":
			if err := s.gatherOrganizationData(ctx, req, entityData, plan); err != nil {
				log.Printf("Failed to gather organization data: %v", err)
			}
		}
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

// buildRetrievalPlanPrompt builds a prompt for generating retrieval plans
func (s *AIService) buildRetrievalPlanPrompt(req *types.AIQueryRequest) string {
	return fmt.Sprintf(`You are a data retrieval planning assistant. Based on the user's query, generate a JSON plan for what data to retrieve.

Available data sources:
- source_control: Pull requests, reviews, metrics, code contributions
- conversations: Manager conversations, feedback, action items
- member_data: Basic member information
- team_data: Team composition, performance, structure
- organization_data: Organization-wide metrics, teams, members

Query: "%s"
Entity Type: %s
Context: %s

Generate a JSON response with this exact structure:
{
  "data_sources": ["source_control", "conversations"],
  "time": {
    "from": "2025-09-01",
    "to": "2025-09-27", 
    "interval": "daily"
  },
  "filters": {
    "member_ids": ["member_id_1"],
    "limit": 100
  },
  "priority": "high",
  "reasoning": "Brief explanation of why this plan was chosen"
}

Rules:
1. Only include data sources that are relevant to the query
2. When the user is requesting an overview use all the data sources directly related to the entity. Ex. If member use source_control and conversations. If team use team_data and member_data. If organization use organization_data and team_data and member_data and source_control.
3. Set appropriate time ranges based on query context
4. If no start and end date are provided, use the last 30 days
5. Use "daily" interval for recent data, "weekly" for trends, "monthly" for overviews
6. Set priority based on query urgency
7. Provide clear reasoning for your choices
8. Return ONLY valid JSON, no other text`, req.Query, req.EntityType, req.Context)
}

// validateRetrievalPlan validates a generated retrieval plan
func (s *AIService) validateRetrievalPlan(plan *types.RetrievalPlan, req *types.AIQueryRequest) error {
	// Check if data sources are valid
	validSources := map[string]bool{
		"source_control":    true,
		"conversations":     true,
		"member_data":       true,
		"team_data":         true,
		"organization_data": true,
	}

	for _, source := range plan.DataSources {
		if !validSources[source] {
			return fmt.Errorf("invalid data source: %s", source)
		}
	}

	// Validate time range if provided
	if plan.Time != nil {
		if plan.Time.From == "" || plan.Time.To == "" {
			return fmt.Errorf("time range must have both 'from' and 'to' dates")
		}

		validIntervals := map[string]bool{
			"daily":   true,
			"weekly":  true,
			"monthly": true,
			"yearly":  true,
		}
		if !validIntervals[plan.Time.Interval] {
			return fmt.Errorf("invalid time interval: %s", plan.Time.Interval)
		}
	}

	// Validate priority
	if plan.Priority != "" {
		validPriorities := map[string]bool{
			"high":   true,
			"medium": true,
			"low":    true,
		}
		if !validPriorities[plan.Priority] {
			return fmt.Errorf("invalid priority: %s", plan.Priority)
		}
	}

	return nil
}

// getDefaultRetrievalPlan returns a default retrieval plan based on entity type
func (s *AIService) getDefaultRetrievalPlan(req *types.AIQueryRequest) *types.RetrievalPlan {
	plan := &types.RetrievalPlan{
		DataSources: []string{"member_data"},
		Priority:    "medium",
		Reasoning:   "Default plan based on entity type",
	}

	// Add time range for recent data
	now := time.Now()
	plan.Time = &types.TimeRange{
		From:     now.AddDate(0, 0, -30).Format("2006-01-02"),
		To:       now.Format("2006-01-02"),
		Interval: "day",
	}

	// Add relevant data sources based on entity type
	switch req.EntityType {
	case "member":
		plan.DataSources = []string{"member_data", "source_control", "conversations"}
	case "team":
		plan.DataSources = []string{"team_data", "member_data", "source_control"}
	case "organization":
		plan.DataSources = []string{"organization_data", "team_data", "member_data"}
	}

	return plan
}

// gatherBasicEntityData gathers basic entity information
func (s *AIService) gatherBasicEntityData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) error {
	switch req.EntityType {
	case "member":
		return s.gatherMemberBasicData(ctx, req, entityData)
	case "team":
		return s.gatherTeamBasicData(ctx, req, entityData)
	case "organization":
		return s.gatherOrganizationBasicData(ctx, req, entityData)
	default:
		return fmt.Errorf("unsupported entity type: %s", req.EntityType)
	}
}

// gatherMemberBasicData gathers basic member information
func (s *AIService) gatherMemberBasicData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) error {
	members, err := s.memberAPI.GetOrganizationMembers(ctx, req.OrganizationID, nil)
	if err != nil {
		return err
	}

	var member *membertypes.OrganizationMember
	for _, m := range members {
		if m.ID == req.EntityID {
			member = &m
			break
		}
	}

	if member == nil {
		return fmt.Errorf("member not found")
	}

	entityData.Data["member"] = member
	return nil
}

// gatherTeamBasicData gathers basic team information
func (s *AIService) gatherTeamBasicData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) error {
	team, err := s.teamAPI.GetTeam(ctx, req.EntityID)
	if err != nil {
		return err
	}

	if team == nil {
		return fmt.Errorf("team not found")
	}

	entityData.Data["team"] = team
	return nil
}

// gatherOrganizationBasicData gathers basic organization information
func (s *AIService) gatherOrganizationBasicData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData) error {
	org, err := s.organizationAPI.GetOrganizationByID(ctx, req.OrganizationID)
	if err != nil {
		return err
	}

	if org == nil {
		return fmt.Errorf("organization not found")
	}

	entityData.Data["organization"] = org
	return nil
}

// gatherSourceControlData gathers source control data based on plan
func (s *AIService) gatherSourceControlData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData, plan *types.RetrievalPlan) error {
	// Get member's source control accounts
	sourceControlAccounts, err := s.sourceControlAPI.GetSourceControlAccounts(ctx, &sourcecontroltypes.SourceControlAccountParams{
		OrganizationID: req.OrganizationID,
		MemberIDs:      []string{req.EntityID},
	})
	if err != nil || len(sourceControlAccounts) == 0 {
		return err
	}

	entityData.Data["source_control_accounts"] = sourceControlAccounts

	// Determine time range
	startDate := time.Now().AddDate(0, 0, -90) // Default to 90 days
	endDate := time.Now()
	if plan.Time != nil {
		if parsedStartTime, err := time.Parse("2006-01-02", plan.Time.From); err == nil {
			startDate = parsedStartTime
		}
		if parsedEndTime, err := time.Parse("2006-01-02", plan.Time.To); err == nil {
			endDate = parsedEndTime
		}
	}

	// Get pull requests
	pullRequests, err := s.sourceControlAPI.GetMemberPullRequests(ctx, &sourcecontroltypes.MemberPullRequestParams{
		MemberID:        req.EntityID,
		StartDate:       &startDate,
		EndDate:         &endDate,
		IncludeComments: &[]bool{true}[0],
	})
	if err == nil {
		entityData.Data["pull_requests"] = pullRequests
	}

	// Get pull request reviews
	reviews, err := s.sourceControlAPI.GetMemberPullRequestReviews(ctx, &sourcecontroltypes.MemberPullRequestReviewsParams{
		MemberID:  req.EntityID,
		StartDate: &startDate,
		EndDate:   &endDate,
		HasBody:   &[]bool{true}[0],
	})
	if err == nil {
		entityData.Data["member_pull_request_reviews"] = reviews
	}

	// Calculate metrics for the target member
	metricsResponse, err := s.memberAPI.CalculateSourceControlMemberMetrics(ctx, req.OrganizationID, req.EntityID, sourcecontroltypes.MemberMetricsParams{
		StartDate: &startDate,
		EndDate:   &endDate,
		Interval:  plan.Time.Interval,
	})

	if err == nil {
		entityData.Data["metrics"] = metricsResponse
	}

	return nil
}

// gatherConversationData gathers conversation data based on plan
func (s *AIService) gatherConversationData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData, plan *types.RetrievalPlan) error {
	conversations, err := s.conversationAPI.ListConversations(ctx, req.OrganizationID, &conversationtypes.ListConversationsQuery{
		DirectMemberID: &req.EntityID,
	})
	if err != nil {
		return err
	}

	entityData.Data["conversations"] = conversations
	return nil
}

// gatherMemberData gathers member data based on plan (updated signature)
func (s *AIService) gatherMemberData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData, plan *types.RetrievalPlan) error {
	// This method can be expanded to include additional member-specific data gathering
	// based on the retrieval plan filters and requirements
	return nil
}

// gatherTeamData gathers team data based on plan (updated signature)
func (s *AIService) gatherTeamData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData, plan *types.RetrievalPlan) error {
	// Get team members
	members, err := s.memberAPI.GetOrganizationMembers(ctx, req.OrganizationID, nil)
	if err == nil {
		entityData.Data["team_members"] = members
	}
	return err
}

// gatherOrganizationData gathers organization data based on plan (updated signature)
func (s *AIService) gatherOrganizationData(ctx context.Context, req *types.AIQueryRequest, entityData *types.EntityData, plan *types.RetrievalPlan) error {
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

	return nil
}

// formatAccountIDsForJSON formats source control account IDs as JSON array string
func (s *AIService) formatAccountIDsForJSON(accountIDs []string) string {
	if len(accountIDs) == 0 {
		return "[]"
	}

	jsonStr := "["
	for i, id := range accountIDs {
		if i > 0 {
			jsonStr += ","
		}
		jsonStr += `"` + id + `"`
	}
	jsonStr += "]"

	return jsonStr
}
