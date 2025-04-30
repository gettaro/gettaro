-- Drop the index
DROP INDEX IF EXISTS "pull_requests_repository_name_idx";

-- Drop the repository_name column
ALTER TABLE "pull_requests" DROP COLUMN IF EXISTS "repository_name"; 