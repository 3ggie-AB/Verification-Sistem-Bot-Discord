package main

import (
	"fmt"
	"crypto-member/config"
	"crypto-member/service"
	"crypto-member/db"
	"crypto-member/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.LoadEnv()
	db.Connect()

	app := fiber.New()
	app.Use(logger.New())
	app.Static("/bukti", "./bukti")
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080, https://www.cryptolabsakademi.site",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    	AllowCredentials: true,
	}))

	routes.Register(app)

	// Jalankan bot di background
	go service.StartBot()

	fmt.Println("Bot is running...")

	// Jalankan API
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting Fiber:", err)
	}
}