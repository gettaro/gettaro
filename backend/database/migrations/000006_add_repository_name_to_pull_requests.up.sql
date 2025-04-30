-- Add repository_name column to pull_requests table
ALTER TABLE "pull_requests" ADD COLUMN "repository_name" TEXT NOT NULL;

-- Create an index for better query performance
CREATE INDEX "pull_requests_repository_name_idx" ON "pull_requests"("repository_name"); 