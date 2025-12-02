export interface Item {
    id: string;
    sku?: string;
    name: string;
    quality?: string;
    quantity: number;
    price: number;
    department?: string;
    warehouse_id: string;
    batch?: string;
    created_at?: string;
    updated_at?: string;
    // Legacy fields for backwards compatibility
    category?: string;
    description?: string;
}

export interface Warehouse {
    id: string;
    name: string;
    location: string;
    capacity?: number;
    supervisor_id?: string;
    staff_ids?: string[];
    items_count?: number;
    total_value?: number;
    created_at?: string;
    updated_at?: string;
}

export interface Employee {
    id: string;
    full_name: string;
    email: string;
    role: 'Manager' | 'Supervisor' | 'Staff' | 'Auditor';
    warehouse_id?: string;
    company_id: string;
    created_at?: string;
}

export interface AuditLog {
    id: string;
    action: string;
    user_id: string;
    user_name?: string;
    resource_type: string;
    resource_id: string;
    timestamp: string;
    details?: string;
    ip_address?: string;
}

export interface WarehouseWithDetails extends Warehouse {
    supervisor_name?: string;
    staff_count?: number;
    items?: Item[];
}
