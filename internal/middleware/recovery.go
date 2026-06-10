package middleware

import (
	"chenze-faka/internal/pkg/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %v\n%s", err, string(debug.Stack()))
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  "服务器内部错误",
					"data": nil,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
