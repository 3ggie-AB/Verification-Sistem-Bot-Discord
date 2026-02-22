package routes

import (
	"crypto-member/handlers"
	"crypto-member/middleware"

	"github.com/gofiber/fiber/v2"
)

func RoueteApi(api fiber.Router) {
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	
	api.Use(middleware.AuthRequired)
    api.Put("/update-profile", handlers.UpdateProfile)
    api.Get("/me", handlers.Me)
    api.Get("/notif", handlers.SSEHandler)

	api.Get("/membership/pricing", handlers.GetMembershipPricing)
	api.Post("/rule-pricing", handlers.CreateRulePricing)
	api.Get("/rule-pricing", handlers.GetAllRulePricing)
	api.Get("/rule-pricing/:id", handlers.GetRulePricingByID)
	api.Put("/rule-pricing/:id", handlers.UpdateRulePricing)
	api.Delete("/rule-pricing/:id", handlers.DeleteRulePricing)

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

	expense := api.Group("/expenses")
	expense.Post("/", handlers.CreateExpense)
	expense.Get("/", handlers.GetExpenses)
	expense.Get("/:id", handlers.GetExpenseByID)
	expense.Put("/:id", handlers.UpdateExpense)
	expense.Delete("/:id", handlers.DeleteExpense) 

	api.Post("/announcements", handlers.CreateAnnouncement)
	api.Get("/announcements", handlers.GetAllAnnouncements)

	progress := api.Group("/module-progress")
	progress.Post("/", handlers.UpdateModuleProgress)

	api.Post("/automessager", handlers.CreateAutoMessager)
	api.Get("/automessager", handlers.GetAutoMessagers) 
	api.Get("/automessager/:id", handlers.GetAutoMessagerByID)
	api.Put("/automessager/:id", handlers.UpdateAutoMessager)
	api.Patch("/automessager/:id/toggle", handlers.ToggleAutoMessager)
	api.Delete("/automessager/:id", handlers.DeleteAutoMessager)

	api.Post("/bots", handlers.CreateBot)
	api.Get("/bots", handlers.GetBots)
	api.Get("/bots/:id", handlers.GetBotByID)
	api.Put("/bots/:id", handlers.UpdateBot)
	api.Patch("/bots/:id/toggle", handlers.ToggleBot)
	api.Delete("/bots/:id", handlers.DeleteBot)
}
