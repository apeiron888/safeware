package handlers

import (
	"context"
	"net/http"

	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ListWarehouseEmployees returns employees for the requester's warehouse
func (h *ManagerHandler) ListWarehouseEmployees(c *gin.Context) {
	companyID := c.GetString("company_id")
	warehouseID := c.GetString("warehouse_id")

	if companyID == "" || warehouseID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or no warehouse assigned"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("users")
	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	warehouseObjectID, _ := primitive.ObjectIDFromHex(warehouseID)

	// Find users in this warehouse with roles Staff or Supervisor
	filter := bson.M{
		"company_id":   companyObjectID,
		"warehouse_id": warehouseObjectID,
		"role":         bson.M{"$in": []string{"Supervisor", "Staff"}},
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch employees"})
		return
	}
	defer cursor.Close(ctx)

	var employees []models.User
	if err = cursor.All(ctx, &employees); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode employees"})
		return
	}

	// Sanitize password hashes
	for i := range employees {
		employees[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, employees)
}
