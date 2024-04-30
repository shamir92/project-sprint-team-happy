package routes

import (
	"gin-mvc/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", controllers.Ping)
	v1 := r.Group("/v1")
	{
		v1.POST("/login", controllers.Login)

	}

}
