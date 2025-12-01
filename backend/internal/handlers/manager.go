package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/a2sv/safeware/internal/auth"
	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ManagerHandler struct{}

func NewManagerHandler() *ManagerHandler {
	return &ManagerHandler{}
}

type CreateEmployeeRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	WarehouseID string `json:"warehouse_id"` // Required for Staff/Supervisor
}

// CreateEmployee handles creation of Supervisor, Staff, or Auditor
func (h *ManagerHandler) CreateEmployee(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		companyID := c.GetString("company_id")

		var req CreateEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate warehouse for Staff/Supervisor
		var warehouseObjectID *primitive.ObjectID
		if role == "Staff" || role == "Supervisor" {
			if req.WarehouseID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Warehouse ID is required for " + role})
				return
			}
			objID, err := primitive.ObjectIDFromHex(req.WarehouseID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
				return
			}
			warehouseObjectID = &objID

			// Verify warehouse exists and belongs to company
			ctx := context.Background()
			whCollection := database.GetCollection("warehouses")
			companyObjID, _ := primitive.ObjectIDFromHex(companyID)
			count, _ := whCollection.CountDocuments(ctx, bson.M{"_id": objID, "company_id": companyObjID})
			if count == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Warehouse not found or does not belong to company"})
				return
			}

			// For Supervisor, check if warehouse already has one
			if role == "Supervisor" {
				// This check should ideally be on the Warehouse model's SupervisorID field
				// But for now checking users is okay too, or we update warehouse later
			}
		}

		// Hash password
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
			return
		}

		companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

		user := models.User{
			ID:             primitive.NewObjectID(),
			CompanyID:      companyObjectID,
			FullName:       req.FullName,
			Email:          req.Email,
			PasswordHash:   hashedPassword,
			Role:           role,
			WarehouseID:    warehouseObjectID,
			ClearanceLevel: 0,
			IsVerified:     true, // Manager-created employees are auto-verified
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		ctx := context.Background()
		usersCollection := database.GetCollection("users")

		// Check email uniqueness
		count, _ := usersCollection.CountDocuments(ctx, bson.M{"email": req.Email})
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		_, err = usersCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
			return
		}

		// If Supervisor, update Warehouse
		if role == "Supervisor" && warehouseObjectID != nil {
			whCollection := database.GetCollection("warehouses")
			whCollection.UpdateOne(ctx, bson.M{"_id": *warehouseObjectID}, bson.M{"$set": bson.M{"supervisor_id": user.ID}})
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": role + " created successfully",
			"user_id": user.ID.Hex(),
		})
	}
}
