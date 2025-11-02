-- Migration: Add ai-code-assistant account type to member_external_accounts
-- This migration updates the check constraint to allow 'ai-code-assistant' as a valid account_type

-- Drop the existing constraint
ALTER TABLE member_external_accounts 
DROP CONSTRAINT IF EXISTS member_external_accounts_account_type_check;

-- Add the new constraint with ai-code-assistant
ALTER TABLE member_external_accounts 
ADD CONSTRAINT member_external_accounts_account_type_check 
CHECK (account_type IN ('sourcecontrol', 'ai-code-assistant'));

