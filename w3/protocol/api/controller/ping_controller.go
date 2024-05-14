package controller

import (
	"halosuster/domain/usecase"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type pingController struct {
	pingUsecase usecase.IPingUsecase
}

// TODO: add all function under ping controller to inferface. this will make it easier to test
type IPingController interface {
	GetPingController(c *fiber.Ctx) error
}

func NewPingController(pingUsecase usecase.IPingUsecase) *pingController {
	return &pingController{
		pingUsecase: pingUsecase,
	}
}

func (pc *pingController) GetPingController(c *fiber.Ctx) error {
	value, _ := pc.pingUsecase.GetPing()
	if !value {
		return c.Status(http.StatusServiceUnavailable).JSON("pong")
	}
	return c.Status(http.StatusOK).JSON("pong")
}
