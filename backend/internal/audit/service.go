package audit

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"

	"github.com/a2sv/safeware/internal/database"
	"github.com/a2sv/safeware/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuditService struct {
	encryptionKey []byte
}

// NewAuditService creates a new audit service
// key must be 32 bytes for AES-256
func NewAuditService(key string) *AuditService {
	// Ensure key is 32 bytes, pad or trim if necessary (for MVP simplicity)
	// In production, this should be validated strictly
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		keyBytes = padded
	} else if len(keyBytes) > 32 {
		keyBytes = keyBytes[:32]
	}

	return &AuditService{
		encryptionKey: keyBytes,
	}
}

// LogAction records an action asynchronously
func (s *AuditService) LogAction(ctx context.Context, actorID, companyID primitive.ObjectID, username, action, resourceType string, resourceID *primitive.ObjectID, details map[string]interface{}, ip, userAgent string, status string) {
	// Run in goroutine to not block the main request
	go func() {
		// Create a detached context with timeout for the background operation
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		encryptedDetails := ""
		if details != nil {
			jsonBytes, err := json.Marshal(details)
			if err == nil {
				encrypted, err := s.encrypt(jsonBytes)
				if err == nil {
					encryptedDetails = encrypted
				} else {
					log.Printf("Error encrypting audit details: %v", err)
				}
			}
		}

		logEntry := models.AuditLog{
			ID:               primitive.NewObjectID(),
			CompanyID:        companyID,
			UserID:           actorID,
			Username:         username,
			Action:           action,
			ResourceType:     resourceType,
			Status:           status,
			DetailsEncrypted: encryptedDetails,
			IPAddress:        ip,
			UserAgent:        userAgent,
			CreatedAt:        time.Now(),
		}

		if resourceID != nil {
			logEntry.ResourceID = *resourceID
		}

		collection := database.GetCollection("audit_logs")
		_, err := collection.InsertOne(bgCtx, logEntry)
		if err != nil {
			log.Printf("Error writing audit log: %v", err)
		}
	}()
}

// encrypt encrypts data using AES-GCM
func (s *AuditService) encrypt(data []byte) (string, error) {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts audit details (for viewing logs)
func (s *AuditService) Decrypt(encryptedString string) (map[string]interface{}, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedString)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	var details map[string]interface{}
	if err := json.Unmarshal(plaintext, &details); err != nil {
		return nil, err
	}

	return details, nil
}

// GetLogs retrieves audit logs with optional filtering
func (s *AuditService) GetLogs(ctx context.Context, companyID primitive.ObjectID, filter map[string]interface{}) ([]map[string]interface{}, error) {
	collection := database.GetCollection("audit_logs")

	// Base filter: Company ID
	query := map[string]interface{}{
		"company_id": companyID,
	}

	// Apply additional filters with case-insensitive regex
	if action, ok := filter["action"].(string); ok && action != "" {
		query["action"] = map[string]interface{}{
			"$regex":   action,
			"$options": "i", // case-insensitive
		}
	}
	if resourceType, ok := filter["resource_type"].(string); ok && resourceType != "" {
		query["resource_type"] = map[string]interface{}{
			"$regex":   resourceType,
			"$options": "i", // case-insensitive
		}
	}
	if userID, ok := filter["user_id"].(string); ok && userID != "" {
		oid, err := primitive.ObjectIDFromHex(userID)
		if err == nil {
			query["user_id"] = oid
		}
	}

	// Date range filter
	if fromStr, ok := filter["from_date"].(string); ok && fromStr != "" {
		from, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			if _, exists := query["created_at"]; !exists {
				query["created_at"] = map[string]interface{}{}
			}
			query["created_at"].(map[string]interface{})["$gte"] = from
		}
	}
	if toStr, ok := filter["to_date"].(string); ok && toStr != "" {
		to, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			if _, exists := query["created_at"]; !exists {
				query["created_at"] = map[string]interface{}{}
			}
			query["created_at"].(map[string]interface{})["$lte"] = to
		}
	}

	// Create find options with sorting (most recent first)
	findOptions := options.Find()
	findOptions.SetSort(map[string]interface{}{"created_at": -1}) // Sort descending

	cursor, err := collection.Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}

	// Decrypt details for each log
	result := make([]map[string]interface{}, 0, len(logs))
	for _, logEntry := range logs {
		entry := map[string]interface{}{
			"id":            logEntry.ID,
			"user_id":       logEntry.UserID,
			"username":      logEntry.Username,
			"action":        logEntry.Action,
			"resource_type": logEntry.ResourceType,
			"resource_id":   logEntry.ResourceID,
			"status":        logEntry.Status,
			"ip_address":    logEntry.IPAddress,
			"user_agent":    logEntry.UserAgent,
			"timestamp":     logEntry.CreatedAt, // Changed to "timestamp" to match frontend
		}

		if logEntry.DetailsEncrypted != "" {
			details, err := s.Decrypt(logEntry.DetailsEncrypted)
			if err == nil {
				entry["details"] = details
			} else {
				entry["details_error"] = "Failed to decrypt details"
			}
		}

		result = append(result, entry)
	}

	return result, nil
}
