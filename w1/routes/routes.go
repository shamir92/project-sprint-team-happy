package routes

import (
	"gin-mvc/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/", controllers.Index)
	r.GET("/users", controllers.GetUsers)
}
