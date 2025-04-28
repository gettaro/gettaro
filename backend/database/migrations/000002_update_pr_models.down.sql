-- Remove new columns from pull_requests table
ALTER TABLE pull_requests
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS comments,
    DROP COLUMN IF EXISTS review_comments,
    DROP COLUMN IF EXISTS additions,
    DROP COLUMN IF EXISTS deletions,
    DROP COLUMN IF EXISTS changed_files,
    DROP COLUMN IF EXISTS metadata;

-- Revert pr_comments table changes
ALTER TABLE pr_comments
    DROP CONSTRAINT IF EXISTS pr_comments_source_control_account_id_fkey,
    DROP COLUMN IF EXISTS source_control_account_id,
    ADD COLUMN IF NOT EXISTS author_id UUID;

-- Recreate pr_reviewers table
CREATE TABLE pr_reviewers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pr_id UUID NOT NULL,
    reviewer_id UUID NOT NULL,
    reviewed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    review_state VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pr_id) REFERENCES pull_requests(id) ON DELETE CASCADE,
    FOREIGN KEY (reviewer_id) REFERENCES source_control_accounts(id) ON DELETE CASCADE
); 