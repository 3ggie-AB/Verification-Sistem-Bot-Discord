package handlers

import (
	"github.com/gofiber/fiber/v2"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/google/uuid"
)

func GetModulesByGroup(c *fiber.Ctx) error {
	groupID := c.Params("group_id")

	var modules []models.Module
	if err := db.DB.
		Where("module_group_id = ? AND is_active = ?", groupID, true).
		Order("published_at ASC").
		Find(&modules).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal ambil module"})
	}

	return c.JSON(fiber.Map{"data": modules})
}

func CreateModule(c *fiber.Ctx) error {
	var body struct {
		ModuleGroupID string  `json:"module_group_id"`
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		YoutubeID     string  `json:"youtube_id"`
		IsActive      bool    `json:"is_active"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	module := models.Module{
		ModuleGroupID: uuid.MustParse(body.ModuleGroupID),
		Title:         body.Title,
		Description:   body.Description,
		YoutubeID:     body.YoutubeID,
		IsActive:      body.IsActive,
	}

	if err := db.DB.Create(&module).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal buat module"})
	}

	return c.JSON(fiber.Map{
		"message": "Module berhasil dibuat",
		"data":    module,
	})
}

func UpdateModule(c *fiber.Ctx) error {
	id := c.Params("id")

	var module models.Module
	if err := db.DB.First(&module, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Module tidak ditemukan"})
	}

	var body struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
		YoutubeID   string  `json:"youtube_id"`
		IsActive    bool    `json:"is_active"`
	}

	c.BodyParser(&body)

	db.DB.Model(&module).Updates(body)

	return c.JSON(fiber.Map{
		"message": "Module berhasil diupdate",
	})
}

func DeleteModule(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := db.DB.Delete(&models.Module{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus module"})
	}

	return c.JSON(fiber.Map{
		"message": "Module berhasil dihapus",
	})
}
