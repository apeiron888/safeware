package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/a2sv/safeware/internal/audit"
	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WarehouseHandler struct {
	auditService *audit.AuditService
}

func NewWarehouseHandler(auditService *audit.AuditService) *WarehouseHandler {
	return &WarehouseHandler{
		auditService: auditService,
	}
}

type CreateWarehouseRequest struct {
	Name        string   `json:"name" binding:"required"`
	Location    string   `json:"location" binding:"required"`
	IPWhitelist []string `json:"ip_whitelist"`
}

type UpdateWarehouseRequest struct {
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	IPWhitelist []string `json:"ip_whitelist"`
	IsActive    *bool    `json:"is_active"`
}

// List returns all warehouses for the user's company
func (h *WarehouseHandler) List(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("warehouses")

	objectID, _ := primitive.ObjectIDFromHex(companyID)
	cursor, err := collection.Find(ctx, bson.M{"company_id": objectID, "is_active": true})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch warehouses"})
		return
	}
	defer cursor.Close(ctx)

	var warehouses []models.Warehouse
	if err = cursor.All(ctx, &warehouses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode warehouses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"warehouses": warehouses})
}

// Get returns a single warehouse by ID
func (h *WarehouseHandler) Get(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")
	warehouseID := c.Param("id")

	ctx := context.Background()
	collection := database.GetCollection("warehouses")

	objectID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	var warehouse models.Warehouse
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}).Decode(&warehouse)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Log audit for reading warehouse details
	go h.auditService.LogAction(
		context.Background(),
		userObjectID,
		companyObjectID,
		username,
		"READ",
		"WAREHOUSE",
		&objectID,
		map[string]interface{}{
			"warehouse_name": warehouse.Name,
		},
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	c.JSON(http.StatusOK, warehouse)
}

// Create creates a new warehouse
func (h *WarehouseHandler) Create(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")

	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	warehouse := models.Warehouse{
		ID:          primitive.NewObjectID(),
		CompanyID:   companyObjectID,
		Name:        req.Name,
		Location:    req.Location,
		IPWhitelist: req.IPWhitelist,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ctx := context.Background()
	collection := database.GetCollection("warehouses")

	_, err := collection.InsertOne(ctx, warehouse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create warehouse"})
		return
	}

	// Log audit
	go h.auditService.LogAction(
		context.Background(),
		userObjectID,
		companyObjectID,
		username,
		"CREATE",
		"WAREHOUSE",
		&warehouse.ID,
		map[string]interface{}{
			"name":     warehouse.Name,
			"location": warehouse.Location,
		},
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	c.JSON(http.StatusCreated, warehouse)
}

// Update updates an existing warehouse
func (h *WarehouseHandler) Update(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")
	warehouseID := c.Param("id")

	var req UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Location != "" {
		update["$set"].(bson.M)["location"] = req.Location
	}
	if req.IPWhitelist != nil {
		update["$set"].(bson.M)["ip_whitelist"] = req.IPWhitelist
	}
	if req.IsActive != nil {
		update["$set"].(bson.M)["is_active"] = *req.IsActive
	}

	ctx := context.Background()
	collection := database.GetCollection("warehouses")

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	// Fetch updated warehouse
	var warehouse models.Warehouse
	collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&warehouse)

	// Log audit
	go h.auditService.LogAction(
		context.Background(),
		userObjectID,
		companyObjectID,
		username,
		"UPDATE",
		"WAREHOUSE",
		&objectID,
		map[string]interface{}{
			"updates": req,
		},
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	c.JSON(http.StatusOK, warehouse)
}

// Delete soft deletes a warehouse
func (h *WarehouseHandler) Delete(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")
	warehouseID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(warehouseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	userObjectID, _ := primitive.ObjectIDFromHex(userID)

	ctx := context.Background()
	collection := database.GetCollection("warehouses")

	update := bson.M{"$set": bson.M{"is_active": false, "updated_at": time.Now()}}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	// Log audit
	go h.auditService.LogAction(
		context.Background(),
		userObjectID,
		companyObjectID,
		username,
		"DELETE",
		"WAREHOUSE",
		&objectID,
		nil,
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
}
