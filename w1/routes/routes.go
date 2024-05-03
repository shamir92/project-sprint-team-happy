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

		// cats
		v1.GET("/cat", middlewares.AuthMiddleware(), controllers.GetCats)
		v1.POST("/cat", middlewares.AuthMiddleware(), controllers.CreateCat)
		v1.PUT("/cat/:catId", middlewares.AuthMiddleware(), controllers.EditCatById)
		v1.DELETE("/cat/:catId", middlewares.AuthMiddleware(), controllers.DeleteCatById)

		match := v1.Group("/cat/match")
		match.GET("", middlewares.AuthMiddleware(), controllers.MatchBrowse)
		match.POST("", middlewares.AuthMiddleware(), controllers.MatchCreate)
		match.POST("/approve", middlewares.AuthMiddleware(), controllers.MatchAnswerApprove)
		match.POST("/reject", middlewares.AuthMiddleware(), controllers.MatchAnswerReject)
		match.DELETE("/:matchId", middlewares.AuthMiddleware(), controllers.MatchDelete)

	}
}
