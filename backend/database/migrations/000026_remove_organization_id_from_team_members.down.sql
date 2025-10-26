-- Rollback: Add organization_id column back to team_members table

-- Add the organization_id column back
ALTER TABLE "team_members" ADD COLUMN "organization_id" UUID;

-- Create foreign key constraint for organization_id referencing organizations table
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_organization_id_fkey" 
FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE CASCADE;

-- Migrate existing data: Update organization_id based on team_id
UPDATE "team_members" 
SET "organization_id" = (
    SELECT t.organization_id 
    FROM "teams" t 
    WHERE t.id = "team_members".team_id
)
WHERE "team_members".team_id IS NOT NULL;

-- Drop the unique constraint without organization_id
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_member_id_team_id_key";

-- Add back the unique constraint with organization_id
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_member_id_team_id_organization_id_key" 
UNIQUE ("member_id", "team_id", "organization_id");
