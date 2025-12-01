package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Warehouse represents a storage location
type Warehouse struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID    primitive.ObjectID     `bson:"company_id" json:"company_id"`
	Name         string                 `bson:"name" json:"name"`
	Location     string                 `bson:"location" json:"location"`
	SupervisorID *primitive.ObjectID    `bson:"supervisor_id,omitempty" json:"supervisor_id,omitempty"`
	IPWhitelist  []string               `bson:"ip_whitelist,omitempty" json:"ip_whitelist,omitempty"`
	Attributes   map[string]interface{} `bson:"attributes,omitempty" json:"attributes,omitempty"`
	IsActive     bool                   `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time              `bson:"updated_at" json:"updated_at"`
}

// Item represents an inventory item
type Item struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID      primitive.ObjectID     `bson:"company_id" json:"company_id"`
	SKU            string                 `bson:"sku" json:"sku"`
	Name           string                 `bson:"name" json:"name"`
	Quality        string                 `bson:"quality" json:"quality"` // New, Used, Damaged
	Price          float64                `bson:"price,omitempty" json:"price,omitempty"`
	Classification string                 `bson:"classification,omitempty" json:"classification,omitempty"` // Deprecated
	OwnerUserID    primitive.ObjectID     `bson:"owner_user_id,omitempty" json:"owner_user_id,omitempty"`
	Department     string                 `bson:"department,omitempty" json:"department,omitempty"`
	Attributes     map[string]interface{} `bson:"attributes,omitempty" json:"attributes,omitempty"`
	IsArchived     bool                   `bson:"is_archived" json:"is_archived"`
	CreatedAt      time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time              `bson:"updated_at" json:"updated_at"`
}

// ItemLocation tracks where items are stored
type ItemLocation struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	ItemID      primitive.ObjectID `bson:"item_id" json:"item_id"`
	WarehouseID primitive.ObjectID `bson:"warehouse_id" json:"warehouse_id"`
	Quantity    int                `bson:"quantity" json:"quantity"`
	Batch       string             `bson:"batch,omitempty" json:"batch,omitempty"`
	UpdatedBy   primitive.ObjectID `bson:"updated_by,omitempty" json:"updated_by,omitempty"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Transfer represents an item transfer request
type Transfer struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID       primitive.ObjectID `bson:"company_id" json:"company_id"`
	ItemID          primitive.ObjectID `bson:"item_id" json:"item_id"`
	FromWarehouseID primitive.ObjectID `bson:"from_warehouse_id" json:"from_warehouse_id"`
	ToWarehouseID   primitive.ObjectID `bson:"to_warehouse_id" json:"to_warehouse_id"`
	Quantity        int                `bson:"quantity" json:"quantity"`
	Status          string             `bson:"status" json:"status"` // requested, approved, rejected, completed, cancelled
	Reason          string             `bson:"reason,omitempty" json:"reason,omitempty"`
	RequestedBy     primitive.ObjectID `bson:"requested_by" json:"requested_by"`
	ApprovedBy      primitive.ObjectID `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	RejectionReason string             `bson:"rejection_reason,omitempty" json:"rejection_reason,omitempty"`
	CreatedAt       time.Time          `bson:"created_at" json:"created_at"`
	ApprovedAt      *time.Time         `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
	CompletedAt     *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	UpdatedAt       time.Time          `bson:"updated_at" json:"updated_at"`
}
