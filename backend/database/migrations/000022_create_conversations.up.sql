-- Create conversations table
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    template_id UUID REFERENCES conversation_templates(id) ON DELETE SET NULL,
    manager_id UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    direct_report_id UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content JSONB, -- Stores template questions and answers
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index for efficient queries
CREATE INDEX idx_conversations_organization_id ON conversations(organization_id);
CREATE INDEX idx_conversations_manager_id ON conversations(manager_id);
CREATE INDEX idx_conversations_direct_report_id ON conversations(direct_report_id);
CREATE INDEX idx_conversations_template_id ON conversations(template_id);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_conversations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_conversations_updated_at();
