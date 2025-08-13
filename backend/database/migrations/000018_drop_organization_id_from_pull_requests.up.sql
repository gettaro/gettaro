-- Drop the organization_id column from pull_requests table
ALTER TABLE "pull_requests" DROP CONSTRAINT IF EXISTS "pull_requests_organization_id_fkey";
ALTER TABLE "pull_requests" DROP COLUMN IF EXISTS "organization_id"; 