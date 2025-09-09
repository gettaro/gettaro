CREATE TABLE conversation_templates (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  template_fields JSONB NOT NULL DEFAULT '[]'::jsonb,
  is_active BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_conversation_templates_organization_id ON conversation_templates(organization_id);
CREATE INDEX idx_conversation_templates_is_active ON conversation_templates(is_active);

-- Add updated_at trigger
CREATE OR REPLACE FUNCTION update_conversation_templates_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_conversation_templates_updated_at
  BEFORE UPDATE ON conversation_templates
  FOR EACH ROW
  EXECUTE FUNCTION update_conversation_templates_updated_at();
