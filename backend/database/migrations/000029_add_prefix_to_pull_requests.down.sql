DROP INDEX IF EXISTS idx_pull_requests_prefix;
ALTER TABLE pull_requests DROP COLUMN prefix;

