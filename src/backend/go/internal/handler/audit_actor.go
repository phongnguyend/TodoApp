package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/security"
)

const auditUserIDKey = "auditUserID"

// CaptureOptionalAuditUser records a valid bearer-token user without making
// otherwise-public endpoints require authentication.
func CaptureOptionalAuditUser(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != "" {
			if id, err := security.AuthenticatedUserID(
				c.GetHeader("Authorization"), cfg.JWTSecretKey, time.Now().UTC(),
			); err == nil {
				c.Set(auditUserIDKey, id)
			}
		}
		c.Next()
	}
}

// RequireAuthenticatedAPI protects every API endpoint except the account
// bootstrap and password-recovery endpoints.
func RequireAuthenticatedAPI(cfg *config.Config) gin.HandlerFunc {
	publicEndpoints := map[string]struct{}{
		"/api/tokens":                 {},
		"/api/users/signup":           {},
		"/api/users/password/reset":   {},
		"/api/users/password/confirm": {},
	}

	return func(c *gin.Context) {
		_, public := publicEndpoints[c.Request.URL.Path]
		if public && c.Request.Method == http.MethodPost {
			c.Next()
			return
		}
		if len(c.Request.URL.Path) < 5 || c.Request.URL.Path[:5] != "/api/" {
			c.Next()
			return
		}

		id, err := security.AuthenticatedUserID(
			c.GetHeader("Authorization"), cfg.JWTSecretKey, time.Now().UTC(),
		)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}
		c.Set(auditUserIDKey, id)
		c.Next()
	}
}

func auditUserID(c *gin.Context) *uint {
	value, ok := c.Get(auditUserIDKey)
	if !ok {
		return nil
	}
	id, ok := value.(uint)
	if !ok || id == 0 {
		return nil
	}
	return &id
}
