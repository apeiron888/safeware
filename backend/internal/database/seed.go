package database

import (
	"context"
	"log"
	"time"

	"github.com/a2sv/safeware/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SeedDefaultPermissions creates default system permissions
func SeedDefaultPermissions() error {
	ctx := context.Background()
	collection := GetCollection("permissions")

	// Check if permissions already exist
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Permissions already seeded, skipping...")
		return nil
	}

	// Default permissions
	permissions := []interface{}{
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "users.create",
			Description:  "Create new users",
			ResourceType: "user",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "users.read",
			Description:  "View users",
			ResourceType: "user",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "users.update",
			Description:  "Update users",
			ResourceType: "user",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "users.delete",
			Description:  "Delete users",
			ResourceType: "user",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "roles.manage",
			Description:  "Manage roles",
			ResourceType: "role",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "warehouses.create",
			Description:  "Create warehouses",
			ResourceType: "warehouse",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "warehouses.read",
			Description:  "View warehouses",
			ResourceType: "warehouse",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "warehouses.update",
			Description:  "Update warehouses",
			ResourceType: "warehouse",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "warehouses.delete",
			Description:  "Delete warehouses",
			ResourceType: "warehouse",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "items.create",
			Description:  "Create items",
			ResourceType: "item",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "items.read",
			Description:  "View items",
			ResourceType: "item",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "items.update",
			Description:  "Update items",
			ResourceType: "item",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "items.delete",
			Description:  "Delete items",
			ResourceType: "item",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "transfers.approve",
			Description:  "Approve transfers",
			ResourceType: "transfer",
			CreatedAt:    time.Now(),
		},
		models.Permission{
			ID:           primitive.NewObjectID(),
			Name:         "audit.read",
			Description:  "View audit logs",
			ResourceType: "audit",
			CreatedAt:    time.Now(),
		},
	}

	_, err = collection.InsertMany(ctx, permissions)
	if err != nil {
		return err
	}

	log.Println("✅ Seeded 15 default permissions")
	return nil
}

// SeedDefaultRoles creates default roles for a company
func SeedDefaultRoles(companyID primitive.ObjectID) error {
	ctx := context.Background()
	rolesCollection := GetCollection("roles")
	permissionsCollection := GetCollection("permissions")

	// Check if roles already exist for this company
	count, err := rolesCollection.CountDocuments(ctx, bson.M{"company_id": companyID})
	if err != nil {
		return err
	}
	if count > 0 {
		log.Printf("Roles already exist for company %s, skipping...\n", companyID.Hex())
		return nil
	}

	// Get all permissions
	cursor, err := permissionsCollection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	var allPermissions []models.Permission
	if err = cursor.All(ctx, &allPermissions); err != nil {
		return err
	}

	// Create permission ID maps
	var allPermissionIDs []primitive.ObjectID
	readPermissions := []primitive.ObjectID{}
	writePermissions := []primitive.ObjectID{}

	for _, p := range allPermissions {
		allPermissionIDs = append(allPermissionIDs, p.ID)
		if p.Name == "items.read" || p.Name == "warehouses.read" || p.Name == "users.read" {
			readPermissions = append(readPermissions, p.ID)
		}
		if p.Name == "items.create" || p.Name == "items.update" ||
			p.Name == "warehouses.create" || p.Name == "warehouses.update" {
			writePermissions = append(writePermissions, p.ID)
		}
	}

	// Default roles
	roles := []interface{}{
		models.Role{
			ID:             primitive.NewObjectID(),
			CompanyID:      companyID,
			Name:           "Manager",
			Description:    "Full access to all resources",
			HierarchyLevel: 3,
			PermissionIDs:  allPermissionIDs,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		models.Role{
			ID:             primitive.NewObjectID(),
			CompanyID:      companyID,
			Name:           "Staff",
			Description:    "Can create and manage items and warehouses",
			HierarchyLevel: 2,
			PermissionIDs:  append(readPermissions, writePermissions...),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
		models.Role{
			ID:             primitive.NewObjectID(),
			CompanyID:      companyID,
			Name:           "Auditor",
			Description:    "Read-only access for auditing purposes",
			HierarchyLevel: 1,
			PermissionIDs:  readPermissions,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	_, err = rolesCollection.InsertMany(ctx, roles)
	if err != nil {
		return err
	}

	log.Printf("✅ Seeded 3 default roles for company %s\n", companyID.Hex())
	return nil
}

// GetCompanyRoles returns all roles for a company
func GetCompanyRoles(companyID primitive.ObjectID) ([]models.Role, error) {
	ctx := context.Background()
	collection := GetCollection("roles")

	cursor, err := collection.Find(ctx, bson.M{"company_id": companyID})
	if err != nil {
		return nil, err
	}

	var roles []models.Role
	if err = cursor.All(ctx, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}
