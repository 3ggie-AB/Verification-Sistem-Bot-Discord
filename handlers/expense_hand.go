package handlers

import (
	"time"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"crypto-member/db"
	"crypto-member/models"
)

// ========================
// CREATE EXPENSE
// ========================
func CreateExpense(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	type req struct {
		Description string   `json:"description"`
		Amount      float64  `json:"amount"`
		Category    *string  `json:"category"`
		SpentAt     *string  `json:"spent_at"` // format: "YYYY-MM-DD"
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	if body.Description == "" || body.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Description dan Amount wajib diisi"})
	}

	spentAt := time.Now()
	if body.SpentAt != nil && *body.SpentAt != "" {
		t, err := time.Parse("2006-01-02", *body.SpentAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format spent_at salah, gunakan YYYY-MM-DD"})
		}
		spentAt = t
	}

	expense := models.Expense{
		Description: body.Description,
		Amount:      body.Amount,
		Category:    body.Category,
		SpentAt:     spentAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.DB.Create(&expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan pengeluaran"})
	}

	return c.JSON(fiber.Map{"message": "Pengeluaran berhasil dibuat", "expense": expense})
}

// ========================
// GET EXPENSES (FILTERABLE)
// ========================
func GetExpenses(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	// query optional: month=2026-01&category=operasional
	monthQuery := c.Query("month")
	categoryQuery := c.Query("category")

	var expenses []models.Expense
	q := db.DB

	if monthQuery != "" {
		parts := strings.Split(monthQuery, "-")
		if len(parts) == 2 {
			year, err1 := strconv.Atoi(parts[0])
			month, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
				end := start.AddDate(0, 1, 0)
				q = q.Where("spent_at >= ? AND spent_at < ?", start, end)
			}
		}
	}

	if categoryQuery != "" {
		q = q.Where("category = ?", categoryQuery)
	}

	if err := q.Order("spent_at DESC").Find(&expenses).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pengeluaran"})
	}

	return c.JSON(fiber.Map{"data": expenses})
}

// ========================
// GET EXPENSE BY ID
// ========================
func GetExpenseByID(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	id := c.Params("id")
	var expense models.Expense
	if err := db.DB.First(&expense, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pengeluaran tidak ditemukan"})
	}
	return c.JSON(fiber.Map{"expense": expense})
}

// ========================
// UPDATE EXPENSE
// ========================
func UpdateExpense(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	id := c.Params("id")
	var expense models.Expense
	if err := db.DB.First(&expense, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pengeluaran tidak ditemukan"})
	}

	type req struct {
		Description *string  `json:"description"`
		Amount      *float64 `json:"amount"`
		Category    *string  `json:"category"`
		SpentAt     *string  `json:"spent_at"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	updates := map[string]interface{}{}
	if body.Description != nil {
		updates["description"] = *body.Description
	}
	if body.Amount != nil && *body.Amount > 0 {
		updates["amount"] = *body.Amount
	}
	if body.Category != nil {
		updates["category"] = body.Category
	}
	if body.SpentAt != nil && *body.SpentAt != "" {
		t, err := time.Parse("2006-01-02", *body.SpentAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format spent_at salah"})
		}
		updates["spent_at"] = t
	}
	updates["updated_at"] = time.Now()

	if err := db.DB.Model(&expense).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update pengeluaran"})
	}

	return c.JSON(fiber.Map{"message": "Pengeluaran berhasil diupdate", "expense": expense})
}

// ========================
// DELETE EXPENSE
// ========================
func DeleteExpense(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}
	id := c.Params("id")
	if err := db.DB.Delete(&models.Expense{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus pengeluaran"})
	}
	return c.JSON(fiber.Map{"message": "Pengeluaran berhasil dihapus"})
}
