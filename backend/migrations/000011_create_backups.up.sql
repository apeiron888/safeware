CREATE TABLE IF NOT EXISTS backups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_size BIGINT,
    backup_type VARCHAR(50) DEFAULT 'full',
    encrypted BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    retention_until TIMESTAMP NOT NULL,
    restored_at TIMESTAMP,
    notes TEXT
);

CREATE INDEX idx_backups_company ON backups(company_id);
CREATE INDEX idx_backups_created ON backups(created_at DESC);
CREATE INDEX idx_backups_retention ON backups(retention_until);
