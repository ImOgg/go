package test

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handler - 示範變數和常數的用法
func Handler(c *gin.Context) {
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

// CalculateSum - 示範：不需要 Gin Context 的普通函數
func CalculateSum() int {
	const tax = 0.05  // 常數：稅率
	var total = 0     // 變數
	
	price := 100
	quantity := 2
	
	total = price * quantity
	totalWithTax := float64(total) * (1 + tax)
	
	return int(totalWithTax)
}
