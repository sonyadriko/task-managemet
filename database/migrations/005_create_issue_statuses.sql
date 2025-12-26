-- Migration: Create issue_statuses table
-- Description: Flexible workflow per organization (shared by all teams)

CREATE TABLE issue_statuses (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    position INTEGER NOT NULL,
    is_final BOOLEAN DEFAULT FALSE,
    color VARCHAR(7) DEFAULT '#6B7280',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(organization_id, name),
    UNIQUE(organization_id, position)
);

CREATE INDEX idx_issue_statuses_org ON issue_statuses(organization_id);
