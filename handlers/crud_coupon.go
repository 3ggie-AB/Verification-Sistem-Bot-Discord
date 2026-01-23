package handlers

import (
	"github.com/gofiber/fiber/v2"
	"crypto-member/db"
	"crypto-member/models"
	"time"
	"strconv"
	"fmt"
)

func mustAdmin(c *fiber.Ctx) (*models.User, error) {
	user := c.Locals("user").(*models.User)
	if user.Role != "admin" {
		return nil, fiber.NewError(403, "Akses ditolak")
	}
	return user, nil
}

func GetCoupons(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}

	var coupons []models.Coupon
	if err := db.DB.Find(&coupons).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data coupon"})
	}

	return c.JSON(fiber.Map{"data": coupons})
}

func GetCouponByID(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}

	id := c.Params("id")
	var coupon models.Coupon
	if err := db.DB.First(&coupon, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Coupon tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"coupon": coupon})
}

func CreateCoupon(c *fiber.Ctx) error {
	if _, err := mustAdmin(c); err != nil {
		return err
	}

	type req struct {
		Code        string   `json:"code"`
		Type        string   `json:"type"`
		Value       float64  `json:"value"`
		MaxDiscount *float64 `json:"max_discount"`
		Quota       uint     `json:"quota"`
		Trigger     *string  `json:"trigger"`     // 🔥
		ExpiredAt   *string  `json:"expired_at"`
		IsActive    *bool    `json:"is_active"`
		MinMonth    *uint    `json:"min_month"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	if body.Code == "" || body.Value <= 0 || body.Quota == 0 ||
		(body.Type != "percent" && body.Type != "fixed") {
		return c.Status(400).JSON(fiber.Map{"error": "Field wajib tidak valid"})
	}

	var expiredAt *time.Time
	if body.ExpiredAt != nil {
		t, err := time.Parse("2006-01-02", *body.ExpiredAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format expired_at salah"})
		}
		expiredAt = &t
	}

	isActive := true
	if body.IsActive != nil {
		isActive = *body.IsActive
	}

	minMonth := uint(0)
	if body.MinMonth != nil {
		minMonth = *body.MinMonth
	}

	coupon := models.Coupon{
		Code:        body.Code,
		Type:        body.Type,
		Value:       body.Value,
		MaxDiscount: body.MaxDiscount,
		Quota:       body.Quota,
		Trigger:     body.Trigger, // 🔥
		IsActive:    isActive,
		MinMonth:    minMonth,
		ExpiredAt:   expiredAt,
		CreatedAt:   time.Now(),
	}

	if err := db.DB.Create(&coupon).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal membuat coupon"})
	}

	return c.JSON(fiber.Map{"message": "Coupon berhasil dibuat", "coupon": coupon})
}

func UpdateCoupon(c *fiber.Ctx) error {
	if _, err := mustAdmin(c); err != nil {
		return err
	}

	id := c.Params("id")
	var coupon models.Coupon
	if err := db.DB.First(&coupon, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Coupon tidak ditemukan"})
	}

	type req struct {
		Code        *string  `json:"code"`
		Type        *string  `json:"type"`
		Value       *float64 `json:"value"`
		MaxDiscount *float64 `json:"max_discount"`
		Quota       *uint    `json:"quota"`
		Trigger     *string  `json:"trigger"` // 🔥
		ExpiredAt   *string  `json:"expired_at"`
		IsActive    *bool    `json:"is_active"`
	    MinMonth    *uint    `json:"min_month"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	updates := map[string]interface{}{}

	if body.Code != nil {
		updates["code"] = *body.Code
	}
	if body.Type != nil {
		updates["type"] = *body.Type
	}
	if body.Value != nil {
		updates["value"] = *body.Value
	}
	if body.MaxDiscount != nil {
		updates["max_discount"] = body.MaxDiscount
	}
	if body.Quota != nil {
		updates["quota"] = *body.Quota
	}
	if body.Trigger != nil {
		updates["trigger"] = body.Trigger // 🔥
	}
	if body.ExpiredAt != nil {
		t, err := time.Parse("2006-01-02", *body.ExpiredAt)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format expired_at salah"})
		}
		updates["expired_at"] = t
	}
	if body.IsActive != nil {
		updates["is_active"] = *body.IsActive
	}
	if body.MinMonth != nil {
		updates["min_month"] = *body.MinMonth
	}

	if err := db.DB.Model(&coupon).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal update coupon"})
	}

	return c.JSON(fiber.Map{"message": "Coupon berhasil diupdate"})
}

func DeleteCoupon(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)
	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Akses ditolak"})
	}

	id := c.Params("id")
	if err := db.DB.Delete(&models.Coupon{}, id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menghapus coupon"})
	}

	return c.JSON(fiber.Map{"message": "Coupon berhasil dihapus"})
}

func CheckCouponByCode(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Kode coupon wajib diisi",
		})
	}

	monthStr := c.Query("month")
	if monthStr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bulan wajib diisi",
		})
	}

	monthInt, err := strconv.Atoi(monthStr)
	if err != nil || monthInt < 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bulan harus berupa angka",
		})
	}

	var coupon models.Coupon
	if err := db.DB.Where("code = ?", code).First(&coupon).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Coupon tidak ditemukan",
		})
	}

	now := time.Now()

	// validasi coupon
	if !coupon.IsActive {
		return c.Status(400).JSON(fiber.Map{
			"error": "Coupon tidak aktif",
		})
	}

	month := uint(monthInt)
	if coupon.MinMonth > 0 && month < coupon.MinMonth {
		textBulan := fmt.Sprintf("%d Bulan", month)
		if month >= 1000 {
			textBulan = "Lifetime"
		}
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("Coupon ini tidak bisa di Gunakan minimal Berlangganan %s", textBulan),
		})
	}

	if coupon.ExpiredAt != nil && coupon.ExpiredAt.Before(now) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Coupon sudah kadaluarsa",
		})
	}

	if coupon.UsedCount >= coupon.Quota {
		return c.Status(400).JSON(fiber.Map{
			"error": "Quota coupon sudah habis",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Coupon valid",
		"coupon": fiber.Map{
			"id":            coupon.ID,
			"code":          coupon.Code,
			"type":          coupon.Type,
			"value":         coupon.Value,
			"max_discount":  coupon.MaxDiscount,
			"quota":         coupon.Quota,
			"used_count":    coupon.UsedCount,
			"expired_at":    coupon.ExpiredAt,
		},
	})
}
