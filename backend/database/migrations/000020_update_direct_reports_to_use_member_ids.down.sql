-- Revert direct_reports table back to using user IDs
-- Drop the new foreign key constraints
ALTER TABLE direct_reports DROP CONSTRAINT IF EXISTS direct_reports_manager_member_id_fkey;
ALTER TABLE direct_reports DROP CONSTRAINT IF EXISTS direct_reports_report_member_id_fkey;

-- Rename the columns back to original names
ALTER TABLE direct_reports RENAME COLUMN manager_member_id TO manager_id;
ALTER TABLE direct_reports RENAME COLUMN report_member_id TO report_id;

-- Add back the original foreign key constraints referencing users
ALTER TABLE direct_reports ADD CONSTRAINT direct_reports_manager_id_fkey 
    FOREIGN KEY (manager_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE direct_reports ADD CONSTRAINT direct_reports_report_id_fkey 
    FOREIGN KEY (report_id) REFERENCES users(id) ON DELETE CASCADE;

-- Update the unique constraint back to original names
ALTER TABLE direct_reports DROP CONSTRAINT IF EXISTS direct_reports_manager_member_id_report_member_id_organization_id_key;
ALTER TABLE direct_reports ADD CONSTRAINT direct_reports_manager_id_report_id_organization_id_key 
    UNIQUE (manager_id, report_id, organization_id);
