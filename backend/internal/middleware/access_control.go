package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TimeEnforcementMiddleware enforces 8AM-6PM access for non-managers
func TimeEnforcementMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role == "Manager" {
			c.Next()
			return
		}

		now := time.Now()
		hour := now.Hour()

		// 8:00 AM to 6:00 PM (18:00)
		if hour < 8 || hour >= 18 {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. System is only available between 8:00 AM and 6:00 PM.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleEnforcementMiddleware ensures user has one of the allowed roles
func RoleEnforcementMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions. Required role: " + userRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// WarehouseEnforcementMiddleware ensures Staff/Supervisor only access their assigned warehouse
func WarehouseEnforcementMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		// Managers and Auditors can access any warehouse
		if role == "Manager" || role == "Auditor" {
			c.Next()
			return
		}

		// Get user's assigned warehouse
		userWarehouseID := c.GetString("warehouse_id")
		if userWarehouseID == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "No warehouse assigned to user"})
			c.Abort()
			return
		}

		// Check if request is targeting a specific warehouse
		// This can be in URL param :warehouse_id or :id (if resource is warehouse)
		// or in the JSON body

		// For simple URL param check
		requestedWarehouseID := c.Param("warehouse_id")
		if requestedWarehouseID == "" {
			// Try to get from query or body if needed, but for now strict URL param
			// If the resource itself implies a warehouse (like item), we need to check that item's warehouse
			// This might be complex for generic middleware, so we might need handler-level checks too.
			// For now, let's assume if warehouse_id is present in path, we check it.
			c.Next()
			return
		}

		if requestedWarehouseID != userWarehouseID {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. You can only access your assigned warehouse.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
