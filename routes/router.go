package routes

import (
    "github.com/gin-gonic/gin"
    "my-api/controllers"
)

func InitRoutes(r *gin.Engine) {
    // 把剛才那一坨 Group 搬到這裡
    u := r.Group("/users")
    {
        u.GET("/:name", controllers.GetUserByName)
    }
    
    // 你可以在這裡一直加...
    r.GET("/health", func(c *gin.Context) {
        c.String(200, "OK")
    })
}