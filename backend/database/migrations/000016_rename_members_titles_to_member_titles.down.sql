-- Rename member_titles table back to members_titles
ALTER TABLE "member_titles" RENAME TO "members_titles";

-- Rename the index back to match the original table name
ALTER INDEX "member_titles_member_id_idx" RENAME TO "members_titles_member_id_idx";

-- Rename the foreign key constraint back to match the original table name
ALTER TABLE "members_titles" DROP CONSTRAINT IF EXISTS "member_titles_member_id_fkey";
ALTER TABLE "members_titles" ADD CONSTRAINT "members_titles_member_id_fkey" 
FOREIGN KEY ("member_id") REFERENCES "organization_members"("id") ON DELETE CASCADE;

-- Rename the unique constraint back to match the original table name
ALTER TABLE "members_titles" DROP CONSTRAINT IF EXISTS "member_titles_member_id_organization_id_key";
ALTER TABLE "members_titles" ADD CONSTRAINT "members_titles_member_id_organization_id_key" 
UNIQUE ("member_id", "organization_id"); 