-- Update team_members table to use member_id instead of user_id
-- This aligns with the pattern used in other tables like conversations and direct_reports

-- Add the new member_id column
ALTER TABLE "team_members" ADD COLUMN "member_id" UUID;

-- Create foreign key constraint for member_id referencing organization_members table
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_member_id_fkey" 
FOREIGN KEY ("member_id") REFERENCES "organization_members"("id") ON DELETE CASCADE;

-- Create index for member_id
CREATE INDEX "team_members_member_id_idx" ON "team_members"("member_id");

-- Migrate existing data: Update member_id based on user_id and organization_id
UPDATE "team_members" 
SET "member_id" = (
    SELECT om.id 
    FROM "organization_members" om 
    WHERE om.user_id = "team_members".user_id 
    AND om.organization_id = "team_members".organization_id
)
WHERE "team_members".user_id IS NOT NULL;

-- Update the unique constraint to use member_id instead of user_id
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_user_id_team_id_organization_id_key";
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_member_id_team_id_organization_id_key" 
UNIQUE ("member_id", "team_id", "organization_id");

-- Drop the old user_id foreign key constraint and column
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_user_id_fkey";
ALTER TABLE "team_members" DROP COLUMN "user_id";

-- Drop the old user_id index
DROP INDEX IF EXISTS "team_members_user_id_idx";
