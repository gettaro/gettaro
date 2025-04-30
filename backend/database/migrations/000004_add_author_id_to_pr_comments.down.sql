-- Remove author_id from pr_comments table
ALTER TABLE pr_comments
    DROP CONSTRAINT IF EXISTS pr_comments_author_id_fkey,
    DROP COLUMN author_id; 