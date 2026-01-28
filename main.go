package main

import (
	"github.com/gin-gonic/gin"
	"my-api/config"
	"my-api/database"
	"my-api/routes"
)

func main() {
	// 載入配置
	config.LoadConfig()

	// 方式一：使用 GORM（推薦用於快速開發）
	database.InitDB()

	// 方式二：使用原生 SQL（可選，適合複雜查詢）
	// database.InitRawDB()
	// defer database.CloseRawDB()

	// 初始化 Redis（可選）
	// database.InitRedis()

	r := gin.Default()

	// 呼叫路由設定
	routes.InitRoutes(r)

	r.Run(":" + config.GlobalConfig.App.Port)
}