package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"my-api/app/pkg/logger"
	"my-api/bootstrap"
)

// Logger - 結構化 HTTP 請求日誌中間件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 取得 Request ID
		requestID := logger.GetRequestID(c)

		// 建立帶有 request_id 的子 Logger
		reqLogger := bootstrap.Log.WithRequestID(requestID)

		// 將 Logger 存入 context，供後續使用
		logger.ToGinContext(c, reqLogger)

		// 記錄請求開始時間
		startTime := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 處理請求
		c.Next()

		// 計算執行時間
		latency := time.Since(startTime)

		// 取得狀態碼
		statusCode := c.Writer.Status()

		// 建立日誌 context
		logContext := map[string]interface{}{
			"status":     statusCode,
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
			"body_size":  c.Writer.Size(),
		}

		// 如果有錯誤，加入錯誤訊息
		if len(c.Errors) > 0 {
			logContext["errors"] = c.Errors.String()
		}

		// 根據狀態碼選擇日誌等級
		switch {
		case statusCode >= 500:
			reqLogger.Error("Server Error", logContext)
		case statusCode >= 400:
			reqLogger.Warning("Client Error", logContext)
		case statusCode >= 300:
			reqLogger.Info("Redirect", logContext)
		default:
			reqLogger.Info("Request", logContext)
		}
	}
}
