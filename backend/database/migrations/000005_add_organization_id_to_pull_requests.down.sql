-- Drop the index
DROP INDEX IF EXISTS "pull_requests_organization_id_idx";

-- Drop the foreign key constraint
ALTER TABLE "pull_requests" DROP CONSTRAINT IF EXISTS "pull_requests_organization_id_fkey";

-- Drop the organization_id column
ALTER TABLE "pull_requests" DROP COLUMN IF EXISTS "organization_id"; 