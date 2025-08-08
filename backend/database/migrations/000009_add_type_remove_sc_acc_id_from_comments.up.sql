-- Add type column to pr_comments table
ALTER TABLE "pr_comments" ADD COLUMN "type" VARCHAR(255) NOT NULL DEFAULT 'COMMENT';

-- Drop the source_control_account_id column and its foreign key constraint
ALTER TABLE "pr_comments" DROP CONSTRAINT IF EXISTS "pr_comments_source_control_account_id_fkey";
ALTER TABLE "pr_comments" DROP COLUMN IF EXISTS "source_control_account_id";

-- Create an index on the type column for better query performance
CREATE INDEX "pr_comments_type_idx" ON "pr_comments"("type");
