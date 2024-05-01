package routes

import (
	"gin-mvc/controllers"
	"gin-mvc/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", controllers.Ping)
	v1 := r.Group("/v1")
	{
		// public route without auth middleware
		v1.POST("/user/login", controllers.Login)
		v1.POST("/user/register", controllers.Register)
		// private route with auth middleware
		v1.GET("/ping-private", middlewares.AuthMiddleware(), controllers.Ping)
	}
}
