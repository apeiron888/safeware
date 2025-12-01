CREATE TYPE transfer_status AS ENUM ('requested', 'approved', 'rejected', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    from_warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    to_warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status transfer_status DEFAULT 'requested',
    reason TEXT,
    requested_by UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
    rejection_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    approved_at TIMESTAMP,
    completed_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (from_warehouse_id != to_warehouse_id)
);

CREATE INDEX idx_transfers_company ON transfers(company_id);
CREATE INDEX idx_transfers_item ON transfers(item_id);
CREATE INDEX idx_transfers_status ON transfers(status);
CREATE INDEX idx_transfers_requester ON transfers(requested_by);
CREATE INDEX idx_transfers_approver ON transfers(approved_by);
CREATE INDEX idx_transfers_created ON transfers(created_at);
