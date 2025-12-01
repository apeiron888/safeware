package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoleHandler struct{}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{}
}

type CreateRoleRequest struct {
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	HierarchyLevel int      `json:"hierarchy_level"`
	PermissionIDs  []string `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	HierarchyLevel *int     `json:"hierarchy_level"`
	PermissionIDs  []string `json:"permission_ids"`
}

type AssignRoleRequest struct {
	UserID string `json:"user_id" binding:"required"`
	RoleID string `json:"role_id" binding:"required"`
}

// ListRoles returns all roles for the company
func (h *RoleHandler) List(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("roles")

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	cursor, err := collection.Find(ctx, bson.M{"company_id": companyObjectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch roles"})
		return
	}
	defer cursor.Close(ctx)

	var roles []models.Role
	if err = cursor.All(ctx, &roles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// Get returns a single role by ID
func (h *RoleHandler) Get(c *gin.Context) {
	companyID := c.GetString("company_id")
	roleID := c.Param("id")

	ctx := context.Background()
	collection := database.GetCollection("roles")

	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	var role models.Role
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}).Decode(&role)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Create creates a new role
func (h *RoleHandler) Create(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	// Convert permission IDs
	var permissionIDs []primitive.ObjectID
	for _, idStr := range req.PermissionIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err == nil {
			permissionIDs = append(permissionIDs, id)
		}
	}

	role := models.Role{
		ID:             primitive.NewObjectID(),
		CompanyID:      companyObjectID,
		Name:           req.Name,
		Description:    req.Description,
		HierarchyLevel: req.HierarchyLevel,
		PermissionIDs:  permissionIDs,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	ctx := context.Background()
	collection := database.GetCollection("roles")

	_, err := collection.InsertOne(ctx, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create role"})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// Update updates an existing role
func (h *RoleHandler) Update(c *gin.Context) {
	companyID := c.GetString("company_id")
	roleID := c.Param("id")

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Description != "" {
		update["$set"].(bson.M)["description"] = req.Description
	}
	if req.HierarchyLevel != nil {
		update["$set"].(bson.M)["hierarchy_level"] = *req.HierarchyLevel
	}
	if req.PermissionIDs != nil {
		var permissionIDs []primitive.ObjectID
		for _, idStr := range req.PermissionIDs {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				permissionIDs = append(permissionIDs, id)
			}
		}
		update["$set"].(bson.M)["permission_ids"] = permissionIDs
	}

	ctx := context.Background()
	collection := database.GetCollection("roles")

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	var role models.Role
	collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&role)

	c.JSON(http.StatusOK, role)
}

// Delete deletes a role
func (h *RoleHandler) Delete(c *gin.Context) {
	companyID := c.GetString("company_id")
	roleID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	ctx := context.Background()
	collection := database.GetCollection("roles")

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID})
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role deleted successfully"})
}

// AssignRole assigns a role to a user
func (h *RoleHandler) AssignRole(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleObjectID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	// Add role to user's role_ids array
	update := bson.M{
		"$addToSet": bson.M{"role_ids": roleObjectID},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// RemoveRole removes a role from a user
func (h *RoleHandler) RemoveRole(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userObjectID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	roleObjectID, err := primitive.ObjectIDFromHex(req.RoleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	// Remove role from user's role_ids array
	update := bson.M{
		"$pull": bson.M{"role_ids": roleObjectID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	result, err := usersCollection.UpdateOne(ctx, bson.M{"_id": userObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}
