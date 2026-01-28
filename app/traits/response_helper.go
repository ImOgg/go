package traits

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RespondJSON - 統一的 JSON 回應輔助函式
func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// RespondSuccess - 成功回應
func RespondSuccess(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// RespondCreated - 建立成功回應
func RespondCreated(c *gin.Context, data interface{}, message string) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// RespondError - 錯誤回應
func RespondError(c *gin.Context, status int, message string, errors interface{}) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"errors":  errors,
	})
}

// RespondValidationError - 驗證錯誤回應
func RespondValidationError(c *gin.Context, errors interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "驗證失敗",
		"errors":  errors,
	})
}

// RespondNotFound - 找不到資源回應
func RespondNotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": message,
	})
}

// RespondUnauthorized - 未授權回應
func RespondUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": message,
	})
}

// RespondForbidden - 禁止存取回應
func RespondForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": message,
	})
}
