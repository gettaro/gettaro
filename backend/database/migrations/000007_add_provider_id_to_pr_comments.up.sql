-- Add provider_id column to pr_comments table
ALTER TABLE "pr_comments" ADD COLUMN "provider_id" VARCHAR(255) NOT NULL;

-- Create an index for better query performance
CREATE INDEX "pr_comments_provider_id_idx" ON "pr_comments"("provider_id"); 