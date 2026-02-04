package routes

import (
	"github.com/gofiber/fiber/v2"
)

func WebhookApi(webhook fiber.Router) {
	webhook.Post("/test", func(c *fiber.Ctx) error {
		body := c.Body()
		return c.JSON(fiber.Map{
			"message": "Webhook test successful",
			"data":    string(body),
		})
	})
}
