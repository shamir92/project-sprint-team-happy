package controller

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type errHttpStatusCodeProvider interface {
	HTTPStatusCode() int
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
	Code    int    `json:"-"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	errResp := ErrorResponse{
		Message: fiber.ErrInternalServerError.Message,
		Error:   err.Error(),
	}

	switch e := err.(type) {
	case errHttpStatusCodeProvider:
		errResp.Code = e.HTTPStatusCode()
		errResp.Message = err.Error()
	case *fiber.Error:
		errResp.Code = e.Code
		errResp.Message = e.Message
	}

	// go-validator erros
	if validationErrors, ok := err.(validator.ValidationErrors); ok && len(validationErrors) > 0 {
		err := validationErrors[0]

		errResp.Code = fiber.StatusBadRequest
		errResp.Message = fiber.ErrBadRequest.Message
		errResp.Error = fmt.Sprintf("'%s' failed on %s validation", err.Field(), err.ActualTag())
	}

	return c.Status(errResp.Code).JSON(errResp)
}
