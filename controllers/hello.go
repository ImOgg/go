package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 注意：Function 開頭要大寫，外部才抓得到 (Public)
func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "這就是從 Controller 噴出來的資料！",
	})
}

func GetUserByName(c *gin.Context) {
	// 抓取網址路徑裡的 :name
	name := c.Param("name")

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "歡迎回來！",
		"data":    name,
	})
}

// 2. 進階：Query 參數 (網址後面的 ?)
// 如果你想處理像 http://localhost:8080/search?keyword=golang 這種參數，寫法又不太一樣：
func Search(c *gin.Context) {
	// 抓取 ?keyword=xxx 裡的值
	keyword := c.DefaultQuery("keyword", "預設值")

	c.JSON(http.StatusOK, gin.H{
		"search_for": keyword,
	})
}