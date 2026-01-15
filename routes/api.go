package routes

import (
	"crypto-member/handlers"
	"crypto-member/middleware"

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

	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	
	api.Use(middleware.AuthRequired)
    api.Put("/update-profile", handlers.UpdateProfile)
    api.Get("/me", handlers.Me)

	api.Post("/checkout", handlers.CheckoutMembership)

	api.Get("/payments", handlers.GetPayments)
	api.Post("/payments/:id/approve", handlers.ApprovePayment)
	api.Post("/payments/:id/reject", handlers.RejectPayment)
	api.Delete("/payments/:id", handlers.DeletePayment)

	coupon := api.Group("/coupons")
	coupon.Get("/check", handlers.CheckCouponByCode)
	coupon.Post("/", handlers.CreateCoupon)
	coupon.Get("/", handlers.GetCoupons)
	coupon.Get("/:id", handlers.GetCouponByID)
	coupon.Put("/:id", handlers.UpdateCoupon)
	coupon.Delete("/:id", handlers.DeleteCoupon)

	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Post("/", handlers.CreateUser)
	users.Put("/:id", handlers.UpdateUser)
	users.Delete("/:id", handlers.DeleteUser)

	api.Get("/stream/:module_id", handlers.StreamModule)

	moduleGroup := api.Group("/module-groups")
	moduleGroup.Get("/", handlers.GetModuleGroups)
	moduleGroup.Post("/", handlers.CreateModuleGroup)
	moduleGroup.Put("/:id", handlers.UpdateModuleGroup)
	moduleGroup.Delete("/:id", handlers.DeleteModuleGroup)

	module := api.Group("/modules")
	module.Get("/group/:group_id", handlers.GetModulesByGroup)
	module.Post("/", handlers.CreateModule)
	module.Put("/:id", handlers.UpdateModule)
	module.Delete("/:id", handlers.DeleteModule)

	progress := api.Group("/module-progress")
	progress.Post("/", handlers.UpdateModuleProgress)
}
