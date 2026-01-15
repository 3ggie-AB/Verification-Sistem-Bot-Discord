package handlers

import (
	"crypto-member/db"
	"crypto-member/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func StreamModule(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	moduleID := c.Params("module_id")

	// 1. cek module
	var module models.Module
	if err := db.DB.
		Where("id = ? AND is_active = ?", moduleID, true).
		First(&module).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Module tidak ditemukan",
		})
	}

	// 2. cek membership
	now := time.Now()
	if user.MemberExpiredAt == nil || user.MemberExpiredAt.Before(now) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Membership tidak aktif",
		})
	}

	// 3. return embed url
	embedURL := "https://www.youtube.com/embed/" + module.YoutubeID + "?enablejsapi=1"

	return c.JSON(fiber.Map{
		"title":     module.Title,
		"embed_url": embedURL,
	})
}
