package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongkailiu/test-go/swagger/swagger/models"
	log "github.com/sirupsen/logrus"
)

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

func IsAuthorized(c *gin.Context) bool {
	method := c.Request.Method
	path := c.Request.URL.Path
	aHeader := c.GetHeader("Authorization")
	log.WithFields(log.Fields{"aHeader": aHeader, "method": method, "path": path}).Debug("isAuthorized")
	return isAuthorized(getLocalID(aHeader), method, path, c)
}
func getLocalID(s string) string {
	if strings.HasPrefix(s, "Bearer ") {
		tokenString := strings.TrimPrefix(s, "Bearer ")
		localID, err := getLocalIDFromToken(tokenString, sessionKey)
		if err != nil {
			log.Warnf("found error when getLocalIDFromToken(s), %s", err.Error())
			return ""
		}
		return localID
	}
	log.Warnf("malformed Authorization header: %s", s)
	return ""
}

func isAuthorized(localID, method, path string, c *gin.Context) bool {
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
