ALTER TABLE pull_requests ADD COLUMN prefix VARCHAR(50);
CREATE INDEX idx_pull_requests_prefix ON pull_requests(prefix);

