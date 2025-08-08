-- Drop the index on the type column
DROP INDEX IF EXISTS "pr_comments_type_idx";

-- Add back the source_control_account_id column
ALTER TABLE "pr_comments" ADD COLUMN "source_control_account_id" UUID;

-- Add back the foreign key constraint
ALTER TABLE "pr_comments" ADD CONSTRAINT "pr_comments_source_control_account_id_fkey"
FOREIGN KEY ("source_control_account_id") REFERENCES "source_control_accounts"("id") ON DELETE SET NULL;

-- Drop the type column
ALTER TABLE "pr_comments" DROP COLUMN IF EXISTS "type";
