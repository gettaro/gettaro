-- Add back the user_id column
ALTER TABLE "source_control_accounts" ADD COLUMN "user_id" UUID;

-- Create foreign key constraint for user_id referencing users table
ALTER TABLE "source_control_accounts" ADD CONSTRAINT "source_control_accounts_user_id_fkey" 
FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE SET NULL;

-- Migrate existing data: Update user_id based on member_id
UPDATE "source_control_accounts" 
SET "user_id" = (
    SELECT om.user_id 
    FROM "organization_members" om 
    WHERE om.id = "source_control_accounts".member_id
    LIMIT 1
)
WHERE "source_control_accounts".member_id IS NOT NULL;

-- Drop the member_id foreign key constraint and column
ALTER TABLE "source_control_accounts" DROP CONSTRAINT IF EXISTS "source_control_accounts_member_id_fkey";
ALTER TABLE "source_control_accounts" DROP COLUMN "member_id";

-- Drop the member_id index
DROP INDEX IF EXISTS "source_control_accounts_member_id_idx"; 