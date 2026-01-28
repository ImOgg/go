package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - 驗證中間件
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

		token := parts[1]

		// TODO: 實作真正的 Token 驗證（JWT、Session 等）
		// 這裡先做簡單驗證示範
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "無效的授權憑證",
			})
			c.Abort()
			return
		}

		// 驗證成功，可以將用戶資訊存入 context
		// user, err := validateToken(token)
		// if err != nil {
		//     c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "授權已過期"})
		//     c.Abort()
		//     return
		// }
		// c.Set("user", user)

		// 暫時示範：任何有 token 的都通過
		c.Set("token", token)
		
		c.Next()
	}
}
