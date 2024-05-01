package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type httpStatusCodeProvider interface {
	HTTPStatusCode() int
}

func handleError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	message := err.Error()

	// go-validator
	if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
		code = http.StatusBadRequest
		message = err.Error()
	}

	// Domain Error
	pr, ok := err.(httpStatusCodeProvider)
	if ok {
		code = pr.HTTPStatusCode()
		message = err.Error()
	}

	c.JSON(code, gin.H{
		"message": message,
	})
}
