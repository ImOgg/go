package main

import (
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

	// 初始化 Logger（必須在其他初始化之前）
	bootstrap.InitLogger()

	// 初始化資料庫連接
	bootstrap.InitDB()

	// 自動執行 migrations
	if err := database.RunMigrations(); err != nil {
		bootstrap.Log.Warning("Migration 警告", map[string]interface{}{
			"error": err.Error(),
		})
		// 不中斷程式，繼續啟動
	}

	// 初始化 Redis（可選）
	// bootstrap.InitRedis()

	// 建立應用程式容器（Laravel 風格）
	application := app.NewApp(bootstrap.DB, bootstrap.Log)

	// 設定 Gin 模式
	if config.GlobalConfig.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // 使用 gin.New() 而非 gin.Default()，避免重複的日誌

	// 設定所有路由
	routes.SetupRoutes(r, application)

	bootstrap.Log.Info("Server starting", map[string]interface{}{
		"port": config.GlobalConfig.App.Port,
		"env":  config.GlobalConfig.App.Env,
	})

	r.Run(":" + config.GlobalConfig.App.Port)
}