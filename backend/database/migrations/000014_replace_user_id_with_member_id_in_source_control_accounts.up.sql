-- Drop the existing foreign key constraint for user_id
ALTER TABLE "source_control_accounts" DROP CONSTRAINT IF EXISTS "source_control_accounts_user_id_fkey";

-- Add the new member_id column
ALTER TABLE "source_control_accounts" ADD COLUMN "member_id" UUID;

-- Create foreign key constraint for member_id referencing organization_members table
ALTER TABLE "source_control_accounts" ADD CONSTRAINT "source_control_accounts_member_id_fkey" 
FOREIGN KEY ("member_id") REFERENCES "organization_members"("id") ON DELETE SET NULL;

-- Create an index for better query performance
CREATE INDEX "source_control_accounts_member_id_idx" ON "source_control_accounts"("member_id");

-- Migrate existing data: Update member_id based on user_id and organization_id
UPDATE "source_control_accounts" 
SET "member_id" = (
    SELECT om.id 
    FROM "organization_members" om 
    WHERE om.user_id = "source_control_accounts".user_id 
    AND om.organization_id = "source_control_accounts".organization_id
    LIMIT 1
)
WHERE "source_control_accounts".user_id IS NOT NULL 
AND "source_control_accounts".organization_id IS NOT NULL;

-- Drop the old user_id column
ALTER TABLE "source_control_accounts" DROP COLUMN "user_id"; 