-- Add back the user_id column
ALTER TABLE "members_titles" ADD COLUMN "user_id" UUID;

-- Create foreign key constraint for user_id referencing users table
ALTER TABLE "members_titles" ADD CONSTRAINT "members_titles_user_id_fkey" 
FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;

-- Migrate existing data: Update user_id based on member_id
UPDATE "members_titles" 
SET "user_id" = (
    SELECT om.user_id 
    FROM "organization_members" om 
    WHERE om.id = "members_titles".member_id
    LIMIT 1
)
WHERE "members_titles".member_id IS NOT NULL;

-- Drop the member_id foreign key constraint and column
ALTER TABLE "members_titles" DROP CONSTRAINT IF EXISTS "members_titles_member_id_fkey";
ALTER TABLE "members_titles" DROP COLUMN "member_id";

-- Drop the member_id index
DROP INDEX IF EXISTS "members_titles_member_id_idx";

-- Update the unique constraint to use user_id instead of member_id
ALTER TABLE "members_titles" DROP CONSTRAINT IF EXISTS "members_titles_member_id_organization_id_key";
ALTER TABLE "members_titles" ADD CONSTRAINT "user_titles_user_id_organization_id_key" 
UNIQUE ("user_id", "organization_id");

-- Rename members_titles table back to user_titles
ALTER TABLE "members_titles" RENAME TO "user_titles"; 