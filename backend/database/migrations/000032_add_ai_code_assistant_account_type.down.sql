-- Rollback: Revert account_type constraint to only allow sourcecontrol
ALTER TABLE member_external_accounts 
DROP CONSTRAINT IF EXISTS member_external_accounts_account_type_check;

ALTER TABLE member_external_accounts 
ADD CONSTRAINT member_external_accounts_account_type_check 
CHECK (account_type IN ('sourcecontrol'));

