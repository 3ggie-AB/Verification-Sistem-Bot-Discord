package service

import (
	"crypto-member/db"
	"crypto-member/models"
	"strings"
	"encoding/json"
	"log"
)

func CreateAnnouncement(a models.Announcement, channels []string) error {

	// Simpan ke DB dulu
	if err := db.DB.Create(&a).Error; err != nil {
		return err
	}

	// Kirim async supaya tidak blocking
	go ProcessAnnouncement(a, channels)

	return nil
}

func GetAllAnnouncements() ([]models.Announcement, error) {
	var data []models.Announcement

	if err := db.DB.
		Order("created_at desc").
		Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func ProcessAnnouncement(a models.Announcement, channels []string) {

	users := GetUsersByTarget(a.Target)

	for _, user := range users {

		message := strings.ReplaceAll(a.Content, "[NAMA]", user.Username)

		for _, ch := range channels {

			switch ch {

			case "email":
				if err := SendEmail(user.Email, a.Title, message); err != nil {
					log.Println("Email gagal:", err)
				}

			case "discord":
				if user.IDDiscord != nil && *user.IDDiscord != "" {
					if err := SendDiscordDM(*user.IDDiscord, message); err != nil {
						log.Println("Discord DM gagal:", err)
					}
				}
			}
		}
	}
}

func GetUsersByTarget(targetJSON []byte) []models.User {

	var target map[string]interface{}
	json.Unmarshal(targetJSON, &target)

	var users []models.User

	switch target["audience"] {

	case "all":
		db.DB.Find(&users)

	case "active":
		db.DB.Where("status = ?", "active").Find(&users)

	case "expired":
		db.DB.Where("status = ?", "expired").Find(&users)

	default:
		db.DB.Find(&users)
	}

	return users
}

func SendDiscordDM(userID, message string) error {

	dm, err := Discord.UserChannelCreate(userID)
	if err != nil {
		return err
	}

	_, err = Discord.ChannelMessageSend(dm.ID, message)
	return err
}

func SendDiscordChannelMessage(channelID, message string) error {
	_, err := Discord.ChannelMessageSend(channelID, message)
	return err
}