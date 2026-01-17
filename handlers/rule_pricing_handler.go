package handlers

import (
	"strconv"
	"strings"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"crypto-member/config"
	"crypto-member/db"
	"crypto-member/models"
)

//
// ========================
// UTIL & CONFIG
// ========================
//

func isAturanDBEnabled() bool {
	return strings.ToLower(config.Get("ATURAN_DB")) == "true"
}

// fallback hardcode
func priceFallback(month int) float64 {
	switch {
	case month >= 12:
		return 1_800_000
	case month >= 6:
		return 1_000_000
	case month >= 3:
		return 550_000
	default:
		return 200_000
	}
}

func priceFromDB(monthCount int) (float64, error) {
	var rule models.RulePricing

	// Cari rule yang cocok
	q := db.DB.
		Where("is_active = true").
		Where("min_month <= ?", monthCount).
		Order("min_month DESC")

	// Kalau max_month ada, cek max_month >= monthCount
	q = q.Where("(max_month IS NULL OR max_month >= ?)", monthCount)

	if err := q.First(&rule).Error; err != nil {
		return 0, fmt.Errorf("Aturan pricing tidak ditemukan")
	}

	// Hitung proporsional
	baseMonths := rule.MinMonth
	if baseMonths == 0 {
		baseMonths = 1
	}

	price := rule.TotalPrice / float64(baseMonths) * float64(monthCount)
	return price, nil
}

type UpdateRulePricingRequest struct {
    MinMonth   int
    MaxMonth   *int
    TotalPrice float64
    IsActive   bool
}

// GET /membership/pricing?months=1,3,6,12,1000
func GetMembershipPricing(c *fiber.Ctx) error {
	raw := c.Query("months")
	if raw == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Query months wajib diisi",
		})
	}

	monthStrings := strings.Split(raw, ",")
	result := make(map[string]interface{})

	for _, m := range monthStrings {
		month, err := strconv.Atoi(strings.TrimSpace(m))
		if err != nil || month <= 0 {
			continue
		}

		var price float64
		if isAturanDBEnabled() {
			price, err = priceFromDB(month)
			if err != nil {
				result[m] = fiber.Map{
					"error": err.Error(),
				}
				continue
			}
		} else {
			price = priceFallback(month)
		}

		result[m] = fiber.Map{
			"price":         price,
			"price_rupiah":  formatRupiah(price),
		}
	}

	return c.JSON(fiber.Map{
		"aturan_db": isAturanDBEnabled(),
		"data":      result,
	})
}

//
// ========================
// ADMIN CRUD RULE PRICING
// ========================
//

// POST /admin/rule-pricing
func CreateRulePricing(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user.Role != "admin" {
		return fiber.ErrForbidden
	}

	var req models.RulePricing
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if req.MinMonth <= 0 || req.TotalPrice <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "min_month dan total_price wajib valid",
		})
	}

	max := 9999
	if req.MaxMonth != nil {
		max = *req.MaxMonth
	}
	
	ok, err := IsRangeAvailable(req.MinMonth, max, 0)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Range bulan %d–%d sudah ada, ubah rule lama dulu", req.MinMonth, max)
	}

	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()

	if err := db.DB.Create(&req).Error; err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"data": req,
	})
}

// GET /admin/rule-pricing
func GetAllRulePricing(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user.Role != "admin" {
		return fiber.ErrForbidden
	}

	var rules []models.RulePricing
	db.DB.Order("min_month ASC").Find(&rules)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Berhasil Mendapatkan Data Aturan Pembayaran",
		"data" : rules,
	})
}

// GET /admin/rule-pricing/:id
func GetRulePricingByID(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user.Role != "admin" {
		return fiber.ErrForbidden
	}

	id := c.Params("id")
	var rule models.RulePricing

	if err := db.DB.First(&rule, id).Error; err != nil {
		return fiber.ErrNotFound
	}

	return c.JSON(rule)
}

// PUT /admin/rule-pricing/:id
func UpdateRulePricing(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user.Role != "admin" {
		return fiber.ErrForbidden
	}

	id := c.Params("id")
	var rule models.RulePricing

	if err := db.DB.First(&rule, id).Error; err != nil {
		return fiber.ErrNotFound
	}

	var req UpdateRulePricingRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	// cek range dulu (opsional)
	max := 9999
	if req.MaxMonth != nil {
		max = *req.MaxMonth
	}
	ok, err := IsRangeAvailable(req.MinMonth, max, int(rule.ID))
	if err != nil {
		return err
	}
	if !ok {
		errms := fmt.Errorf("Range bulan %d–%d sudah ada, ubah rule lama dulu", req.MinMonth, max)
		c.JSON(fiber.Map{"error": errms})
	}

	// update field
	rule.MinMonth = req.MinMonth
	rule.MaxMonth = req.MaxMonth
	rule.TotalPrice = req.TotalPrice
	rule.IsActive = req.IsActive
	rule.UpdatedAt = time.Now()

	if err := db.DB.Save(&rule).Error; err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"data": rule})
}

// DELETE /admin/rule-pricing/:id
func DeleteRulePricing(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	if user.Role != "admin" {
		return fiber.ErrForbidden
	}

	id := c.Params("id")

	if err := db.DB.Delete(&models.RulePricing{}, id).Error; err != nil {
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{
		"message": "Rule pricing berhasil dihapus",
	})
}

func IsRangeAvailable(minMonth, maxMonth, ignoreID int) (bool, error) {
	var count int64
	query := db.DB.Model(&models.RulePricing{}).
		Where("is_active = true").
		Where("(min_month <= ? AND (max_month IS NULL OR max_month >= ?))", maxMonth, minMonth)
	
	if ignoreID > 0 {
		query = query.Where("id != ?", ignoreID)
	}

	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}
