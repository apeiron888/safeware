package middleware

import (
	"bytes"
	"io"
	"strings"

	"github.com/a2sv/safeware/internal/audit"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditMiddleware intercepts requests and logs them
func AuditMiddleware(auditService *audit.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read body for logging (and restore it for handlers)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Process request
		c.Next()

		// Extract info after request is processed
		status := "SUCCESS"
		if c.Writer.Status() >= 400 {
			status = "FAILURE"
		}

		// Get user info from context (set by AuthMiddleware)
		userIDStr := c.GetString("user_id")
		companyIDStr := c.GetString("company_id")
		email := c.GetString("email")

		// Skip logging for unauthenticated requests (unless we want to log login failures, which is good practice)
		// For MVP, we'll log if we have a user ID, or if it's a login/register attempt
		path := c.Request.URL.Path
		method := c.Request.Method

		// Determine Action
		action := method
		if strings.Contains(path, "/login") {
			action = "LOGIN"
		} else if strings.Contains(path, "/register") {
			action = "REGISTER"
		} else if strings.Contains(path, "/logout") {
			action = "LOGOUT"
		} else {
			// Map HTTP methods to actions
			switch method {
			case "POST":
				action = "CREATE"
			case "PUT", "PATCH":
				action = "UPDATE"
			case "DELETE":
				action = "DELETE"
			case "GET":
				action = "READ"
			}
		}

		// Determine Resource Type
		resourceType := "API"
		if strings.Contains(path, "/items") || strings.Contains(path, "/item") {
			resourceType = "ITEM"
		} else if strings.Contains(path, "/warehouses") || strings.Contains(path, "/warehouse") {
			resourceType = "WAREHOUSE"
		} else if strings.Contains(path, "/users") || strings.Contains(path, "/staff") || strings.Contains(path, "/supervisor") {
			resourceType = "USER"
		} else if strings.Contains(path, "/roles") {
			resourceType = "ROLE"
		}

		// Prepare details
		details := map[string]interface{}{
			"path":   path,
			"method": method,
			"status": c.Writer.Status(),
		}

		// Parse IDs
		var userObjID, companyObjID primitive.ObjectID
		if userIDStr != "" {
			userObjID, _ = primitive.ObjectIDFromHex(userIDStr)
		}
		if companyIDStr != "" {
			companyObjID, _ = primitive.ObjectIDFromHex(companyIDStr)
		}

		// Log it
		// If no user ID (e.g. failed login), we log with empty ID but capture IP/Email if possible
		auditService.LogAction(
			c.Request.Context(),
			userObjID,
			companyObjID,
			email,
			action,
			resourceType,
			nil, // Resource ID could be extracted from params if needed
			details,
			c.ClientIP(),
			c.Request.UserAgent(),
			status,
		)
	}
}
