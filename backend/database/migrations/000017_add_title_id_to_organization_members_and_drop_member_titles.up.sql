-- Add title_id column to organization_members table
ALTER TABLE "organization_members" ADD COLUMN "title_id" UUID;

-- Create foreign key constraint for title_id referencing titles table
ALTER TABLE "organization_members" ADD CONSTRAINT "organization_members_title_id_fkey" 
FOREIGN KEY ("title_id") REFERENCES "titles"("id") ON DELETE SET NULL;

-- Create an index for better query performance
CREATE INDEX "organization_members_title_id_idx" ON "organization_members"("title_id");

-- Migrate existing data: Update title_id based on member_titles table
UPDATE "organization_members" 
SET "title_id" = (
    SELECT mt.title_id 
    FROM "member_titles" mt 
    WHERE mt.member_id = "organization_members".id
    LIMIT 1
)
WHERE EXISTS (
    SELECT 1 FROM "member_titles" mt 
    WHERE mt.member_id = "organization_members".id
);

-- Drop the member_titles table and its constraints
DROP TABLE IF EXISTS "member_titles" CASCADE; 