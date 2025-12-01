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

type PermissionHandler struct{}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{}
}

type CreatePermissionRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	ResourceType string `json:"resource_type"`
}

// ListPermissions returns all permissions
func (h *PermissionHandler) List(c *gin.Context) {
	ctx := context.Background()
	collection := database.GetCollection("permissions")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch permissions"})
		return
	}
	defer cursor.Close(ctx)

	var permissions []models.Permission
	if err = cursor.All(ctx, &permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode permissions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// Get returns a single permission by ID
func (h *PermissionHandler) Get(c *gin.Context) {
	permissionID := c.Param("id")

	ctx := context.Background()
	collection := database.GetCollection("permissions")

	objectID, err := primitive.ObjectIDFromHex(permissionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid permission ID"})
		return
	}

	var permission models.Permission
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&permission)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Permission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, permission)
}

// Create creates a new permission (admin only)
func (h *PermissionHandler) Create(c *gin.Context) {
	var req CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	permission := models.Permission{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Description:  req.Description,
		ResourceType: req.ResourceType,
		CreatedAt:    time.Now(),
	}

	ctx := context.Background()
	collection := database.GetCollection("permissions")

	_, err := collection.InsertOne(ctx, permission)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create permission"})
		return
	}

	c.JSON(http.StatusCreated, permission)
}
