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

	app := fiber.New(fiber.Config{
		BodyLimit: 4 * 1024 * 1024,
	})

	app.Use(logger.New())
	app.Static("/bukti", "./bukti")
	
	app.Use(func(c *fiber.Ctx) error {
		ua := c.Get("User-Agent")
		if ua == "" || ua == "python-requests/2.31.0" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Bot detected",
			})
		}
		return c.Next()
	})

	// Rate limiter: max 10 request per 10 detik per IP
	app.Use(limiter.New(limiter.Config{
		Max:        10,
		Expiration: 10 * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			token := c.Get("Authorization") // per token if exists
			if token != "" {
				return c.IP() + ":" + token
			}
			return c.IP() // fallback anonymous
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, try later",
			})
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080,https://www.cryptolabsakademi.site,https://cryptolabsakademi.site",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    	AllowCredentials: true,
	}))

	app.Get("/hidden-admin-trap", func(c *fiber.Ctx) error {
		// log IP, could block temporarily
		fmt.Println("Bot hit honeypot:", c.IP())
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	})

	blockedIPs := map[string]bool{
		"1.2.3.4": true,
	}
	app.Use(func(c *fiber.Ctx) error {
		if blockedIPs[c.IP()] {
			return c.Status(fiber.StatusForbidden).SendString("Blocked IP")
		}
		return c.Next()
	})

	routes.Register(app)

	// Jalankan bot di background
	go service.StartBot()

	fmt.Println("Bot is running...")

	// Jalankan API
	if err := app.Listen(":3000"); err != nil {
		fmt.Println("Error starting Fiber:", err)
	}
}