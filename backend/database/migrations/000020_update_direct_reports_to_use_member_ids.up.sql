
-- Update direct_reports table to use member IDs instead of user IDs
-- First, drop the existing foreign key constraints
ALTER TABLE direct_reports DROP CONSTRAINT IF EXISTS direct_reports_manager_id_fkey;
ALTER TABLE direct_reports DROP CONSTRAINT IF EXISTS direct_reports_report_id_fkey;

-- Rename the columns to be more explicit
ALTER TABLE direct_reports RENAME COLUMN manager_id TO manager_member_id;
ALTER TABLE direct_reports RENAME COLUMN report_id TO report_member_id;

-- Add new foreign key constraints referencing organization_members
ALTER TABLE direct_reports ADD CONSTRAINT direct_reports_manager_member_id_fkey 
    FOREIGN KEY (manager_member_id) REFERENCES organization_members(id) ON DELETE CASCADE;
ALTER TABLE direct_reports ADD CONSTRAINT direct_reports_report_member_id_fkey 
    FOREIGN KEY (report_member_id) REFERENCES organization_members(id) ON DELETE CASCADE;

-- Update the unique constraint to use the new column names
ALTER TABLE direct_reports DROP CONSTRAINT IF EXISTS direct_reports_manager_id_report_id_organization_id_key;
ALTER TABLE direct_reports ADD CONSTRAINT direct_reports_manager_member_id_report_member_id_organization_id_key 
    UNIQUE (manager_member_id, report_member_id, organization_id);
