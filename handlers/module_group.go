package handlers

import (
	"github.com/gofiber/fiber/v2"

	"crypto-member/db"
	"crypto-member/models"
)

func GetModuleGroups(c *fiber.Ctx) error {
	var groups []models.ModuleGroup

	if err := db.DB.
		// Preload("Modules").
		Order("created_at DESC").
		Find(&groups).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal ambil module group",
		})
	}

	return c.JSON(fiber.Map{
		"data": groups,
	})
}

func CreateModuleGroup(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	if user.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}

	var body struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		IsActive    bool    `json:"is_active"`
	}

	if err := c.BodyParser(&body); err != nil || body.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	group := models.ModuleGroup{
		Title:       body.Title,
		Description: body.Description,
		IsActive:    body.IsActive,
	}

	if err := db.DB.Create(&group).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal buat module group"})
	}

	return c.JSON(fiber.Map{
		"message": "Module group berhasil dibuat",
		"data":    group,
	})
}

func UpdateModuleGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}

	var group models.ModuleGroup
	if err := db.DB.First(&group, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Module group tidak ditemukan"})
	}

	var body struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		IsActive    bool    `json:"is_active"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	group.Title = body.Title
	group.Description = body.Description
	group.IsActive = body.IsActive

	db.DB.Save(&group)

	return c.JSON(fiber.Map{
		"message": "Module group berhasil diupdate",
	})
}

func DeleteModuleGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}

	if err := db.DB.Delete(&models.ModuleGroup{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus module group"})
	}

	return c.JSON(fiber.Map{
		"message": "Module group berhasil dihapus",
	})
}
