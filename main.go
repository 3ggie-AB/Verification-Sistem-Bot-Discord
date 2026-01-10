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

	// Jalankan bot di background
	go service.StartBot()

	fmt.Println("Bot is running...")

	// Jalankan API
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting Fiber:", err)
	}
}