package handlers

import (
	"context"
	"net/http"

	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// ListUsers returns all users for the company
func (h *UserHandler) List(c *gin.Context) {
	companyID := c.GetString("company_id")
	if companyID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	collection := database.GetCollection("users")

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	cursor, err := collection.Find(ctx, bson.M{"company_id": companyObjectID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode users"})
		return
	}

	// Remove sensitive fields
	for i := range users {
		users[i].PasswordHash = ""
		users[i].TOTPSecret = ""
		users[i].VerificationToken = ""
		users[i].ResetPasswordToken = ""
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// GetUser returns a single user by ID
func (h *UserHandler) Get(c *gin.Context) {
	companyID := c.GetString("company_id")
	userID := c.Param("id")

	ctx := context.Background()
	collection := database.GetCollection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	companyObjectID, _ := primitive.ObjectIDFromHex(companyID)
	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID, "company_id": companyObjectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Remove sensitive fields
	user.PasswordHash = ""
	user.TOTPSecret = ""
	user.VerificationToken = ""
	user.ResetPasswordToken = ""

	c.JSON(http.StatusOK, user)
}
