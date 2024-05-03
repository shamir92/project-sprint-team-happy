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

type matchAnswerIn struct {
	MatchID string `json:"matchId" binding:"required"`
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
	if err := c.BindJSON(&req); err != nil {
		// handleError(c, err)
		return
	}

	statusCode, err := models.MatchCreate(userId, models.MatchCreateIn{
		IssuerCatID:   req.IssuerCatID,
		ReceiverCatID: req.ReceiverCatID,
		Message:       req.Message,
	})
	if err != nil {
		c.JSON(statusCode, gin.H{
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func MatchAnswerReject(c *gin.Context) {
	userId := c.GetString("userId")

	var req matchAnswerIn
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err)
		return
	}

	if err := models.MatchReject(userId, models.MatchAnswerIn{MatchID: req.MatchID}); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func MatchAnswerApprove(c *gin.Context) {
	userId := c.GetString("userId")

	var req matchAnswerIn
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, err)
		return
	}

	if err := models.MatchApprove(userId, models.MatchAnswerIn{MatchID: req.MatchID}); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func MatchDelete(c *gin.Context) {
	userId := c.GetString("userId")
	matchId := c.Param("matchId")

	err := models.DeleteMatch(matchId, userId)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
