CREATE TYPE grant_target_type AS ENUM ('user', 'role');

CREATE TABLE IF NOT EXISTS dac_grants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    target_type grant_target_type NOT NULL,
    target_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    target_role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    permission_id UUID NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    granted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    CHECK (
        (target_type = 'user' AND target_user_id IS NOT NULL AND target_role_id IS NULL) OR
        (target_type = 'role' AND target_role_id IS NOT NULL AND target_user_id IS NULL)
    )
);

CREATE INDEX idx_dac_grants_owner ON dac_grants(owner_user_id);
CREATE INDEX idx_dac_grants_target_user ON dac_grants(target_user_id);
CREATE INDEX idx_dac_grants_target_role ON dac_grants(target_role_id);
CREATE INDEX idx_dac_grants_resource ON dac_grants(resource_type, resource_id);
CREATE INDEX idx_dac_grants_permission ON dac_grants(permission_id);
CREATE INDEX idx_dac_grants_active ON dac_grants(is_active);
CREATE INDEX idx_dac_grants_expires ON dac_grants(expires_at);
