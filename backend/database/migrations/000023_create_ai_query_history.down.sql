-- Drop AI query history table
DROP TABLE IF EXISTS ai_query_history CASCADE;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_ai_query_history_updated_at();
