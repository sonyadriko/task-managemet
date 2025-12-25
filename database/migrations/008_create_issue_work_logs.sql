-- Migration: Create issue_work_logs table
-- Description: Tracks actual work done (planning vs reality)

CREATE TABLE issue_work_logs (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER REFERENCES issues(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    work_date DATE NOT NULL,
    minutes_spent INTEGER NOT NULL CHECK (minutes_spent > 0),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_work_logs_issue ON issue_work_logs(issue_id);
CREATE INDEX idx_work_logs_user ON issue_work_logs(user_id);
CREATE INDEX idx_work_logs_date ON issue_work_logs(work_date);
