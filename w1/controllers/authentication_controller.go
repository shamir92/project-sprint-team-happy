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

func generateJwtToken(user models.User) (string, error) {
	jwtConfig := config.NewJWT()

	data := map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	}

	jwtToken, err := jwtConfig.GenerateJWT(data)

	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func Login(c *gin.Context) {
	jwtConfig := config.NewJWT()
	var db = internal.GetDB()
	var reqData dtos.LoginRequest
	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{}
	err := db.Get(&user, "SELECT id,  email, fullname, password FROM users WHERE email=$1", reqData.Email)
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
		"id":    user.ID,
		"email": user.Email,
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

func Register(c *gin.Context) {
	// TODO: adjust based on project requirement
	reqBody := struct {
		Email    string `json:"email" binding:"required,email"`
		Name     string `json:"name" binding:"required,min=5,max=50"`
		Password string `json:"password" binding:"required,min=5,max=15"`
	}{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		handleError(c, err)
		return
	}

	createdUser, err := models.CreateUser(models.RegisterUser{
		Email:    reqBody.Email,
		Password: reqBody.Password,
		FullName: reqBody.Name,
	})

	// TODO: make an error handler
	if err != nil {
		handleError(c, err)
		return
	}

	accessToken, err := generateJwtToken(createdUser)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": gin.H{
			"email":       createdUser.Email,
			"name":        createdUser.FullName,
			"accessToken": accessToken,
		},
	})
}
