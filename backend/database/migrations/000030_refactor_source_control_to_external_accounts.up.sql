-- Migration: Refactor source_control_accounts to member_external_accounts
-- This migration:
-- 1. Creates the new member_external_accounts table
-- 2. Migrates data from source_control_accounts
-- 3. Updates foreign key references in pull_requests and pr_comments
-- 4. Drops the old table

-- Step 1: Create new member_external_accounts table
CREATE TABLE member_external_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID,
    organization_id UUID,
    account_type VARCHAR(50) NOT NULL CHECK (account_type IN ('sourcecontrol')),
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

-- Step 2: Create indexes
CREATE INDEX member_external_accounts_member_id_idx ON member_external_accounts(member_id);
CREATE INDEX member_external_accounts_organization_id_idx ON member_external_accounts(organization_id);
CREATE INDEX member_external_accounts_account_type_idx ON member_external_accounts(account_type);
CREATE INDEX member_external_accounts_member_account_type_idx ON member_external_accounts(member_id, account_type);

-- Step 3: Migrate data from source_control_accounts
INSERT INTO member_external_accounts (
    id, member_id, organization_id, account_type, provider_name, provider_id, 
    username, metadata, last_synced_at, created_at, updated_at
)
SELECT 
    id, member_id, organization_id, 'sourcecontrol', provider_name, provider_id,
    username, metadata, last_synced_at, created_at, updated_at
FROM source_control_accounts;

-- Step 4: Update pull_requests table
-- Rename the column from source_control_account_id to external_account_id
ALTER TABLE pull_requests 
    RENAME COLUMN source_control_account_id TO external_account_id;

-- Drop old foreign key constraint
ALTER TABLE pull_requests
    DROP CONSTRAINT IF EXISTS pull_requests_source_control_account_id_fkey;

-- Add new foreign key constraint
ALTER TABLE pull_requests
    ADD CONSTRAINT pull_requests_external_account_id_fkey
    FOREIGN KEY (external_account_id) REFERENCES member_external_accounts(id) ON DELETE CASCADE;

-- Step 5: Update pr_comments table (if it has the column)
-- Check first if column exists, then update
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'pr_comments' AND column_name = 'source_control_account_id'
    ) THEN
        ALTER TABLE pr_comments 
            RENAME COLUMN source_control_account_id TO external_account_id;
        
        ALTER TABLE pr_comments
            DROP CONSTRAINT IF EXISTS pr_comments_source_control_account_id_fkey;
        
        ALTER TABLE pr_comments
            ADD CONSTRAINT pr_comments_external_account_id_fkey
            FOREIGN KEY (external_account_id) REFERENCES member_external_accounts(id) ON DELETE SET NULL;
    END IF;
END $$;

-- Step 6: Drop old table and its constraints
ALTER TABLE source_control_accounts 
    DROP CONSTRAINT IF EXISTS source_control_accounts_member_id_fkey;
ALTER TABLE source_control_accounts 
    DROP CONSTRAINT IF EXISTS source_control_accounts_organization_id_fkey;
DROP TABLE IF EXISTS source_control_accounts;

