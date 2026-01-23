package handlers

import (
	"github.com/gofiber/fiber/v2"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/google/uuid"
	"time"
	"gorm.io/gorm"
)

type ModuleGroupResponse struct {
	ID          uuid.UUID
	Title       string
	Description *string
	IsActive    bool
	ForMember   bool
	Modules     []ModuleResponse
}

type ModuleResponse struct {
	ID            uuid.UUID
	ModuleGroupID uuid.UUID
	Title         string
	Description   *string
	YoutubeID     string
	ForMember     bool
}

func GetModulesByGroup(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	now := time.Now()
	isAdmin := user.Role == "admin"

	areMemberActive := true
	if user.Role != "admin" && (user.MemberExpiredAt == nil || user.MemberExpiredAt.Before(now)) {
		areMemberActive = false
	}

	groupID := c.Params("group_id")

	var group models.ModuleGroup

	query := db.DB.
		Where("id = ?", groupID)

	// Filter group kalau bukan member
	if !areMemberActive {
		query = query.Where("for_member = ?", false)
	}

	// Preload modules + filter modules juga
	query = query.Preload("Modules", func(db *gorm.DB) *gorm.DB {
		if !areMemberActive {
			return db.Where("for_member = ?", false)
		}
		return db
	})

	if err := query.First(&group).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal ambil module group",
		})
	}

	// mapping ke DTO
	res := ModuleGroupResponse{
		ID:          group.ID,
		Title:       group.Title,
		Description: group.Description,
		IsActive:    group.IsActive,
		ForMember:   group.ForMember,
		Modules:     []ModuleResponse{},
	}

	for _, m := range group.Modules {
		moduleRes := ModuleResponse{
			ID:            m.ID,
			ModuleGroupID: m.ModuleGroupID,
			Title:         m.Title,
			Description:   m.Description,
			ForMember:     m.ForMember,
		}

		// 🔐 hanya admin yang dapat youtube_id
		if isAdmin {
			moduleRes.YoutubeID = m.YoutubeID
		}

		res.Modules = append(res.Modules, moduleRes)
	}

	return c.JSON(fiber.Map{
		"data": res,
	})
}

func CreateModule(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	var body struct {
		ModuleGroupID string  `json:"module_group_id"`
		Title         string  `json:"title"`
		Description   *string `json:"description"`
		YoutubeID     string  `json:"youtube_id"`
		IsActive      bool    `json:"is_active"`
		ForMember     bool    `json:"for_member"`
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
		ForMember:     body.ForMember,
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
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
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
		ForMember   bool    `json:"for_member"`
	}

	c.BodyParser(&body)

	db.DB.Model(&module).Updates(body)

	return c.JSON(fiber.Map{
		"message": "Module berhasil diupdate",
	})
}

func DeleteModule(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	id := c.Params("id")

	if err := db.DB.Delete(&models.Module{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal hapus module"})
	}

	return c.JSON(fiber.Map{
		"message": "Module berhasil dihapus",
	})
}
