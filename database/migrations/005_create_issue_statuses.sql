-- Migration: Create issue_statuses table
-- Description: Flexible workflow per team (not hardcoded ENUM)

CREATE TABLE issue_statuses (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    position INTEGER NOT NULL,
    is_final BOOLEAN DEFAULT FALSE,
    color VARCHAR(7) DEFAULT '#6B7280',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, name),
    UNIQUE(team_id, position)
);

CREATE INDEX idx_issue_statuses_team ON issue_statuses(team_id);
