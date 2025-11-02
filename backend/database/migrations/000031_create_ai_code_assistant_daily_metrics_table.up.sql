-- Migration: Create ai_code_assistant_daily_metrics table
-- This migration creates a table to store daily aggregated metrics for AI code assistant tools
-- at the user level (overlapping metrics from both Cursor Analytics API and Claude Code Usage Analytics)

-- Create ai_code_assistant_daily_metrics table
CREATE TABLE ai_code_assistant_daily_metrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    external_account_id UUID NOT NULL,
    tool_name VARCHAR(255) NOT NULL,
    metric_date DATE NOT NULL,
    lines_of_code_accepted INTEGER DEFAULT 0,
    total_suggestions INTEGER DEFAULT 0,
    suggestions_accepted INTEGER DEFAULT 0,
    suggestion_accept_rate DECIMAL(5,2),
    active_sessions INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (external_account_id) REFERENCES member_external_accounts(id) ON DELETE CASCADE,
    CONSTRAINT ai_code_assistant_daily_metrics_unique UNIQUE (organization_id, external_account_id, tool_name, metric_date)
);

-- Create indexes for performance
CREATE INDEX idx_ai_code_assistant_daily_metrics_org_id ON ai_code_assistant_daily_metrics(organization_id);
CREATE INDEX idx_ai_code_assistant_daily_metrics_external_account_id ON ai_code_assistant_daily_metrics(external_account_id);
CREATE INDEX idx_ai_code_assistant_daily_metrics_date ON ai_code_assistant_daily_metrics(metric_date);
CREATE INDEX idx_ai_code_assistant_daily_metrics_tool_name ON ai_code_assistant_daily_metrics(tool_name);
CREATE INDEX idx_ai_code_assistant_daily_metrics_user_date ON ai_code_assistant_daily_metrics(external_account_id, metric_date);
CREATE INDEX idx_ai_code_assistant_daily_metrics_org_user_date ON ai_code_assistant_daily_metrics(organization_id, external_account_id, metric_date);

