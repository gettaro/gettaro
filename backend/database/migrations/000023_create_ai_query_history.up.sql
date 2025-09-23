-- Create AI query history table
CREATE TABLE IF NOT EXISTS ai_query_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    user_id UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    query TEXT NOT NULL,
    answer TEXT NOT NULL,
    context VARCHAR(100),
    confidence DECIMAL(3,2) DEFAULT 0.0,
    sources TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_ai_query_history_organization_id ON ai_query_history(organization_id);
CREATE INDEX idx_ai_query_history_user_id ON ai_query_history(user_id);
CREATE INDEX idx_ai_query_history_entity_type ON ai_query_history(entity_type);
CREATE INDEX idx_ai_query_history_entity_id ON ai_query_history(entity_id);
CREATE INDEX idx_ai_query_history_created_at ON ai_query_history(created_at);
CREATE INDEX idx_ai_query_history_org_user ON ai_query_history(organization_id, user_id);
CREATE INDEX idx_ai_query_history_org_entity ON ai_query_history(organization_id, entity_type, entity_id);

-- Add foreign key constraints
ALTER TABLE ai_query_history 
ADD CONSTRAINT fk_ai_query_history_organization 
FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE;

ALTER TABLE ai_query_history 
ADD CONSTRAINT fk_ai_query_history_user 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add updated_at trigger
CREATE OR REPLACE FUNCTION update_ai_query_history_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Note: We don't add updated_at column since query history is immutable
-- But we keep the trigger function for potential future use
