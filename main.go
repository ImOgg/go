package main

import (
	"log"
	
	"github.com/gin-gonic/gin"
	"my-api/app"
	"my-api/config"
	"my-api/database"
	_ "my-api/database/migrations" // 引入 migrations 確保註冊
	"my-api/routes"
)

func main() {
	// 載入配置
	config.LoadConfig()

	// 初始化資料庫連接
	database.InitDB()
	
	// 自動執行 migrations
	if err := database.RunMigrations(); err != nil {
		log.Println("⚠️  Migration 警告:", err)
		// 不中斷程式，繼續啟動
	}

	// 方式二：使用原生 SQL（可選，適合複雜查詢）
	// database.InitRawDB()
	// defer database.CloseRawDB()

	// 初始化 Redis（可選）
	// database.InitRedis()

	// 建立應用程式容器（Laravel 風格）
	application := app.NewApp(database.DB)

	r := gin.Default()

	// 設定所有路由
	routes.SetupRoutes(r, application)

	r.Run(":" + config.GlobalConfig.App.Port)
}