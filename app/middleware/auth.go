package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"my-api/app/utils"
)

// AuthMiddleware - JWT 驗證中間件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Header 取得 Token
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "缺少授權憑證",
			})
			c.Abort()
			return
		}

		// 檢查 Token 格式：Bearer <token>
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "無效的授權格式",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 驗證 JWT Token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "無效或已過期的授權憑證",
			})
			c.Abort()
			return
		}

		// 將用戶資訊存入 context，供後續 handler 使用
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}
