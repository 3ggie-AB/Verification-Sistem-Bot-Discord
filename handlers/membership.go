package handlers

import (
	"github.com/gofiber/fiber/v2"
	"crypto-member/db"
	"crypto-member/models"
	"fmt"
	"strings"
	"strconv"
	"path/filepath"
	"math/rand"
	"time"
)

func CheckoutMembership(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	method := c.FormValue("method")
	monthCount, _ := strconv.Atoi(c.FormValue("month_count"))
	couponCode := c.FormValue("coupon_code")

	if method == "" || monthCount <= 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Method dan month_count wajib diisi",
		})
	}

	// === UPLOAD FILE ===
	file, err := c.FormFile("bukti")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Bukti pembayaran wajib diupload",
		})
	}

	// validasi ext
	ext := filepath.Ext(file.Filename)
	allowed := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".pdf": true,
	}
	if !allowed[strings.ToLower(ext)] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Format bukti tidak didukung",
		})
	}

	// path simpan
	filename := fmt.Sprintf(
		"bukti/%d_%d%s",
		user.ID,
		time.Now().UnixNano(),
		ext,
	)

	if err := c.SaveFile(file, filename); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menyimpan bukti",
		})
	}

	// === HITUNG PAYMENT ===
	pricePerMonth := 200000.0
	originalAmount := pricePerMonth * float64(monthCount)

	finalAmount, coupon, err := applyCoupon(couponCode, originalAmount)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	trx := generateTransactionRef()
	buktiPath := filename

	payment := models.Payment{
		UserID:         user.ID,
		Method:         method,
		Status:         "pending",
		MonthCount:     uint(monthCount),
		OriginalAmount: originalAmount,
		Amount:         finalAmount,
		TransactionRef: &trx,
		Bukti:          &buktiPath,
	}

	if coupon != nil {
		payment.CouponID = &coupon.ID
		payment.Discount = originalAmount - finalAmount
	}

	if err := db.DB.Create(&payment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal membuat pembayaran",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Checkout berhasil, menunggu konfirmasi admin",
		"payment": payment,
	})
}

func generateTransactionRef() string {
	return fmt.Sprintf("TRX-%d", time.Now().UnixNano())
}

func generateDiscordCode() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 8)

	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func generateUniqueDiscordCode() string {
	for {
		code := generateDiscordCode()

		var count int64
		db.DB.Model(&models.DiscordCode{}).
			Where("code = ? AND is_used = true", code).
			Count(&count)

		if count == 0 {
			return code
		}
	}
}

func applyCoupon(code string, amount float64) (float64, *models.Coupon, error) {
	if code == "" {
		return amount, nil, nil
	}

	var coupon models.Coupon
	if err := db.DB.
		Where("code = ? AND is_active = true", code).
		First(&coupon).Error; err != nil {
		return amount, nil, fmt.Errorf("Coupon tidak valid")
	}

	if coupon.ExpiredAt != nil && coupon.ExpiredAt.Before(time.Now()) {
		return amount, nil, fmt.Errorf("Coupon sudah expired")
	}

	if coupon.UsedCount >= coupon.Quota {
		return amount, nil, fmt.Errorf("Coupon sudah habis")
	}

	discount := 0.0

	if coupon.Type == "percent" {
		discount = amount * coupon.Value / 100
		if coupon.MaxDiscount != nil && discount > *coupon.MaxDiscount {
			discount = *coupon.MaxDiscount
		}
	} else {
		discount = coupon.Value
	}

	finalAmount := amount - discount
	if finalAmount < 0 {
		finalAmount = 0
	}

	return finalAmount, &coupon, nil
}

