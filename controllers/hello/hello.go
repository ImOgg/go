package hello

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler - 基本的 Hello Handler
func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "這就是從 Controller 噴出來的資料！",
	})
}

// Search - Query 參數範例：處理 ?keyword=xxx
func Search(c *gin.Context) {
	// 抓取 ?keyword=xxx 裡的值
	keyword := c.DefaultQuery("keyword", "預設值")

	c.JSON(http.StatusOK, gin.H{
		"search_for": keyword,
	})
}
