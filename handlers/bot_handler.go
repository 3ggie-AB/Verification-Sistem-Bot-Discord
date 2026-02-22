package handlers

import (
	"time"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type BotRequest struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

//
// CREATE BOT
//
func CreateBot(c *fiber.Ctx) error {
	var req BotRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	if req.Name == "" || req.Token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name and Token required"})
	}

	bot := models.Bot{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Token:     req.Token,
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&bot).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create bot"})
	}

	// Jangan kirim token full ke client
	bot.Token = "********"

	return c.JSON(bot)
}

//
// GET ALL BOTS
//
func GetBots(c *fiber.Ctx) error {
	var bots []models.Bot

	if err := db.DB.Order("created_at desc").Find(&bots).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch bots"})
	}

	// Mask token
	for i := range bots {
		bots[i].Token = "********"
	}

	return c.JSON(bots)
}

//
// GET BOT BY ID
//
func GetBotByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var bot models.Bot
	if err := db.DB.First(&bot, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Bot not found"})
	}

	bot.Token = "********"

	return c.JSON(bot)
}

//
// UPDATE BOT
//
func UpdateBot(c *fiber.Ctx) error {
	id := c.Params("id")

	var bot models.Bot
	if err := db.DB.First(&bot, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Bot not found"})
	}

	var req BotRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	updateData := map[string]interface{}{
		"name":       req.Name,
		"updated_at": time.Now(),
	}

	// Kalau token dikirim, update
	if req.Token != "" {
		updateData["token"] = req.Token
	}

	if err := db.DB.Model(&bot).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update bot"})
	}

	return c.JSON(fiber.Map{"message": "Bot updated"})
}

//
// TOGGLE ACTIVE
//
func ToggleBot(c *fiber.Ctx) error {
	id := c.Params("id")

	var bot models.Bot
	if err := db.DB.First(&bot, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Bot not found"})
	}

	bot.IsActive = !bot.IsActive

	if err := db.DB.Save(&bot).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to toggle bot"})
	}

	return c.JSON(bot)
}

//
// DELETE BOT
//
func DeleteBot(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := db.DB.Delete(&models.Bot{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete"})
	}

	return c.JSON(fiber.Map{"message": "Bot deleted"})
}