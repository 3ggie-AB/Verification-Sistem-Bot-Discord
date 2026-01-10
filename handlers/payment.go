package handlers

import (
	"github.com/gofiber/fiber/v2"
	"crypto-member/db"
	"crypto-member/models"
	"time"
)

func GetPayments(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var payments []models.Payment
	query := db.DB.Preload("DiscordCode")

	if user.Role == "admin" {
		// admin lihat semua
		if err := query.Find(&payments).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal mengambil data payment",
			})
		}
	} else {
		// user hanya lihat punya sendiri
		if err := query.
			Where("user_id = ?", user.ID).
			Find(&payments).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Gagal mengambil data payment",
			})
		}
	}

	return c.JSON(fiber.Map{
		"data": payments,
	})
}

func ApprovePayment(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	paymentID := c.Params("id")

	var payment models.Payment
	if err := db.DB.First(&payment, paymentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Payment tidak ditemukan",
		})
	}

	if payment.Status == "paid" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Payment sudah di-approve",
		})
	}

	// update payment
	now := time.Now()
	if err := db.DB.Model(&payment).Updates(models.Payment{
		Status: "paid",
		PaidAt: &now,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal approve payment",
		})
	}

	// generate discord code SETELAH PAID
	code := generateUniqueDiscordCode()
	discordCode := models.DiscordCode{
		PaymentID: payment.ID,
		Code:      code,
		IsUsed:    false,
	}

	if err := db.DB.Create(&discordCode).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal membuat Discord Code",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment berhasil di-approve",
		"discord_code": code,
	})
}

func RejectPayment(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	paymentID := c.Params("id")

	type req struct {
		Reason string `json:"reason"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil || body.Reason == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Alasan penolakan wajib diisi",
		})
	}

	var payment models.Payment
	if err := db.DB.First(&payment, paymentID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Payment tidak ditemukan",
		})
	}

	if payment.Status != "pending" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Hanya payment pending yang bisa ditolak",
		})
	}

	if err := db.DB.Model(&payment).Updates(map[string]interface{}{
		"status":        "failed",
		"reject_reason": body.Reason,
	}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menolak payment",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Payment berhasil ditolak",
	})
}
