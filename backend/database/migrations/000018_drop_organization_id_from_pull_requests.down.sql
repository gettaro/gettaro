-- Add back the organization_id column to pull_requests table
ALTER TABLE "pull_requests" ADD COLUMN "organization_id" UUID;
ALTER TABLE "pull_requests" ADD CONSTRAINT "pull_requests_organization_id_fkey" 
FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE CASCADE; 