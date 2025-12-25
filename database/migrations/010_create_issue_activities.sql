-- Migration: Create issue_activities table
-- Description: Comprehensive activity log for all issue events

CREATE TYPE activity_type AS ENUM (
    'created', 
    'assigned', 
    'status_changed', 
    'priority_changed', 
    'commented', 
    'hold', 
    'resumed'
);

CREATE TABLE issue_activities (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER REFERENCES issues(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    activity_type activity_type NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_activities_issue ON issue_activities(issue_id);
CREATE INDEX idx_activities_type ON issue_activities(activity_type);
CREATE INDEX idx_activities_created_at ON issue_activities(created_at);
