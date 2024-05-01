package controllers

import (
	"gin-mvc/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createOrUpdateCatIn struct {
	Name        string   `json:"name" binding:"required,min=1,max=30"`
	Race        string   `json:"race" binding:"required,oneof=Persian MaineCoon Siamese Ragdoll Bengal Sphynx BritishShorthair Abyssinian ScottishFold Birman"`
	Sex         string   `json:"sex" binding:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" binding:"required,min=1,max=120082"`
	Description string   `json:"description" binding:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" binding:"required,min=1,dive,url"`
}

func CreateCat(c *gin.Context) {
	userId := c.GetString("userId")

	reqBody := createOrUpdateCatIn{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		handleError(c, err)
		return
	}

	cat, err := models.CreateCat(models.CreateOrUpdateCatIn{
		Name:        reqBody.Name,
		Race:        reqBody.Race,
		Age:         reqBody.AgeInMonth,
		Description: reqBody.Description,
		ImageURLs:   reqBody.ImageUrls,
	}, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success",
		"data": gin.H{
			"id":        cat.ID,
			"createdAt": cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), // ISO8601
		},
	})
}

func EditCatById(c *gin.Context) {
	catId := c.Param("catId")
	userId := c.GetString("userId")

	reqBody := createOrUpdateCatIn{}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		handleError(c, err)
		return
	}

	cat, err := models.EditCat(models.CreateOrUpdateCatIn{
		ID:          catId,
		Name:        reqBody.Name,
		Race:        reqBody.Race,
		Age:         reqBody.AgeInMonth,
		Description: reqBody.Description,
		ImageURLs:   reqBody.ImageUrls,
	}, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "cat updated successfully",
		"data": gin.H{
			"id": cat.ID,
		},
	})
}
