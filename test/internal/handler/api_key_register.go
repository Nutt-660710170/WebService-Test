package handler

import (
	"net/http"
	"time"

	"oil/pkg/database"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterAPIKeyHandler handles API key registration requests
func RegisterAPIKeyHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type reqBody struct {
			Owner     string     `json:"owner"`
			ExpiredAt *time.Time `json:"expired_at"`
		}
		var body reqBody
		if err := c.BodyParser(&body); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		apiKey, err := database.SaveAPIKey(db, body.Owner, body.ExpiredAt)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate API key"})
		}
		return c.JSON(fiber.Map{
			"api_key": apiKey.Key,
			"expired_at": apiKey.ExpiredAt,
		})
	}
}
