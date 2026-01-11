package handlers

import (
	"crypto-member/models"
	"crypto-member/db"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GET /api/users
func GetUsers(c *fiber.Ctx) error {
	db := db.DB

	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// jangan kirim password ke frontend
	for i := range users {
		users[i].Password = ""
	}

	return c.JSON(fiber.Map{"data": users})
}

// POST /api/users
func CreateUser(c *fiber.Ctx) error {
	db := db.DB

	input := new(struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role"`
		NomorHp  string `json:"phone"`
		Nama     string `json:"name"`
	})

	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}

	user := models.User{
		Username:      input.Username,
		Email:         input.Email,
		Password:      string(hash),
		Role:          input.Role,
		NomorHp:       &input.NomorHp,
		NamaLengkap:   &input.Nama,
		MemberExpiredAt: nil,
	}

	if err := db.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	user.Password = "" // jangan kirim password
	return c.Status(201).JSON(fiber.Map{"data": user})
}

// PUT /api/users/:id
func UpdateUser(c *fiber.Ctx) error {
	db := db.DB
	id := c.Params("id")

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	input := new(struct {
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Password *string `json:"password"`
		Role     *string `json:"role"`
		NomorHp  *string `json:"phone"`
		Nama     *string `json:"name"`
	})

	if err := c.BodyParser(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if input.Username != nil {
		user.Username = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil && *input.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		user.Password = string(hash)
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.NomorHp != nil {
		user.NomorHp = input.NomorHp
	}
	if input.Nama != nil {
		user.NamaLengkap = input.Nama
	}

	if err := db.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	user.Password = ""
	return c.JSON(fiber.Map{"data": user})
}

// DELETE /api/users/:id
func DeleteUser(c *fiber.Ctx) error {
	db := db.DB
	id := c.Params("id")

	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	if user.Role == "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "admin tidak bisa dihapus"})
	}

	if err := db.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "user deleted"})
}
