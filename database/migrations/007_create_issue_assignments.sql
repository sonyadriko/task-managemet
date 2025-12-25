-- Migration: Create issue_assignments table
-- Description: WHO works on WHAT and WHEN (source for calendar view)

CREATE TABLE issue_assignments (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER REFERENCES issues(id) ON DELETE CASCADE,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    CONSTRAINT valid_date_range CHECK (end_date >= start_date)
);

CREATE INDEX idx_assignments_issue ON issue_assignments(issue_id);
CREATE INDEX idx_assignments_user ON issue_assignments(user_id) WHERE is_active = TRUE;
CREATE INDEX idx_assignments_dates ON issue_assignments(start_date, end_date);
