-- Make the updated_at column not nullable again
ALTER TABLE "pr_comments" ALTER COLUMN "updated_at" SET NOT NULL;

-- Update null values back to the default timestamp
UPDATE "pr_comments" 
SET "updated_at" = '0001-01-01 00:00:00+00' 
WHERE "updated_at" IS NULL; 