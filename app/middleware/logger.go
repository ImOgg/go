package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger - 自訂日誌中間件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 記錄請求開始時間
		startTime := time.Now()

		// 處理請求
		c.Next()

		// 計算執行時間
		latency := time.Since(startTime)

		// 取得狀態碼
		statusCode := c.Writer.Status()

		// 輸出日誌
		log.Printf("[%s] %s %s %d %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			statusCode,
			latency,
		)
	}
}
