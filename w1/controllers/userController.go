package controllers

import (
	"gin-mvc/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Home Page",
	})
}

func GetUsers(c *gin.Context) {
	users := models.GetAllUsers() // Assume this method fetches users from the database
	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
