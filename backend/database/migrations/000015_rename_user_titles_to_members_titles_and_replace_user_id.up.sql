-- Rename user_titles table to members_titles
ALTER TABLE "user_titles" RENAME TO "members_titles";

-- Add the new member_id column
ALTER TABLE "members_titles" ADD COLUMN "member_id" UUID;

-- Create foreign key constraint for member_id referencing organization_members table
ALTER TABLE "members_titles" ADD CONSTRAINT "members_titles_member_id_fkey" 
FOREIGN KEY ("member_id") REFERENCES "organization_members"("id") ON DELETE CASCADE;

-- Create an index for better query performance
CREATE INDEX "members_titles_member_id_idx" ON "members_titles"("member_id");

-- Migrate existing data: Update member_id based on user_id and organization_id
UPDATE "members_titles" 
SET "member_id" = (
    SELECT om.id 
    FROM "organization_members" om 
    WHERE om.user_id = "members_titles".user_id 
    AND om.organization_id = "members_titles".organization_id
    LIMIT 1
)
WHERE "members_titles".user_id IS NOT NULL 
AND "members_titles".organization_id IS NOT NULL;

-- Drop the old user_id column and its foreign key constraint
ALTER TABLE "members_titles" DROP CONSTRAINT IF EXISTS "members_titles_user_id_fkey";
ALTER TABLE "members_titles" DROP COLUMN "user_id";

-- Update the unique constraint to use member_id instead of user_id
ALTER TABLE "members_titles" DROP CONSTRAINT IF EXISTS "user_titles_user_id_organization_id_key";
ALTER TABLE "members_titles" ADD CONSTRAINT "members_titles_member_id_organization_id_key" 
UNIQUE ("member_id", "organization_id"); 