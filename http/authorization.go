package http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongkailiu/test-go/swagger/swagger/models"
	log "github.com/sirupsen/logrus"
)

// AuthorizationMiddleware handles if the context is not authorized
func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAuthorized(c) {
			msg := "unauthorized operation"
			c.JSON(http.StatusUnauthorized, models.Error{Code: int64(http.StatusUnauthorized), Message: &msg})
			c.Abort()
			return
		}
	}
}

// IsAuthorized returns true if the context is authorized
func IsAuthorized(c *gin.Context) bool {
	method := c.Request.Method
	path := c.Request.URL.Path
	tokenString, err := getTokenString(c)
	if err != nil {
		log.Warnf("found error when getTokenString(c), %s", err.Error())
	}
	localID, err := getLocalIDFromToken(tokenString, sessionKey)
	if err != nil {
		log.Warnf("found error when getLocalIDFromToken(tokenString, sessionKey), %s", err.Error())
	}
	return isAuthorized(localID, method, path)
}

func getTokenString(c *gin.Context) (string, error) {
	headerValue := c.GetHeader("Authorization")
	if strings.HasPrefix(headerValue, "Bearer ") {
		return strings.TrimPrefix(headerValue, "Bearer "), nil
	}
	return "", fmt.Errorf("malformed Authorization header: %s", headerValue)
}

func isAuthorized(localID, method, path string) bool {
	if localID == "" {
		return false
	}
	if strings.HasPrefix(path, "/api/v1") {
		log.WithFields(log.Fields{"localID": localID, "method": method, "path": path}).Debug("isAuthorized")
		// TODO Implement Role based authorization with storage
		// Simulate authorization
		if strings.HasPrefix(path, "/api/v1/user/") {
			return true
		}
		if strings.HasPrefix(path, "/api/v1/users") {
			return false
		}
	}
	return true
}
