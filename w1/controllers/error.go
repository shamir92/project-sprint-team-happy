package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"

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
	}

	// Domain Error
	pr, ok := err.(httpStatusCodeProvider)
	if ok {
		code = pr.HTTPStatusCode()
	}

	if code == http.StatusInternalServerError {
		errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

		trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

		errLog.Output(2, trace)

		message = "Internal Server Error"
	}

	c.JSON(code, gin.H{
		"message": message,
	})
}
