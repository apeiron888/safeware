package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Company represents an organization using the system
type Company struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string                 `bson:"name" json:"name"`
	Domain    string                 `bson:"domain,omitempty" json:"domain,omitempty"`
	Settings  map[string]interface{} `bson:"settings,omitempty" json:"settings,omitempty"`
	CreatedAt time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time              `bson:"updated_at" json:"updated_at"`
}

// User represents a system user
type User struct {
	ID                       primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID                primitive.ObjectID     `bson:"company_id" json:"company_id"`
	FullName                 string                 `bson:"full_name" json:"full_name"`
	Email                    string                 `bson:"email" json:"email"`
	Phone                    string                 `bson:"phone,omitempty" json:"phone,omitempty"`
	PasswordHash             string                 `bson:"password_hash" json:"-"`
	TOTPSecret               string                 `bson:"totp_secret,omitempty" json:"-"`
	ClearanceLevel           int                    `bson:"clearance_level" json:"clearance_level"`
	Attributes               map[string]interface{} `bson:"attributes,omitempty" json:"attributes,omitempty"`
	IsVerified               bool                   `bson:"is_verified" json:"is_verified"`
	VerificationToken        string                 `bson:"verification_token,omitempty" json:"-"`
	VerificationTokenExpires *time.Time             `bson:"verification_token_expires,omitempty" json:"-"`
	ResetPasswordToken       string                 `bson:"reset_password_token,omitempty" json:"-"`
	ResetPasswordExpires     *time.Time             `bson:"reset_password_expires,omitempty" json:"-"`
	FailedLogins             int                    `bson:"failed_logins" json:"-"`
	LockedUntil              *time.Time             `bson:"locked_until,omitempty" json:"-"`
	LastLogin                *time.Time             `bson:"last_login,omitempty" json:"last_login,omitempty"`
	CreatedAt                time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt                time.Time              `bson:"updated_at" json:"updated_at"`
	Role                     string                 `bson:"role" json:"role"` // Manager, Supervisor, Staff, Auditor
	WarehouseID              *primitive.ObjectID    `bson:"warehouse_id,omitempty" json:"warehouse_id,omitempty"`
	Roles                    []primitive.ObjectID   `bson:"role_ids,omitempty" json:"-"` // Deprecated: keeping for backward compat if needed, but logic will use Role field
}

// Role represents a user role
type Role struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID      primitive.ObjectID   `bson:"company_id" json:"company_id"`
	Name           string               `bson:"name" json:"name"`
	Description    string               `bson:"description,omitempty" json:"description,omitempty"`
	HierarchyLevel int                  `bson:"hierarchy_level" json:"hierarchy_level"`
	CreatedAt      time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updated_at"`
	PermissionIDs  []primitive.ObjectID `bson:"permission_ids,omitempty" json:"permission_ids,omitempty"`
}

// Permission represents an action permission
type Permission struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Description  string             `bson:"description,omitempty" json:"description,omitempty"`
	ResourceType string             `bson:"resource_type,omitempty" json:"resource_type,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
