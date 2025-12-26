-- Migration: Create issue_attachments table
-- Description: Store metadata for files attached to issues (files stored in Cloudflare R2)

CREATE TABLE issue_attachments (
    id SERIAL PRIMARY KEY,
    issue_id INTEGER NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
    filename VARCHAR(255) NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    storage_key VARCHAR(500) NOT NULL,
    uploaded_by INTEGER REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_attachments_issue ON issue_attachments(issue_id);
CREATE INDEX idx_attachments_uploaded_by ON issue_attachments(uploaded_by);
