-- Drop the indexes
DROP INDEX IF EXISTS "user_organizations_username_idx";
DROP INDEX IF EXISTS "user_organizations_email_idx";

-- Drop the username and email columns
ALTER TABLE "user_organizations" DROP COLUMN IF EXISTS "username";
ALTER TABLE "user_organizations" DROP COLUMN IF EXISTS "email"; 