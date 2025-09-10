-- Create conversations table
CREATE TABLE conversations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    template_id UUID REFERENCES conversation_templates(id) ON DELETE SET NULL,
    manager_member_id UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    direct_member_id UUID NOT NULL REFERENCES organization_members(id) ON DELETE CASCADE,
    conversation_date DATE,
    status VARCHAR(50) NOT NULL DEFAULT 'draft', -- draft, completed
    content JSONB, -- Filled template data
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for efficient queries
CREATE INDEX idx_conversations_organization_id ON conversations(organization_id);
CREATE INDEX idx_conversations_manager_member_id ON conversations(manager_member_id);
CREATE INDEX idx_conversations_direct_member_id ON conversations(direct_member_id);
CREATE INDEX idx_conversations_template_id ON conversations(template_id);
CREATE INDEX idx_conversations_status ON conversations(status);
CREATE INDEX idx_conversations_conversation_date ON conversations(conversation_date);

-- Create updated_at trigger
CREATE OR REPLACE FUNCTION update_conversations_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_conversations_updated_at
    BEFORE UPDATE ON conversations
    FOR EACH ROW
    EXECUTE FUNCTION update_conversations_updated_at();
