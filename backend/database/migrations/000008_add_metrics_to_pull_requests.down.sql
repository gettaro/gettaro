-- Drop metrics column
ALTER TABLE "pull_requests" DROP COLUMN IF EXISTS "metrics";

-- Restore residual columns
ALTER TABLE "pull_requests" ADD COLUMN "time_to_first_review" INTEGER;
ALTER TABLE "pull_requests" ADD COLUMN "time_to_merge" INTEGER; 