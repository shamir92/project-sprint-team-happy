package controllers

import (
	"fmt"
	"gin-mvc/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MatchBrowse(c *gin.Context) {
	// userId := c.GetString("userId")
	userId := "a0ed9574-c63d-4063-bc04-c2945120a67c"

	matches, err := models.MatchAll(userId)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    matches,
	})
}

func MatchCreate(c *gin.Context) {
	userId := "a0ed9574-c63d-4063-bc04-c2945120a67c"
	// userId := c.GetString("userId")

	var req models.MatchSubmit
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err)
		return
	}

	issuerCat, err := models.CatGetById(req.IssuerCatID, nil)
	if err != nil {
		handleError(c, err)
		return
	}

	if issuerCat.OwnerID.String() != userId {
		handleError(c, models.MatchError{Message: models.ErrCantMatchOtherCats.Error(), StatusCode: http.StatusNotFound})
		return
	}

	receiverCat, err := models.CatGetById(req.ReceiverCatID, nil)
	if err != nil {
		handleError(c, err)
		return
	}

	if receiverCat.OwnerID.String() == userId {
		handleError(c, models.MatchError{Message: models.ErrCantMatchYourOwnCats.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	if receiverCat.HasMatched || issuerCat.HasMatched {
		handleError(c, models.MatchError{Message: models.ErrCatAlreadyMatched.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	if receiverCat.Sex == issuerCat.Sex {
		handleError(c, models.MatchError{Message: models.ErrSameGender.Error(), StatusCode: http.StatusBadRequest})
		return
	}

	req.IssuerID = userId
	req.ReceiverID = receiverCat.OwnerID.String()

	fmt.Println(req)
	err = models.MatchCreate(userId, req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}
