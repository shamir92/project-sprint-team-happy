package controller

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"halosuster/domain/entity"
	"halosuster/domain/usecase"
	"halosuster/internal/helper"
	"halosuster/protocol/api/dto"

	"github.com/gofiber/fiber/v2"
)

type userITController struct {
	// pingUsecase usecase.IPingUsecase
	userITUsecase usecase.IUserITUsecase
}

// TODO: add all function under ping controller to inferface. this will make it easier to test
type IUserITController interface {
	RegisterUserIT(c *fiber.Ctx) error
	LoginUserIT(c *fiber.Ctx) error
	GetListUsers(c *fiber.Ctx) error
}

func NewUserITController(userITUsecase usecase.IUserITUsecase) *userITController {
	return &userITController{
		userITUsecase: userITUsecase,
	}
}

func (pc *userITController) RegisterUserIT(c *fiber.Ctx) error {
	//TODO: initization variable
	var request usecase.UserITRegisterRequest

	// parse json
	if err := c.BodyParser(&request); err != nil {
		return fiber.ErrBadRequest
	}

	// TODO: Validate Struct
	if err := helper.ValidateStruct(&request); err != nil {
		return err
	}

	data, err := pc.userITUsecase.RegisterUserIT(request)
	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(dto.UserITRegisterControllerResponse{
		Message: "user registered successfully",
		Data:    data,
	})
}

func (pc *userITController) LoginUserIT(c *fiber.Ctx) error {
	//TODO: initization variable
	var request usecase.UserITLoginRequest

	// parse json
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// TODO: Validate Struct
	if err := helper.ValidateStruct(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// TODO: add logic
	data, err := pc.userITUsecase.LoginUserIT(request)
	if err != nil {
		return err
	}
	// TODO: return response

	return c.Status(http.StatusOK).JSON(dto.UserITRegisterControllerResponse{
		Message: "user login successfully",
		Data:    data,
	})
}

func (pc *userITController) GetListUsers(c *fiber.Ctx) error {
	var req entity.ListUserPayload

	err := c.QueryParser(&req)

	if err != nil {
		// TODO: Adjust based on k6 test cases
		log.Printf("GetListUsers: %v\n", err)
	}

	users, err := pc.userITUsecase.GetUsers(req)

	if err != nil {
		log.Printf("gailed to get list users: %v\n", err)
		return err
	}

	var listUsers = []dto.ListUserItemDto{}

	for _, u := range users {
		integer, err := strconv.Atoi(strings.TrimSpace(u.NIP))
		if err != nil {
			log.Println(err)
		}

		listUsers = append(listUsers, dto.ListUserItemDto{
			ID:        u.ID.String(),
			Name:      u.Name,
			NIP:       integer,
			CreatedAt: u.CreatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "success",
		"data":    listUsers,
	})
}
