package handlers

import (
	"github.com/gofiber/fiber/v2"
	"crypto-member/db"
	"crypto-member/models"
)

type CreateRulePricingRequest struct {
	MinMonth   int   `json:"min_month"`
	MaxMonth   *int  `json:"max_month"`
	TotalPrice float64 `json:"total_price"`
	IsActive   *bool `json:"is_active"`
}
func CreateRulePricing(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	var req CreateRulePricingRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.MinMonth <= 0 || req.TotalPrice <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "MinMonth dan TotalPrice wajib diisi",
		})
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	rule := models.RulePricing{
		MinMonth:   req.MinMonth,
		MaxMonth:   req.MaxMonth,
		TotalPrice: req.TotalPrice,
		IsActive:   isActive,
	}

	if err := db.DB.Create(&rule).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Rule pricing berhasil dibuat",
		"data":    rule,
	})
}

func GetRulePricing(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	var rules []models.RulePricing

	if err := db.DB.
		Where("is_active = ?", true).
		Order("min_month asc").
		Find(&rules).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(rules)
}

func GetPriceByMonth(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	month, err := c.ParamsInt("month")
	if err != nil || month <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Month tidak valid",
		})
	}

	var rule models.RulePricing

	err = db.DB.
		Where("is_active = ?", true).
		Where("min_month <= ?", month).
		Where("(max_month IS NULL OR max_month >= ?)", month).
		Order("min_month desc").
		First(&rule).Error

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Rule pricing tidak ditemukan",
		})
	}

	return c.JSON(fiber.Map{
		"month":       month,
		"total_price": rule.TotalPrice,
		"rule":        rule,
	})
}

func UpdateRulePricing(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	id, _ := c.ParamsInt("id")

	var rule models.RulePricing
	if err := db.DB.First(&rule, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Rule tidak ditemukan",
		})
	}

	var req CreateRulePricingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	if req.MinMonth > 0 {
		rule.MinMonth = req.MinMonth
	}
	if req.TotalPrice > 0 {
		rule.TotalPrice = req.TotalPrice
	}
	if req.MaxMonth != nil {
		rule.MaxMonth = req.MaxMonth
	}
	if req.IsActive != nil {
		rule.IsActive = *req.IsActive
	}

	db.DB.Save(&rule)

	return c.JSON(fiber.Map{
		"message": "Rule pricing berhasil diupdate",
		"data":    rule,
	})
}
