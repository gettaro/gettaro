-- Rename organization_members table back to user_organizations
ALTER TABLE "organization_members" RENAME TO "user_organizations";

-- Rename indexes back
ALTER INDEX "organization_members_username_idx" RENAME TO "user_organizations_username_idx";
ALTER INDEX "organization_members_email_idx" RENAME TO "user_organizations_email_idx"; 