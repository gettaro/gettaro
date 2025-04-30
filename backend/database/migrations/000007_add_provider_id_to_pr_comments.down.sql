-- Drop the index
DROP INDEX IF EXISTS "pr_comments_provider_id_idx";

-- Drop the provider_id column
ALTER TABLE "pr_comments" DROP COLUMN IF EXISTS "provider_id"; 