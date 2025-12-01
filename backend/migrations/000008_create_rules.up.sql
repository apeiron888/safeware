CREATE TYPE rule_effect AS ENUM ('allow', 'deny');

CREATE TABLE IF NOT EXISTS rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    condition_expression JSONB NOT NULL,
    effect rule_effect NOT NULL,
    priority INTEGER DEFAULT 0,
    enabled BOOLEAN DEFAULT TRUE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, name)
);

CREATE INDEX idx_rules_company ON rules(company_id);
CREATE INDEX idx_rules_enabled ON rules(enabled);
CREATE INDEX idx_rules_priority ON rules(priority);
CREATE INDEX idx_rules_effect ON rules(effect);

COMMENT ON COLUMN rules.condition_expression IS 'JSON DSL for RuBAC/ABAC conditions (time, location, attributes, etc.)';
