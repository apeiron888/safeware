package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog represents system audit trail
type AuditLog struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID        primitive.ObjectID `bson:"company_id,omitempty" json:"company_id,omitempty"`
	UserID           primitive.ObjectID `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Username         string             `bson:"username,omitempty" json:"username,omitempty"`
	Action           string             `bson:"action" json:"action"`
	ResourceType     string             `bson:"resource_type,omitempty" json:"resource_type,omitempty"`
	ResourceID       primitive.ObjectID `bson:"resource_id,omitempty" json:"resource_id,omitempty"`
	Status           string             `bson:"status,omitempty" json:"status,omitempty"`
	DetailsEncrypted string             `bson:"details_encrypted,omitempty" json:"-"`
	IPAddress        string             `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	UserAgent        string             `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
}

// Session represents user authentication sessions
type Session struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	UserID           primitive.ObjectID     `bson:"user_id" json:"user_id"`
	RefreshTokenHash string                 `bson:"refresh_token_hash" json:"-"`
	DeviceInfo       map[string]interface{} `bson:"device_info,omitempty" json:"device_info,omitempty"`
	IPAddress        string                 `bson:"ip_address,omitempty" json:"ip_address,omitempty"`
	UserAgent        string                 `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	IsActive         bool                   `bson:"is_active" json:"is_active"`
	ExpiresAt        time.Time              `bson:"expires_at" json:"expires_at"`
	CreatedAt        time.Time              `bson:"created_at" json:"created_at"`
	LastUsedAt       time.Time              `bson:"last_used_at" json:"last_used_at"`
}

// Backup represents database backup metadata
type Backup struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID      primitive.ObjectID `bson:"company_id,omitempty" json:"company_id,omitempty"`
	FilePath       string             `bson:"file_path" json:"file_path"`
	FileSize       int64              `bson:"file_size,omitempty" json:"file_size,omitempty"`
	BackupType     string             `bson:"backup_type" json:"backup_type"`
	Encrypted      bool               `bson:"encrypted" json:"encrypted"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	RetentionUntil time.Time          `bson:"retention_until" json:"retention_until"`
	RestoredAt     *time.Time         `bson:"restored_at,omitempty" json:"restored_at,omitempty"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty"`
}
