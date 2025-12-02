package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/a2sv/safeware/internal/audit"
	"github.com/a2sv/safeware/internal/auth"
	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/email"
	"github.com/a2sv/safeware/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	jwtService   *auth.JWTService
	emailService *email.EmailService
	auditService *audit.AuditService
}

func NewAuthHandler(jwtService *auth.JWTService, emailService *email.EmailService, auditService *audit.AuditService) *AuthHandler {
	return &AuthHandler{
		jwtService:   jwtService,
		emailService: emailService,
		auditService: auditService,
	}
}

// RegisterRequest represents registration payload
type RegisterRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	CompanyName string `json:"company_name"` // Optional - creates new company if provided
}

// LoginRequest represents login payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// TokenResponse represents auth token response
type TokenResponse struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	User         *UserResponse `json:"user"`
}

// UserResponse represents user data in responses
type UserResponse struct {
	ID             string `json:"id"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	CompanyID      string `json:"company_id"`
	Role           string `json:"role"`
	WarehouseID    string `json:"warehouse_id,omitempty"`
	ClearanceLevel int    `json:"clearance_level"`
	IsVerified     bool   `json:"is_verified"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password policy
	if err := auth.ValidatePasswordPolicy(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 12 characters with uppercase, lowercase, number, and symbol"})
		return
	}

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	// Check if user already exists
	var existingUser models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	var companyID primitive.ObjectID
	if req.CompanyName != "" {
		// Create new company
		companiesCollection := database.GetCollection("companies")
		company := models.Company{
			ID:        primitive.NewObjectID(),
			Name:      req.CompanyName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := companiesCollection.InsertOne(ctx, company)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
			return
		}
		companyID = company.ID
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Company name is required"})
		return
	}

	// Generate verification token
	verificationToken, err := auth.GenerateRandomToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verification token"})
		return
	}

	tokenExpiry := time.Now().Add(24 * time.Hour)

	// Create user
	user := models.User{
		ID:                       primitive.NewObjectID(),
		CompanyID:                companyID,
		FullName:                 req.FullName,
		Email:                    req.Email,
		PasswordHash:             hashedPassword,
		Role:                     "Manager", // First user is always Manager
		ClearanceLevel:           0,
		IsVerified:               false,
		VerificationToken:        verificationToken,
		VerificationTokenExpires: &tokenExpiry,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Send verification email (async)
	go h.emailService.SendVerificationEmail(user.Email, verificationToken)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful! Please check your email to verify your account.",
		"user_id": user.ID.Hex(),
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	// Find user
	var user models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Verify password
	if !auth.VerifyPassword(user.PasswordHash, req.Password) {
		// Log failed login audit
		go h.auditService.LogAction(
			context.Background(),
			user.ID,
			user.CompanyID,
			user.FullName,
			"LOGIN",
			"USER",
			&user.ID,
			map[string]interface{}{
				"email": user.Email,
				"role":  user.Role,
			},
			c.ClientIP(),
			c.Request.UserAgent(),
			"FAILED",
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate tokens
	accessToken, err := h.jwtService.GenerateAccessToken(
		user.ID.Hex(),
		user.CompanyID.Hex(),
		user.Email,
		user.Role,
		user.WarehouseID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := h.jwtService.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// Update last login
	now := time.Now()
	usersCollection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{
		"$set": bson.M{"last_login": now},
	})

	// Log successful login audit
	go h.auditService.LogAction(
		context.Background(),
		user.ID,
		user.CompanyID,
		user.FullName,
		"LOGIN",
		"USER",
		&user.ID,
		map[string]interface{}{
			"email": user.Email,
			"role":  user.Role,
		},
		c.ClientIP(),
		c.Request.UserAgent(),
		"SUCCESS",
	)

	var warehouseID string
	if user.WarehouseID != nil {
		warehouseID = user.WarehouseID.Hex()
	}

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: &UserResponse{
			ID:             user.ID.Hex(),
			FullName:       user.FullName,
			Email:          user.Email,
			CompanyID:      user.CompanyID.Hex(),
			Role:           user.Role,
			WarehouseID:    warehouseID,
			ClearanceLevel: user.ClearanceLevel,
			IsVerified:     user.IsVerified,
		},
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Get user
	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	userID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Generate new access token
	accessToken, err := h.jwtService.GenerateAccessToken(
		user.ID.Hex(),
		user.CompanyID.Hex(),
		user.Email,
		user.Role,
		user.WarehouseID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

// GetProfile returns current user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx := context.Background()
	usersCollection := database.GetCollection("users")

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	err = usersCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
