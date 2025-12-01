CREATE TYPE classification_level AS ENUM ('Public', 'Internal', 'Restricted', 'HighValue');

CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    sku VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    classification classification_level DEFAULT 'Public',
    owner_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    department VARCHAR(100),
    attributes JSONB DEFAULT '{}',
    is_archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(company_id, sku)
);

CREATE TABLE IF NOT EXISTS item_locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    batch VARCHAR(100),
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_items_company ON items(company_id);
CREATE INDEX idx_items_sku ON items(sku);
CREATE INDEX idx_items_classification ON items(classification);
CREATE INDEX idx_items_owner ON items(owner_user_id);
CREATE INDEX idx_items_archived ON items(is_archived);
CREATE INDEX idx_item_locations_item ON item_locations(item_id);
CREATE INDEX idx_item_locations_warehouse ON item_locations(warehouse_id);
