package logger

import (
	"context"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyLogger 是存放在 context 中的 Logger key
	ContextKeyLogger = "logger"
	// ContextKeyRequestID 是存放 Request ID 的 key
	ContextKeyRequestID = "request_id"
)

// FromContext 從 context 中取得 Logger
func FromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(ContextKeyLogger).(*Logger); ok {
		return l
	}
	return Global()
}

// FromGinContext 從 Gin context 中取得 Logger
func FromGinContext(c *gin.Context) *Logger {
	if l, exists := c.Get(ContextKeyLogger); exists {
		if logger, ok := l.(*Logger); ok {
			return logger
		}
	}
	return Global()
}

// ToGinContext 將 Logger 存入 Gin context
func ToGinContext(c *gin.Context, l *Logger) {
	c.Set(ContextKeyLogger, l)
}

// GetRequestID 從 Gin context 取得 Request ID
func GetRequestID(c *gin.Context) string {
	if id, exists := c.Get(ContextKeyRequestID); exists {
		return id.(string)
	}
	return ""
}
