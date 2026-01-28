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

// TestHandler 示範變數和常數的用法
func TestHandler(c *gin.Context) {
	// === 常數 (const) ===
	// 常數在編譯時就確定，不能改變
	const AppName = "我的測試應用"
	const MaxRetry = 3
	const Pi = 3.14159
	
	// === 變數 (var) ===
	// 方式1: 完整聲明
	var message string = "這是測試的 Controller！"
	
	// 方式2: 自動推斷型別
	var count = 100
	
	// 方式3: 簡短聲明（最常用）
	userName := "測試使用者"
	isActive := true
	price := 99.99
	
	// 可以修改變數的值
	count = count + 1
	
	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"message":   message,
		"app_name":  AppName,      // 常數
		"max_retry": MaxRetry,     // 常數
		"count":     count,        // 變數
		"user_name": userName,     // 變數
		"is_active": isActive,     // 變數
		"price":     price,        // 變數
		"pi":        Pi,           // 常數
	})
}

// 示範：如果不需要使用 Gin Context 參數
// 這是一個普通函數（不是 HTTP Handler）
func CalculateSum() int {
	const tax = 0.05  // 常數：稅率
	var total = 0     // 變數
	
	price := 100
	quantity := 2
	
	total = price * quantity
	totalWithTax := float64(total) * (1 + tax)
	
	return int(totalWithTax)
}

// 但如果你要用在 Gin 路由中，還是需要有 *gin.Context 參數
// 可以用底線 _ 忽略不用的參數
func SimpleHandler(_ *gin.Context) {
	// 這裡不使用 c，用 _ 表示忽略
	// 但這樣就無法回傳 JSON 了
}