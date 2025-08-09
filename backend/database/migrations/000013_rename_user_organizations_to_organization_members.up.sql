-- Rename user_organizations table to organization_members
ALTER TABLE "user_organizations" RENAME TO "organization_members";

-- Rename indexes
ALTER INDEX "user_organizations_username_idx" RENAME TO "organization_members_username_idx";
ALTER INDEX "user_organizations_email_idx" RENAME TO "organization_members_email_idx"; 