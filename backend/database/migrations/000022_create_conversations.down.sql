-- Drop conversations table and related objects
DROP TRIGGER IF EXISTS conversations_updated_at ON conversations;
DROP FUNCTION IF EXISTS update_conversations_updated_at();
DROP TABLE IF EXISTS conversations;
