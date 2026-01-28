package routes

import (
	"github.com/gin-gonic/gin"
	"my-api/controllers/hello"
	"my-api/controllers/test"
	"my-api/controllers/user"
)

func InitRoutes(r *gin.Engine) {
	// 健康檢查
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	// GORM 方式的 API
	gorm := r.Group("/api/gorm")
	{
		gorm.GET("/users", user.GetUsersGORM)
		gorm.POST("/users", user.CreateUserGORM)
	}

	// 原生 SQL 方式的 API
	sql := r.Group("/api/sql")
	{
		sql.GET("/users", user.GetUsersSQL)
		sql.POST("/users", user.CreateUserSQL)
	}

	// 使用者相關路由
	u := r.Group("/users")
	{
		u.GET("/:name", user.GetByName)
	}

	// Hello 相關路由
	r.GET("/hello", hello.Handler)
	r.GET("/search", hello.Search)

	// 測試路由
	r.GET("/test", test.Handler)
}