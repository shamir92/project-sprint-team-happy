package controllers

import (
	"gin-mvc/internal"
	"gin-mvc/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type createOrUpdateCatIn struct {
	Name        string   `json:"name" binding:"required,min=1,max=30"`
	Race        string   `json:"race" binding:"required"`
	Sex         string   `json:"sex" binding:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" binding:"required,min=1,max=120082"`
	Description string   `json:"description" binding:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" binding:"required,dive,min=1,url"`
}

func CreateCat(c *gin.Context) {
	userId := c.GetString("userId")

	reqBody := createOrUpdateCatIn{}

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	cat, err := models.CreateCat(models.CreateOrUpdateCatIn{
		Name:        reqBody.Name,
		Race:        reqBody.Race,
		Age:         reqBody.AgeInMonth,
		Description: reqBody.Description,
		ImageURLs:   reqBody.ImageUrls,
		Sex:         reqBody.Sex,
	}, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data": gin.H{
			"id":        cat.ID,
			"createdAt": internal.ToISO8601Format(cat.CreatedAt), // ISO8601
		},
	})
}

func EditCatById(c *gin.Context) {
	catId := c.Param("catId")
	userId := c.GetString("userId")
	reqBody := createOrUpdateCatIn{}

	if err := c.BindJSON(&reqBody); err != nil {
		return
	}

	err := models.EditCat(models.CreateOrUpdateCatIn{
		ID:          catId,
		Name:        reqBody.Name,
		Race:        reqBody.Race,
		Age:         reqBody.AgeInMonth,
		Description: reqBody.Description,
		ImageURLs:   reqBody.ImageUrls,
		Sex:         reqBody.Sex,
	}, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cat updated successfully",
		"data": gin.H{
			"id": catId,
		},
	})
}

func DeleteCatById(c *gin.Context) {
	catId := c.Param("catId")
	userId := c.GetString("userId")

	err := models.DeleteCatById(catId, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func GetCats(c *gin.Context) {
	var query models.GetCatOption
	if err := c.ShouldBindWith(&query, binding.Query); err != nil {
		handleError(c, err)
		return
	}

	userId := c.GetString("userId")

	cats, err := models.GetCats(query, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get cats",
		"data":    cats,
	})
}
