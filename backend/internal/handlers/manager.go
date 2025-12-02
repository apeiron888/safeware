package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/a2sv/safeware/internal/audit"
	"github.com/a2sv/safeware/internal/auth"
	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ManagerHandler struct {
	auditService *audit.AuditService
}

func NewManagerHandler(auditService *audit.AuditService) *ManagerHandler {
	return &ManagerHandler{
		auditService: auditService,
	}
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
		userID := c.GetString("user_id")
		username := c.GetString("username")

		var req CreateEmployeeRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// If Supervisor, force warehouse ID
		if c.GetString("role") == "Supervisor" {
			req.WarehouseID = c.GetString("warehouse_id")
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
		managerObjectID, _ := primitive.ObjectIDFromHex(userID)

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

		// Log audit
		go h.auditService.LogAction(
			context.Background(),
			managerObjectID,
			companyObjectID,
			username,
			"CREATE",
			"EMPLOYEE",
			&user.ID,
			map[string]interface{}{
				"role":         role,
				"email":        user.Email,
				"warehouse_id": req.WarehouseID,
			},
			c.ClientIP(),
			c.Request.UserAgent(),
			"SUCCESS",
		)

		c.JSON(http.StatusCreated, gin.H{
			"message": role + " created successfully",
			"user_id": user.ID.Hex(),
		})
	}
}

// ListEmployees returns all employees for the company
func (h *ManagerHandler) ListEmployees(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("users")
	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)

	// Find all users for this company, excluding the current user (Manager) if desired,
	// but usually manager wants to see everyone including other managers?
	// Let's filter by roles: Supervisor, Staff, Auditor
	filter := bson.M{
		"company_id": companyObjectID,
		"role":       bson.M{"$in": []string{"Supervisor", "Staff", "Auditor"}},
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

	// Transform to response format if needed, or just return users
	// We might want to sanitize password hashes
	for i := range employees {
		employees[i].PasswordHash = ""
	}

	c.JSON(http.StatusOK, employees)
}

type UpdateEmployeeRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	WarehouseID string `json:"warehouse_id"`
}

// UpdateEmployee handles updating employee details
func (h *ManagerHandler) UpdateEmployee(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")
	employeeID := c.Param("id")

	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employeeObjectID, err := primitive.ObjectIDFromHex(employeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	managerObjectID, _ := primitive.ObjectIDFromHex(userID)

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	// Build update document
	update := bson.M{"$set": bson.M{"updated_at": time.Now()}}
	if req.FullName != "" {
		update["$set"].(bson.M)["full_name"] = req.FullName
	}
	if req.Email != "" {
		update["$set"].(bson.M)["email"] = req.Email
	}
	if req.WarehouseID != "" && c.GetString("role") != "Supervisor" {
		warehouseObjID, err := primitive.ObjectIDFromHex(req.WarehouseID)
		if err == nil {
			update["$set"].(bson.M)["warehouse_id"] = warehouseObjID
		}
	}

	filter := bson.M{"_id": employeeObjectID, "company_id": companyObjectID}
	if c.GetString("role") == "Supervisor" {
		whID := c.GetString("warehouse_id")
		if whID != "" {
			whObjID, _ := primitive.ObjectIDFromHex(whID)
			filter["warehouse_id"] = whObjID
		}
	}

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil || result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Log audit
	go h.auditService.LogAction(
		context.Background(),
		managerObjectID,
		companyObjectID,
		username,
		"UPDATE",
		"EMPLOYEE",
		&employeeObjectID,
		map[string]interface{}{
			"updates": req,
		},
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

// DeleteEmployee handles employee deletion
func (h *ManagerHandler) DeleteEmployee(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.GetString("user_id")
	username := c.GetString("username")
	employeeID := c.Param("id")

	employeeObjectID, err := primitive.ObjectIDFromHex(employeeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	managerObjectID, _ := primitive.ObjectIDFromHex(userID)

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	// Delete the employee
	filter := bson.M{"_id": employeeObjectID, "company_id": companyObjectID}
	if c.GetString("role") == "Supervisor" {
		whID := c.GetString("warehouse_id")
		if whID != "" {
			whObjID, _ := primitive.ObjectIDFromHex(whID)
			filter["warehouse_id"] = whObjID
		}
	}

	result, err := usersCollection.DeleteOne(ctx, filter)
	if err != nil || result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	// Log audit
	go h.auditService.LogAction(
		context.Background(),
		managerObjectID,
		companyObjectID,
		username,
		"DELETE",
		"EMPLOYEE",
		&employeeObjectID,
		nil,
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
