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

type ItemHandler struct{}

func NewItemHandler() *ItemHandler {
	return &ItemHandler{}
}

type CreateItemRequest struct {
	SKU         string                 `json:"sku" binding:"required"`
	Name        string                 `json:"name" binding:"required"`
	Quality     string                 `json:"quality" binding:"required"` // New, Used, Damaged
	Price       float64                `json:"price"`
	Department  string                 `json:"department"`
	Attributes  map[string]interface{} `json:"attributes"`
	WarehouseID string                 `json:"warehouse_id" binding:"required"`
	Quantity    int                    `json:"quantity" binding:"required,min=0"`
	Batch       string                 `json:"batch"`
}

type UpdateItemRequest struct {
	Name       string                 `json:"name"`
	Quality    string                 `json:"quality"`
	Price      *float64               `json:"price"`
	Department string                 `json:"department"`
	Attributes map[string]interface{} `json:"attributes"`
}

// List returns all items for the user's company
func (h *ItemHandler) List(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("items")

	filter := bson.M{"company_id": companyID, "is_archived": false}

	// Apply filters
	if classification := c.Query("classification"); classification != "" {
		filter["classification"] = classification
	}
	if search := c.Query("search"); search != "" {
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"sku": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer cursor.Close(ctx)

	var items []models.Item
	if err = cursor.All(ctx, &items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode items"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": items})
}

// Get returns a single item by ID
func (h *ItemHandler) Get(c *gin.Context) {
	companyID := c.GetString("company_id")
	itemID := c.Param("id")

	ctx := context.Background()
	collection := database.GetCollection("items")

	objectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	var item models.Item
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Get item locations
	locationsCollection := database.GetCollection("item_locations")
	cursor, _ := locationsCollection.Find(ctx, bson.M{"item_id": objectID})
	var locations []models.ItemLocation
	cursor.All(ctx, &locations)

	c.JSON(http.StatusOK, gin.H{
		"item":      item,
		"locations": locations,
	})
}

// Create creates a new item
func (h *ItemHandler) Create(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	ownerObjectID, _ := primitive.ObjectIDFromHex(userID)
	warehouseObjectID, err := primitive.ObjectIDFromHex(req.WarehouseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
		return
	}

	ctx := context.Background()

	item := models.Item{
		ID:          primitive.NewObjectID(),
		CompanyID:   companyObjectID,
		SKU:         req.SKU,
		Name:        req.Name,
		Quality:     req.Quality,
		Price:       req.Price,
		OwnerUserID: ownerObjectID,
		Department:  req.Department,
		Attributes:  req.Attributes,
		IsArchived:  false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Insert item
	itemsCollection := database.GetCollection("items")
	_, err = itemsCollection.InsertOne(ctx, item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	// Create initial location
	location := models.ItemLocation{
		ID:          primitive.NewObjectID(),
		ItemID:      item.ID,
		WarehouseID: warehouseObjectID,
		Quantity:    req.Quantity,
		Batch:       req.Batch,
		UpdatedBy:   ownerObjectID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	locationsCollection := database.GetCollection("item_locations")
	_, err = locationsCollection.InsertOne(ctx, location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item location"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// Update updates an existing item
func (h *ItemHandler) Update(c *gin.Context) {
	companyID := c.GetString("company_id")
	itemID := c.Param("id")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	if req.Name != "" {
		update["$set"].(bson.M)["name"] = req.Name
	}
	if req.Quality != "" {
		update["$set"].(bson.M)["quality"] = req.Quality
	}
	if req.Price != nil {
		update["$set"].(bson.M)["price"] = *req.Price
	}
	if req.Department != "" {
		update["$set"].(bson.M)["department"] = req.Department
	}
	if req.Attributes != nil {
		update["$set"].(bson.M)["attributes"] = req.Attributes
	}

	ctx := context.Background()
	collection := database.GetCollection("items")

	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Fetch updated item
	var item models.Item
	collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)

	c.JSON(http.StatusOK, item)
}

// Delete archives an item
func (h *ItemHandler) Delete(c *gin.Context) {
	companyID := c.GetString("company_id")
	itemID := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	ctx := context.Background()
	collection := database.GetCollection("items")

	update := bson.M{"$set": bson.M{"is_archived": true, "updated_at": time.Now()}}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item archived successfully"})
}
