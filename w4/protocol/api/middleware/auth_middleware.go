package middleware

import (
	"belimang/internal/helper"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// TODO: Add middleware for role based authorization

func AuthMiddleware(jwtManager helper.IJWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		token := strings.TrimPrefix(tokenString, "Bearer ")

		claim, err := jwtManager.GetClaim(token)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		c.Locals("user", claim)

		return c.Next()
	}
}
