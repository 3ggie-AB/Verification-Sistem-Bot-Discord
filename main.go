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
	
	"time"
    "github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	config.LoadEnv()
	db.Connect()

	app := fiber.New()
	app.Use(logger.New())
	app.Static("/bukti", "./bukti")
	
	// Rate limiter: max 10 request per 10 detik per IP
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later.",
			})
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080,https://www.cryptolabsakademi.site,https://cryptolabsakademi.site",
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