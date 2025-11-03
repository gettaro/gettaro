-- Migration: Update ai_code_assistant_daily_metrics table fields
-- This migration renames total_suggestions to lines_of_code_suggested and removes suggestions_accepted

-- Rename total_suggestions column to lines_of_code_suggested
ALTER TABLE ai_code_assistant_daily_metrics
    RENAME COLUMN total_suggestions TO lines_of_code_suggested;

-- Drop suggestions_accepted column
ALTER TABLE ai_code_assistant_daily_metrics
    DROP COLUMN suggestions_accepted;

