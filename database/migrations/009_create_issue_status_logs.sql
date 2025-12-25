-- Migration: Create issue_status_logs table
-- Description: Audit trail for status changes

CREATE TABLE issue_status_logs (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER REFERENCES issues(id) ON DELETE CASCADE,
    from_status_id INTEGER REFERENCES issue_statuses(id) ON DELETE SET NULL,
    to_status_id INTEGER REFERENCES issue_statuses(id) ON DELETE SET NULL,
    changed_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_status_logs_issue ON issue_status_logs(issue_id);
CREATE INDEX idx_status_logs_changed_at ON issue_status_logs(changed_at);
