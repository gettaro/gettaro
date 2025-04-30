-- Add organization_id column to pull_requests table
ALTER TABLE "pull_requests" ADD COLUMN "organization_id" UUID NOT NULL;

-- Add foreign key constraint
ALTER TABLE "pull_requests" ADD CONSTRAINT "pull_requests_organization_id_fkey" 
FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE CASCADE;

-- Create an index for better query performance
CREATE INDEX "pull_requests_organization_id_idx" ON "pull_requests"("organization_id"); 