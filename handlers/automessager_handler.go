package handlers

import (
	"encoding/json"
	"time"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AutoMessagerRequest struct {
	Name       string   `json:"name"`
	Message    string   `json:"message"`
	BotID      string   `json:"bot_id"`
	ServerID   string   `json:"server_id"`
	ChannelID  string   `json:"channel_id"`
	RunTime    string   `json:"run_time"`    // "08:30"
	DaysOfWeek []string `json:"days_of_week"` // ["Mon","Tue"]
	Timezone   string   `json:"timezone"`
	Image      *string  `json:"image"`
}

//
// CREATE
//
func CreateAutoMessager(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	var req AutoMessagerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.Message == "" || req.BotID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing required fields"})
	}

	daysJSON, err := json.Marshal(req.DaysOfWeek)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid days_of_week"})
	}

	id := uuid.New().String()

	auto := models.AutoMessager{
		ID:         id,
		Name:       req.Name,
		Message:    req.Message,
		Image:      req.Image,
		BotID:      req.BotID,
		ServerID:   req.ServerID,
		ChannelID:  &req.ChannelID,
		RunTime:    &req.RunTime,
		DaysOfWeek: daysJSON,
		Timezone:   req.Timezone,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if auto.Timezone == "" {
		auto.Timezone = "Asia/Jakarta"
	}

	if err := db.DB.Create(&auto).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create automessager"})
	}

	return c.JSON(auto)
}

//
// GET ALL
//
func GetAutoMessagers(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	var list []models.AutoMessager

	if err := db.DB.Order("created_at desc").Find(&list).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch"})
	}

	return c.JSON(list)
}

//
// GET BY ID
//
func GetAutoMessagerByID(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	id := c.Params("id")

	var am models.AutoMessager
	if err := db.DB.First(&am, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}

	return c.JSON(am)
}

//
// UPDATE
//
func UpdateAutoMessager(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	id := c.Params("id")

	var am models.AutoMessager
	if err := db.DB.First(&am, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}

	var req AutoMessagerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	daysJSON, _ := json.Marshal(req.DaysOfWeek)

	updateData := map[string]interface{}{
		"name":        req.Name,
		"message":     req.Message,
		"bot_id":      req.BotID,
		"server_id":   req.ServerID,
		"channel_id":  req.ChannelID,
		"run_time":    req.RunTime,
		"days_of_week": daysJSON,
		"timezone":    req.Timezone,
		"image":       req.Image,
		"updated_at":  time.Now(),
	}

	if err := db.DB.Model(&am).Updates(updateData).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update"})
	}

	return c.JSON(fiber.Map{"message": "Updated successfully"})
}

//
// TOGGLE ACTIVE
//
func ToggleAutoMessager(c *fiber.Ctx) error {
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	id := c.Params("id")

	var am models.AutoMessager
	if err := db.DB.First(&am, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found"})
	}

	am.IsActive = !am.IsActive

	if err := db.DB.Save(&am).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to toggle"})
	}

	return c.JSON(am)
}

//
// DELETE
//
func DeleteAutoMessager(c *fiber.Ctx) error {	
	admin := c.Locals("user").(*models.User)

	if admin.Role != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Akses ditolak",
		})
	}

	id := c.Params("id")

	if err := db.DB.Delete(&models.AutoMessager{}, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete"})
	}

	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}