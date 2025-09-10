-- Add title column to conversations table
ALTER TABLE conversations ADD COLUMN title VARCHAR(255) NOT NULL DEFAULT 'Untitled Conversation';

-- Update existing conversations to use template name as title if available
UPDATE conversations 
SET title = COALESCE(
    (SELECT ct.name FROM conversation_templates ct WHERE ct.id = conversations.template_id::uuid),
    'Untitled Conversation'
)
WHERE title = 'Untitled Conversation';
