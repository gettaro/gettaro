-- Drop pr_reviewers table
DROP TABLE IF EXISTS pr_reviewers;

-- Update pr_comments table
ALTER TABLE pr_comments
    DROP COLUMN IF EXISTS author_id,
    ADD COLUMN IF NOT EXISTS source_control_account_id UUID,
    ADD CONSTRAINT pr_comments_source_control_account_id_fkey
    FOREIGN KEY (source_control_account_id)
    REFERENCES source_control_accounts(id)
    ON DELETE SET NULL;

-- Update pull_requests table
ALTER TABLE pull_requests
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS comments INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS review_comments INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS additions INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS deletions INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS changed_files INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS metadata JSONB; 