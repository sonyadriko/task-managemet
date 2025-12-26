-- Migration: Add deadline column to issues
-- Description: Add optional deadline date for issues

ALTER TABLE issues ADD COLUMN deadline DATE;

CREATE INDEX idx_issues_deadline ON issues(deadline);
