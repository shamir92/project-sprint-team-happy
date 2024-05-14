package controller

import (
	"github.com/gofiber/fiber/v2"
	// "github.com/gin-gonic/gin"
	// usecase "simple-invitation/domain/usecase"
	// "github.com/gin-gonic/gin"
)

type pingController struct {
}

// TODO: add all function under ping controller to inferface. this will make it easier to test
type IPingController interface {
	GetPingController(c *fiber.Ctx) error
}

func NewPingController() *pingController {
	return &pingController{}
}

func (ac *pingController) GetPingController(c *fiber.Ctx) error {
	return c.JSON("pong")
}
