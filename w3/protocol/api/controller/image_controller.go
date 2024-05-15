package controller

import (
	"halosuster/domain/usecase"
	"halosuster/protocol/api/dto"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type imageController struct {
	imageUsecase usecase.IImageUsecase
}

// TODO: add all function under ping controller to inferface. this will make it easier to test
type IImageController interface {
	UploadImage(c *fiber.Ctx) error
}

func NewImageController(imageUsecase usecase.IImageUsecase) *imageController {
	return &imageController{
		imageUsecase: imageUsecase,
	}
}

type UploadImageResponse struct {
	URL string `json:"imageUrl"`
}

func (ic *imageController) UploadImage(c *fiber.Ctx) error {

	// Get the file from the form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Failed to get file from form")
	}

	url, err := ic.imageUsecase.UploadImage(fileHeader)
	if err != nil {
		// tar ganti
		// dirty dulu.
		return c.Status(http.StatusInternalServerError).JSON(err)
	}

	// TODO: return response
	return c.Status(http.StatusCreated).JSON(dto.UserITRegisterControllerResponse{
		Message: "user registered successfully",
		Data:    UploadImageResponse{URL: url},
	})

}
