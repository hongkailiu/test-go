package http

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/hongkailiu/test-go/swagger/swagger/models"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !IsAuthenticated(c) {
			msg := "authentication is required"
			c.JSON(http.StatusUnauthorized, models.Error{Code: int64(http.StatusUnauthorized), Message: &msg})
			c.Abort()
			return
		}
	}
}

func IsAuthenticated(c *gin.Context) bool {
	return WhoAmI(c, "username") != nil
}

func WhoAmI(c *gin.Context, key string) *string {
	if value, exists := c.Get(sessions.DefaultKey); exists {
		session := value.(sessions.Session)
		v := session.Get(key)
		if v == nil {
			return nil
		} else {
			username := v.(string)
			return &username
		}
	}
	return nil
}
