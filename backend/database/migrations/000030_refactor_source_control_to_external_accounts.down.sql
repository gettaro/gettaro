-- Rollback Migration: Revert from member_external_accounts back to source_control_accounts

-- Step 1: Recreate old source_control_accounts table
CREATE TABLE source_control_accounts (
    id UUID PRIMARY KEY,
    member_id UUID,
    organization_id UUID,
    provider_name VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL,
    metadata JSONB,
    last_synced_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES organization_members(id) ON DELETE SET NULL,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE SET NULL
);

-- Step 2: Recreate indexes on source_control_accounts
CREATE INDEX source_control_accounts_member_id_idx ON source_control_accounts(member_id);
CREATE INDEX source_control_accounts_organization_id_idx ON source_control_accounts(organization_id);

-- Step 3: Migrate data back (only sourcecontrol accounts)
INSERT INTO source_control_accounts (
    id, member_id, organization_id, provider_name, provider_id, 
    username, metadata, last_synced_at, created_at, updated_at
)
SELECT 
    id, member_id, organization_id, provider_name, provider_id,
    username, metadata, last_synced_at, created_at, updated_at
FROM member_external_accounts
WHERE account_type = 'sourcecontrol';

-- Step 4: Revert pull_requests table
ALTER TABLE pull_requests 
    RENAME COLUMN external_account_id TO source_control_account_id;

ALTER TABLE pull_requests
    DROP CONSTRAINT IF EXISTS pull_requests_external_account_id_fkey;

ALTER TABLE pull_requests
    ADD CONSTRAINT pull_requests_source_control_account_id_fkey
    FOREIGN KEY (source_control_account_id) REFERENCES source_control_accounts(id) ON DELETE CASCADE;

-- Step 5: Revert pr_comments table if needed
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'pr_comments' AND column_name = 'external_account_id'
    ) THEN
        ALTER TABLE pr_comments 
            RENAME COLUMN external_account_id TO source_control_account_id;
        
        ALTER TABLE pr_comments
            DROP CONSTRAINT IF EXISTS pr_comments_external_account_id_fkey;
        
        ALTER TABLE pr_comments
            ADD CONSTRAINT pr_comments_source_control_account_id_fkey
            FOREIGN KEY (source_control_account_id) REFERENCES source_control_accounts(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Step 6: Drop new table
DROP TABLE IF EXISTS member_external_accounts;

