package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DACGrant represents discretionary access control grants
type DACGrant struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID    primitive.ObjectID `bson:"company_id" json:"company_id"`
	OwnerUserID  primitive.ObjectID `bson:"owner_user_id" json:"owner_user_id"`
	TargetType   string             `bson:"target_type" json:"target_type"` // user or role
	TargetUserID primitive.ObjectID `bson:"target_user_id,omitempty" json:"target_user_id,omitempty"`
	TargetRoleID primitive.ObjectID `bson:"target_role_id,omitempty" json:"target_role_id,omitempty"`
	ResourceType string             `bson:"resource_type" json:"resource_type"`
	ResourceID   primitive.ObjectID `bson:"resource_id" json:"resource_id"`
	PermissionID primitive.ObjectID `bson:"permission_id" json:"permission_id"`
	GrantedAt    time.Time          `bson:"granted_at" json:"granted_at"`
	ExpiresAt    *time.Time         `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	IsActive     bool               `bson:"is_active" json:"is_active"`
}

// Rule represents RuBAC/ABAC rules
type Rule struct {
	ID                  primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	CompanyID           primitive.ObjectID     `bson:"company_id" json:"company_id"`
	Name                string                 `bson:"name" json:"name"`
	Description         string                 `bson:"description,omitempty" json:"description,omitempty"`
	ConditionExpression map[string]interface{} `bson:"condition_expression" json:"condition_expression"`
	Effect              string                 `bson:"effect" json:"effect"` // allow or deny
	Priority            int                    `bson:"priority" json:"priority"`
	Enabled             bool                   `bson:"enabled" json:"enabled"`
	CreatedBy           primitive.ObjectID     `bson:"created_by,omitempty" json:"created_by,omitempty"`
	CreatedAt           time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time              `bson:"updated_at" json:"updated_at"`
}
