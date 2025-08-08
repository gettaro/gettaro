-- Drop the index on source_control_account_id
DROP INDEX IF EXISTS "pr_comments_source_control_account_id_idx";

-- Drop the foreign key constraint for source_control_account_id
ALTER TABLE "pr_comments" DROP CONSTRAINT IF EXISTS "pr_comments_source_control_account_id_fkey";

-- Drop the source_control_account_id column
ALTER TABLE "pr_comments" DROP COLUMN IF EXISTS "source_control_account_id";

-- Add back the author_id column
ALTER TABLE "pr_comments" ADD COLUMN "author_id" UUID;

-- Add back the foreign key constraint for author_id
ALTER TABLE "pr_comments" ADD CONSTRAINT "pr_comments_author_id_fkey"
FOREIGN KEY ("author_id") REFERENCES "users"("id") ON DELETE CASCADE; 