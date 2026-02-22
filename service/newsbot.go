package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"strings"

	"github.com/bwmarrin/discordgo"
	"crypto-member/models"
	"crypto-member/db"
)

func SendNewsToDiscord(dg *discordgo.Session, channelID string) {
	rand.Seed(time.Now().UnixNano()) // seed random

	for {
		var n models.CryptoNews
		if err := db.DB.Where("is_upload = ?", false).First(&n).Error; err != nil {
			time.Sleep(10 * time.Second)
			continue
		}

		lowerContent := strings.ToLower(n.ContentIndo)
		konten = strings.ReplaceAll(n.ContentIndo, "?", " - ")

		invalidPatterns := []string{
			"tidak dapat mengakses",
			"tidak bisa mengakses",
			"tidak bisa membuka link",
			"cannot access",
			"as an ai language model",
		}

		isInvalid := false	
		for _, pattern := range invalidPatterns {
			if strings.Contains(lowerContent, pattern) {
				isInvalid = true
				break
			}
		}

		if isInvalid {
			log.Println("AI gagal generate konten, skip:", n.Title)
			db.DB.Delete(&n)
			continue
		}

		if len(konten) > 4000 {
			konten = konten[:3997] + "..."
		}

		if len(konten) < 200 {
			konten += "\n\nLihat selengkapnya di: " + n.URL
		}

		// generate warna random
		randomColor := rand.Intn(0xFFFFFF + 1) // 0x000000 s/d 0xFFFFFF

		embed := &discordgo.MessageEmbed{
			Title:       n.Title,
			Description: konten,
			URL:         n.URL,
			Color:       randomColor, // <--- warna random
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Sumber: %s | %s", n.SourceName, n.SourceDomain),
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Image: &discordgo.MessageEmbedImage{
				URL: n.Thumbnail,
			},
		}

		_, err := dg.ChannelMessageSendEmbed(channelID, embed)
		if err != nil {
			log.Println("Failed to send news:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		db.DB.Model(&n).Update("is_upload", true)
		fmt.Println("Sent news to Discord:", n.Title)
		time.Sleep(60 * 60 * time.Second)
	}
}
