package routes

import (
	"github.com/gin-gonic/gin"
	"my-api/app"
	"my-api/app/controllers"
	"my-api/app/middleware"
)

// SetupRoutes - 設定所有路由（Laravel 風格）
func SetupRoutes(router *gin.Engine, application *app.App) {
	// 建立 Controllers
	userCtrl := controllers.NewUserController(application)

	// 健康檢查（不需要驗證）
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "API is running",
		})
	})

	// API 路由群組
	api := router.Group("/api")
	{
		// 公開路由（不需要驗證）
		public := api.Group("")
		{
			// 可以加登入、註冊等公開 API
			// public.POST("/login", authCtrl.Login)
			// public.POST("/register", authCtrl.Register)
		}

		// 需要驗證的路由
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware()) // 使用驗證中間件
		{
			// RESTful User 路由
			users := protected.Group("/users")
			{
				users.GET("", userCtrl.Index)          // GET    /api/users
				users.POST("", userCtrl.Store)         // POST   /api/users
				users.GET("/:id", userCtrl.Show)       // GET    /api/users/:id
				users.PUT("/:id", userCtrl.Update)     // PUT    /api/users/:id
				users.PATCH("/:id", userCtrl.Update)   // PATCH  /api/users/:id
				users.DELETE("/:id", userCtrl.Destroy) // DELETE /api/users/:id
			}

			// 其他需要驗證的路由
			// protected.GET("/profile", userCtrl.Profile)
		}
	}
}
