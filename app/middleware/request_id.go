package middleware

import (
	"my-api/app/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID 中間件 - 為每個請求生成唯一 ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 嘗試從 header 取得（支援分散式追蹤）
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 設定到 context
		c.Set(logger.ContextKeyRequestID, requestID)

		// 設定 response header
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}
