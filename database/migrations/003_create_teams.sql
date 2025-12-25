-- Migration: Create teams table
-- Description: Teams/departments with hierarchical support

CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    organization_id INTEGER REFERENCES organizations(id) ON DELETE CASCADE,
    parent_team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TRIGGER update_teams_updated_at BEFORE UPDATE ON teams
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE INDEX idx_teams_org ON teams(organization_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_teams_parent ON teams(parent_team_id) WHERE deleted_at IS NULL;
