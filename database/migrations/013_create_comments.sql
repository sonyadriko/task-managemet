-- Migration: Create issue_comments table
-- Description: Store comments on issues

CREATE TABLE issue_comments (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_issue ON issue_comments(issue_id);
CREATE INDEX idx_comments_user ON issue_comments(user_id);
