-- Make the updated_at column nullable
ALTER TABLE "pr_comments" ALTER COLUMN "updated_at" DROP NOT NULL; 

-- Update records with default timestamp (0001-01-01 00:00:00+00) to null
UPDATE "pr_comments" 
SET "updated_at" = NULL 
WHERE "updated_at" = '0001-01-01 00:00:00+00';
