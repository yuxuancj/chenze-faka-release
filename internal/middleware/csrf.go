package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			c.Next()
			return
		}

		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("_csrf")
		}
		sessionToken, exists := c.Get("csrf_token")
		if !exists || token != sessionToken {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "CSRF token invalid",
				"data": nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func CSRFTokenSetter() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := generateCSRFToken()
		c.Set("csrf_token", token)
		c.SetCookie("csrf_token", token, 3600, "/", "", false, true)
		c.Next()
	}
}
