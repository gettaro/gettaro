-- Add is_manager column to titles table
ALTER TABLE titles ADD COLUMN is_manager BOOLEAN DEFAULT false;
