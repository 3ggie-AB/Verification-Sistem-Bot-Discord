package handlers

import (
	"github.com/gofiber/fiber/v2"
	"crypto-member/db"
	"crypto-member/config"
	"github.com/bwmarrin/discordgo"
	"crypto-member/models"
	"fmt"
	"os"
	"strings"
	"strconv"
	"path/filepath"
	"math/rand"
	"time"
	"math"
)

func formatRupiah(amount float64) string {
    s := fmt.Sprintf("%.0f", amount) // tanpa desimal
    n := len(s)
    if n <= 3 {
        return s
    }

    var parts []string
    for n > 3 {
        parts = append([]string{s[n-3 : n]}, parts...)
        n -= 3
    }
    parts = append([]string{s[:n]}, parts...)
    return strings.Join(parts, ".")
}

func mustOpenFile(path string) *os.File {
    f, err := os.Open(path)
    if err != nil {
        panic(err)
    }
    return f
}

func pricePerMonth(monthCount int) float64 {
	raw := 0.0

	switch {
	case monthCount >= 12:
		raw = 1800000.0 / 12
	case monthCount >= 6:
		raw = 1000000.0 / 6
	case monthCount >= 3:
		raw = 550000.0 / 3
	default:
		raw = 200000.0
	}

	// pembulatan ke atas ke ribuan
	return math.Ceil(raw/1000) * 1000
}

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
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".heic": true, ".pdf": true,
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
	var originalAmount float64
	if monthCount >= 10000 {
		pricePerMonth := pricePerMonth(monthCount)
		originalAmount = pricePerMonth * float64(monthCount)
	}else{
		originalAmount = 2_500_000.0
	}

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

	botToken := config.Get("BOT_TOKEN")
	notifChannel := config.Get("NOTIF_CHANNEL_ID")

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Gagal bikin session Discord:", err)
	} else {
		if err := dg.Open(); err != nil {
			fmt.Println("Gagal buka session Discord:", err)
		} else {
			defer dg.Close()

			f, err := os.Open(filename)
			if err != nil {
				fmt.Println("Gagal buka file bukti:", err)
			} else {
				defer f.Close()

				embed := &discordgo.MessageEmbed{
					Title: "💰 Checkout Membership Baru :",
					Description: fmt.Sprintf(
						"**Nama User : %s\nJumlah bulan: %d\nTotal: Rp %s**\n\nCek Bukti : https://cryptolabsakademi.site/%s",
						user.Username, monthCount, formatRupiah(finalAmount), buktiPath,
					),
					Color:     0x00FF00,
					Timestamp: time.Now().Format(time.RFC3339),
				}

				_, err = dg.ChannelMessageSendComplex(notifChannel, &discordgo.MessageSend{
					Embeds: []*discordgo.MessageEmbed{embed},
				})
				if err != nil {
					fmt.Println("Gagal kirim notif ke channel:", err)
				} else {
					fmt.Println("Notif sukses terkirim!")
				}
			}
		}
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

