package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/google/uuid"
)

func UpdateModuleProgress(c *fiber.Ctx) error {
	user := c.Locals("user").(*models.User)

	var body struct {
		ModuleID string `json:"module_id"`
		Status   string `json:"status"` // watching | completed
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	var progress models.ModuleProgress
	err := db.DB.
		Where("user_id = ? AND module_id = ?", user.ID, body.ModuleID).
		First(&progress).Error

	now := time.Now()

	if err != nil {
		progress = models.ModuleProgress{
			UserID:   user.ID,
			ModuleID: uuid.MustParse(body.ModuleID),
			Status:   body.Status,
			LastWatchedAt: &now,
		}
		db.DB.Create(&progress)
	} else {
		progress.Status = body.Status
		progress.LastWatchedAt = &now
		if body.Status == "completed" {
			progress.CompletedAt = &now
		}
		db.DB.Save(&progress)
	}

	return c.JSON(fiber.Map{
		"message": "Progress tersimpan",
	})
}
