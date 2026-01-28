package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"my-api/app"
	"my-api/bootstrap"
	"my-api/config"
	"my-api/database"
	_ "my-api/database/migrations" // 引入 migrations 確保註冊
	"my-api/routes"
)

func main() {
	// 載入配置
	config.LoadConfig()

	// 初始化資料庫連接
	bootstrap.InitDB()

	// 自動執行 migrations
	if err := database.RunMigrations(); err != nil {
		log.Println("⚠️  Migration 警告:", err)
		// 不中斷程式，繼續啟動
	}

	// 初始化 Redis（可選）
	// bootstrap.InitRedis()

	// 建立應用程式容器（Laravel 風格）
	application := app.NewApp(bootstrap.DB)

	r := gin.Default()

	// 設定所有路由
	routes.SetupRoutes(r, application)

	r.Run(":" + config.GlobalConfig.App.Port)
}