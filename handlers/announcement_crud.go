package handlers

import (
	"crypto-member/models"
	"crypto-member/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CreateAnnouncementRequest struct {
	Title    string          `json:"title"`
	Content  string          `json:"content"`
	Type     string          `json:"type"`
	Channels []string        `json:"channels"`
	Target   datatypes.JSON  `json:"target"`
}

func CreateAnnouncement(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}
	var req CreateAnnouncementRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	announcement := models.Announcement{
		ID:      uuid.NewString(),
		Title:   req.Title,
		Content: req.Content,
		Type:    req.Type,
		Target:  req.Target,
	}

	if err := service.CreateAnnouncement(announcement, req.Channels); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal membuat announcement",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Announcement berhasil dibuat",
	})
}

func GetAllAnnouncements(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}
	data, err := service.GetAllAnnouncements()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mengambil data",
		})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}