-- Rename members_titles table to member_titles
ALTER TABLE "members_titles" RENAME TO "member_titles";

-- Rename the index to match the new table name
ALTER INDEX "members_titles_member_id_idx" RENAME TO "member_titles_member_id_idx";

-- Rename the foreign key constraint to match the new table name
ALTER TABLE "member_titles" DROP CONSTRAINT IF EXISTS "members_titles_member_id_fkey";
ALTER TABLE "member_titles" ADD CONSTRAINT "member_titles_member_id_fkey" 
FOREIGN KEY ("member_id") REFERENCES "organization_members"("id") ON DELETE CASCADE;

-- Rename the unique constraint to match the new table name
ALTER TABLE "member_titles" DROP CONSTRAINT IF EXISTS "members_titles_member_id_organization_id_key";
ALTER TABLE "member_titles" ADD CONSTRAINT "member_titles_member_id_organization_id_key" 
UNIQUE ("member_id", "organization_id"); 