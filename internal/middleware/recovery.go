package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(map[string]interface{}{
					"error": err,
					"time":  time.Now().Format(time.RFC3339),
					"path":  c.Request.URL.Path,
				}).Error("Panic recovered")

				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "服务器内部错误",
					"error":   err,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
