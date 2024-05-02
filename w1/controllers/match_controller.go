package controllers

import (
	"gin-mvc/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type matchCreateIn struct {
	IssuerID      string
	IssuerCatID   string `json:"userCatId" binding:"required"`
	ReceiverID    string
	ReceiverCatID string `json:"matchCatId" binding:"required"`
	Message       string `json:"message" binding:"required,min=5,max=120"`
}

func MatchBrowse(c *gin.Context) {
	userId := c.GetString("userId")

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
	userId := c.GetString("userId")

	var req matchCreateIn
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err)
		return
	}

	if err := models.MatchCreate(userId, models.MatchCreateIn{
		IssuerCatID:   req.IssuerCatID,
		ReceiverCatID: req.ReceiverCatID,
		Message:       req.Message,
	}); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}
