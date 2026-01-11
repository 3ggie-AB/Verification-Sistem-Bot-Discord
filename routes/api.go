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

	coupon := api.Group("/coupons")
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
}
