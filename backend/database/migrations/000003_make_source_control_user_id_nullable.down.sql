-- Revert user_id nullable in source_control_accounts table
ALTER TABLE source_control_accounts 
    DROP CONSTRAINT IF EXISTS source_control_accounts_user_id_fkey,
    ADD CONSTRAINT source_control_accounts_user_id_fkey 
        FOREIGN KEY (user_id) 
        REFERENCES users(id) 
        ON DELETE CASCADE,
    ALTER COLUMN user_id SET NOT NULL; 