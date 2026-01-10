package handlers

import (
	"time"
	"crypto-member/db"
	"crypto-member/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"encoding/hex"
)

// =====================
// Register user
// =====================
func Register(c *fiber.Ctx) error {
	type req struct {
		Email        string  `json:"email"`
		Username     string  `json:"username"`
		Password     string  `json:"password"`
		NamaLengkap  *string `json:"nama_lengkap"`
		NamaDiscord  *string `json:"nama_discord"`
		NomorHp      *string `json:"nomor_hp"`
		From         *string `json:"from"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	// hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		Email:           body.Email,
		Username:        body.Username,
		Password:        string(hash),
		Role:            "user", // default role
		NamaLengkap:     body.NamaLengkap,
		NamaDiscord:     body.NamaDiscord,
		NomorHp:         body.NomorHp,
		From:            body.From,
		MemberExpiredAt: nil, // bisa diatur nanti
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// jangan kirim password kembali
	user.Password = ""
	return c.JSON(user)
}

// =====================
// Login user
// =====================
func Login(c *fiber.Ctx) error {
	type req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	var user models.User
	if err := db.DB.Where("username = ?", body.Username).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid password"})
	}

	// generate token baru
	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	user.Token = &token
	user.UpdatedAt = time.Now()
	db.DB.Save(&user)

	// jangan kirim password
	user.Password = ""
	return c.JSON(fiber.Map{
		"message": "login success",
		"user":    user,
		"token":   token,
	})
}

// =====================
// Update profile
// =====================
func UpdateProfile(c *fiber.Ctx) error {
	type req struct {
		NamaLengkap *string `json:"nama_lengkap"`
		NamaDiscord *string `json:"nama_discord"`
		NomorHp     *string `json:"nomor_hp"`
		From        *string `json:"from"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	user := c.Locals("user").(*models.User)

	// update fields kalau ada
	if body.NamaLengkap != nil {
		user.NamaLengkap = body.NamaLengkap
	}
	if body.NamaDiscord != nil {
		user.NamaDiscord = body.NamaDiscord
	}
	if body.NomorHp != nil {
		user.NomorHp = body.NomorHp
	}
	if body.From != nil {
		user.From = body.From
	}

	user.UpdatedAt = time.Now()
	db.DB.Save(user)

	user.Password = ""
	return c.JSON(user)
}

func Me(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)
	user.Password = ""
	return c.JSON(user)
}
