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

// List returns all items for the user's company, optionally filtered by warehouse
func (h *ItemHandler) List(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.Param("id")

	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("items")
	companyObjID, _ := primitive.ObjectIDFromHex(companyID)

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "company_id", Value: companyObjID}, {Key: "is_archived", Value: false}}}},
	}

	// Lookup locations
	pipeline = append(pipeline, bson.D{
		{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "item_locations"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "item_id"},
			{Key: "as", Value: "locations"},
		}},
	})

	if warehouseID != "" {
		whObjID, err := primitive.ObjectIDFromHex(warehouseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
			return
		}

		// Filter locations to specific warehouse and sum quantity
		pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "quantity", Value: bson.D{
				{Key: "$sum", Value: bson.D{
					{Key: "$map", Value: bson.D{
						{Key: "input", Value: bson.D{
							{Key: "$filter", Value: bson.D{
								{Key: "input", Value: "$locations"},
								{Key: "as", Value: "loc"},
								{Key: "cond", Value: bson.D{{Key: "$eq", Value: bson.A{"$$loc.warehouse_id", whObjID}}}},
							}},
						}},
						{Key: "as", Value: "loc"},
						{Key: "in", Value: "$$loc.quantity"},
					}},
				}},
			}},
			{Key: "batch", Value: bson.D{
				{Key: "$let", Value: bson.D{
					{Key: "vars", Value: bson.D{
						{Key: "loc", Value: bson.D{
							{Key: "$arrayElemAt", Value: bson.A{
								bson.D{
									{Key: "$filter", Value: bson.D{
										{Key: "input", Value: "$locations"},
										{Key: "as", Value: "loc"},
										{Key: "cond", Value: bson.D{{Key: "$eq", Value: bson.A{"$$loc.warehouse_id", whObjID}}}},
									}},
								},
								0,
							}},
						}},
					}},
					{Key: "in", Value: "$$loc.batch"},
				}},
			}},
		}}})

		// Only show items present in this warehouse
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.D{{Key: "quantity", Value: bson.D{{Key: "$gt", Value: 0}}}}}})
	} else {
		// Sum all locations
		pipeline = append(pipeline, bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "quantity", Value: bson.D{
				{Key: "$sum", Value: "$locations.quantity"},
			}},
		}}})
	}

	// Remove locations array
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.D{{Key: "locations", Value: 0}}}})

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	defer cursor.Close(ctx)

	var items []bson.M
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
