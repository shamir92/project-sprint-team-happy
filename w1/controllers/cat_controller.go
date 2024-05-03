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
	Race        string   `json:"race" binding:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Shorthair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
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
	db := internal.GetDB()
	catId := c.Param("catId")
	userId := c.GetString("userId")
	reqBody := createOrUpdateCatIn{}

	if err := c.BindJSON(&reqBody); err != nil {
		handleError(c, err)
		return
	}

	var count int
	db.Get(&count, "SELECT count(id) FROM matches WHERE issuer_cat_id=$1 or receiver_cat_id=$1", catId)
	if count > 0 {
		c.JSON(400, gin.H{
			"message": "sex is edited when cat is already requested to match",
			"data":    gin.H{},
		})
		return
	}

	cat, err := models.EditCat(models.CreateOrUpdateCatIn{
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
			"id": cat.ID,
		},
	})
}

func DeleteCatById(c *gin.Context) {

	db := internal.GetDB()
	catId := c.Param("catId")
	userId := c.GetString("userId")

	var count int
	db.Get(&count, "SELECT count(id) FROM matches WHERE issuer_cat_id=$1 or receiver_cat_id=$1", catId)
	if count > 0 {
		c.JSON(400, gin.H{
			"message": "cat cannot be deleted when cat is already requested to match",
			"data":    gin.H{},
		})
		return
	}

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
