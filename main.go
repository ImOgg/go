package main

import (
	"github.com/gin-gonic/gin"
	"my-api/routes"
)

func main() {
	r := gin.Default()

	// 呼叫路由設定
	routes.InitRoutes(r)

	r.Run(":8080")
}