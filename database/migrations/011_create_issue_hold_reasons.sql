-- Migration: Create issue_hold_reasons table
-- Description: Tracks why tasks are blocked/on hold

CREATE TABLE issue_hold_reasons (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER REFERENCES issues(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP,
    resolved_by INTEGER REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX idx_hold_reasons_issue ON issue_hold_reasons(issue_id);
CREATE INDEX idx_hold_reasons_unresolved ON issue_hold_reasons(issue_id) WHERE resolved_at IS NULL;
