-- Drop the author_id column and its foreign key constraint
ALTER TABLE "pr_comments" DROP CONSTRAINT IF EXISTS "pr_comments_author_id_fkey";
ALTER TABLE "pr_comments" DROP COLUMN IF EXISTS "author_id";

-- Add source_control_account_id column
ALTER TABLE "pr_comments" ADD COLUMN "source_control_account_id" UUID;

-- Add foreign key constraint for source_control_account_id
ALTER TABLE "pr_comments" ADD CONSTRAINT "pr_comments_source_control_account_id_fkey"
FOREIGN KEY ("source_control_account_id") REFERENCES "source_control_accounts"("id") ON DELETE SET NULL;

-- Create an index for better query performance
CREATE INDEX "pr_comments_source_control_account_id_idx" ON "pr_comments"("source_control_account_id"); 