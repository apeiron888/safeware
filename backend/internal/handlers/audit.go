package handlers

import (
	"net/http"

	"github.com/a2sv/safeware/internal/audit"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditHandler struct {
	auditService *audit.AuditService
}

func NewAuditHandler(auditService *audit.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

// List returns a list of audit logs with optional filtering
func (h *AuditHandler) List(c *gin.Context) {
	// Get company ID from context (set by AuthMiddleware)
	companyIDStr := c.GetString("company_id")
	companyID, err := primitive.ObjectIDFromHex(companyIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	// Build filter from query params
	filter := map[string]interface{}{
		"action":        c.Query("action"),
		"resource_type": c.Query("resource_type"),
		"user_id":       c.Query("user_id"),
		"from_date":     c.Query("from_date"),
		"to_date":       c.Query("to_date"),
	}

	logs, err := h.auditService.GetLogs(c.Request.Context(), companyID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve audit logs"})
		return
	}

	c.JSON(http.StatusOK, logs)
}
