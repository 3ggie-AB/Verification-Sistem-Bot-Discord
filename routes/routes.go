package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/", func (c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"AppName": "Backend CryptoLabs Akademi",
			"Pesan": "Selamat Datang di Backend CryptoLabs Akademi",
			"Status": "💚 Online",
		})
	})
	app.Get("/health", func (c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "oke",
		})
	})
	app.Head("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	
	api := app.Group("/api")
	webhook := app.Group("/webhook")

	RoueteApi(api)
	WebhookApi(webhook)
}
