-- Add username and email columns to user_organizations table
ALTER TABLE "user_organizations" ADD COLUMN "username" VARCHAR(255);
ALTER TABLE "user_organizations" ADD COLUMN "email" VARCHAR(255);

-- Create indexes for better query performance
CREATE INDEX "user_organizations_username_idx" ON "user_organizations"("username");
CREATE INDEX "user_organizations_email_idx" ON "user_organizations"("email"); 