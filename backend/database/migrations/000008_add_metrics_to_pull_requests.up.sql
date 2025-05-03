-- Add metrics column to pull_requests table
ALTER TABLE "pull_requests" ADD COLUMN "metrics" JSONB;

-- Drop residual columns
ALTER TABLE "pull_requests" DROP COLUMN IF EXISTS "time_to_first_review";
ALTER TABLE "pull_requests" DROP COLUMN IF EXISTS "time_to_merge"; 