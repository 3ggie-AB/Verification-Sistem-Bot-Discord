package main

import (
	"fmt"
	"crypto-member/config"
	"crypto-member/service"
	"crypto-member/db"
	"crypto-member/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadEnv()
	db.Connect()
	app := fiber.New()
	app.Use(logger.New())
	routes.Register(app)
	service.StartBot()

	fmt.Println("Bot is running...")
	select {}

	// go func() {
	// 	for {
	// 		// service.ExpireMemberships()
	// 		time.Sleep(1 * time.Hour)
	// 	}
	// }()

	// go bot.Start() // bot jalan barengan API

	app.Listen(":3000")
}
