package controllers

import (
	"gin-mvc/config"
	"gin-mvc/dtos"
	"gin-mvc/internal"
	"gin-mvc/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	jwtConfig := config.NewJWT()
	var db = internal.GetDB()
	var reqData dtos.LoginRequest
	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	err := db.Get(&user, "SELECT id,  email, password FROM users WHERE email=$1", reqData.Email)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "user not found",
			"data":    gin.H{},
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqData.Password))
	if err != nil {
		c.JSON(400, gin.H{
			"message": "email and password does not same",
			"data":    gin.H{},
		})
	}
	data := map[string]interface{}{
		"id":   user.ID,
		"emil": user.Email,
	}

	jwtToken, err := jwtConfig.GenerateJWT(data)
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(200, gin.H{
		"message": "User logged successfully",
		"data": gin.H{
			"email":       user.Email,
			"name":        user.FullName,
			"accessToken": jwtToken, // token should have 8 hours until expiration
		},
	})
}
