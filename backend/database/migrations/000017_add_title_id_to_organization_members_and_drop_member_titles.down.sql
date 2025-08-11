-- Recreate the member_titles table
CREATE TABLE "member_titles" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "member_id" UUID NOT NULL,
    "title_id" UUID NOT NULL,
    "organization_id" UUID NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("member_id") REFERENCES "organization_members"("id") ON DELETE CASCADE,
    FOREIGN KEY ("title_id") REFERENCES "titles"("id") ON DELETE CASCADE,
    FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE CASCADE,
    UNIQUE ("member_id", "organization_id")
);

-- Create index for member_id
CREATE INDEX "member_titles_member_id_idx" ON "member_titles"("member_id");

-- Migrate data back: Insert into member_titles based on organization_members.title_id
INSERT INTO "member_titles" ("member_id", "title_id", "organization_id")
SELECT "id", "title_id", "organization_id"
FROM "organization_members"
WHERE "title_id" IS NOT NULL;

-- Drop the title_id column and its constraints from organization_members
ALTER TABLE "organization_members" DROP CONSTRAINT IF EXISTS "organization_members_title_id_fkey";
ALTER TABLE "organization_members" DROP INDEX IF EXISTS "organization_members_title_id_idx";
ALTER TABLE "organization_members" DROP COLUMN IF EXISTS "title_id"; 