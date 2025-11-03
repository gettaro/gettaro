-- Migration: Revert ai_code_assistant_daily_metrics table fields
-- This migration reverts the changes: renames lines_of_code_suggested back to total_suggestions and adds suggestions_accepted

-- Rename lines_of_code_suggested column back to total_suggestions
ALTER TABLE ai_code_assistant_daily_metrics
    RENAME COLUMN lines_of_code_suggested TO total_suggestions;

-- Add suggestions_accepted column back
ALTER TABLE ai_code_assistant_daily_metrics
    ADD COLUMN suggestions_accepted INTEGER DEFAULT 0;

