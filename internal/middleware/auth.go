package middleware

import (
	"chenze-faka/internal/pkg/jwt"
	"chenze-faka/internal/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UsernameKey contextKey = "username"
	IsAdminKey  contextKey = "is_admin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			// fallback to cookie
			cookie, err := c.Cookie("token")
			if err == nil && cookie != "" {
				auth = "Bearer " + cookie
			}
		}
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": response.CodeTokenInvalid, "msg": "未登录"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwt.Parse(tokenStr, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": response.CodeTokenInvalid, "msg": "登录已失效"})
			c.Abort()
			return
		}
		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(UsernameKey), claims.Username)
		c.Set(string(IsAdminKey), claims.IsAdmin)
		c.Next()
	}
}

func AdminAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			cookie, err := c.Cookie("admin_token")
			if err == nil && cookie != "" {
				auth = "Bearer " + cookie
			}
		}
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwt.Parse(tokenStr, secret)
		if err != nil || !claims.IsAdmin {
			c.Redirect(http.StatusFound, "/admin/login")
			c.Abort()
			return
		}
		c.Set(string(UserIDKey), claims.UserID)
		c.Set(string(UsernameKey), claims.Username)
		c.Set(string(IsAdminKey), claims.IsAdmin)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
