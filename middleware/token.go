package middleware

import (
	"crypto-member/db"
	"crypto-member/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthRequired(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing token"})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token format"})
	}

	token := parts[1]

	var user models.User
	if err := db.DB.Where("token = ?", token).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token"})
	}

	c.Locals("user", &user) // simpan user di context
	return c.Next()
}
