-- Add author_id to pr_comments table
ALTER TABLE pr_comments
    ADD COLUMN author_id UUID,
    ADD CONSTRAINT pr_comments_author_id_fkey 
        FOREIGN KEY (author_id) 
        REFERENCES source_control_accounts(id) 
        ON DELETE SET NULL; 