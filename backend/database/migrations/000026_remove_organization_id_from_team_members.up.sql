-- Remove organization_id column from team_members table
-- This column is redundant since we can get the organization through the team relationship

-- Drop the unique constraint that includes organization_id
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_member_id_team_id_organization_id_key";

-- Create new unique constraint without organization_id
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_member_id_team_id_key" 
UNIQUE ("member_id", "team_id");

-- Drop the foreign key constraint for organization_id
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_organization_id_fkey";

-- Drop the organization_id column
ALTER TABLE "team_members" DROP COLUMN "organization_id";
