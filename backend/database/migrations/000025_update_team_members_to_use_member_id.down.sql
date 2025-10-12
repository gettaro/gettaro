-- Rollback: Update team_members table to use user_id instead of member_id

-- Add the old user_id column back
ALTER TABLE "team_members" ADD COLUMN "user_id" UUID;

-- Create foreign key constraint for user_id referencing users table
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_user_id_fkey" 
FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;

-- Create index for user_id
CREATE INDEX "team_members_user_id_idx" ON "team_members"("user_id");

-- Migrate existing data: Update user_id based on member_id
UPDATE "team_members" 
SET "user_id" = (
    SELECT om.user_id 
    FROM "organization_members" om 
    WHERE om.id = "team_members".member_id
)
WHERE "team_members".member_id IS NOT NULL;

-- Update the unique constraint to use user_id instead of member_id
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_member_id_team_id_organization_id_key";
ALTER TABLE "team_members" ADD CONSTRAINT "team_members_user_id_team_id_organization_id_key" 
UNIQUE ("user_id", "team_id", "organization_id");

-- Drop the member_id foreign key constraint and column
ALTER TABLE "team_members" DROP CONSTRAINT IF EXISTS "team_members_member_id_fkey";
ALTER TABLE "team_members" DROP COLUMN "member_id";

-- Drop the member_id index
DROP INDEX IF EXISTS "team_members_member_id_idx";
