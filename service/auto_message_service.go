package service

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/bwmarrin/discordgo"
)

func StartAutoMessager(dg *discordgo.Session) {
	ticker := time.NewTicker(1 * time.Minute)

	for {
		<-ticker.C
		processAutoMessages(dg)
	}
}

func processAutoMessages(dg *discordgo.Session) {
	var list []models.AutoMessager

	err := db.DB.Where("is_active = ?", true).Find(&list).Error
	if err != nil {
		log.Println("Error get automessager:", err)
		return
	}

	now := time.Now()

	for _, am := range list {

		loc, err := time.LoadLocation(am.Timezone)
		if err != nil {
			loc = time.Local
		}

		current := now.In(loc)
		currentDay := current.Weekday().String()[:3]
		currentTime := current.Format("15:04")

		// Parse Days JSON
		var days []string
		if err := json.Unmarshal(am.DaysOfWeek, &days); err != nil {
			continue
		}

		if !containsDay(days, currentDay) {
			continue
		}

		if am.RunTime == nil || *am.RunTime != currentTime {
			continue
		}

		// 🔥 Anti Double Send Protection
		if am.LastRunAt != nil {
			last := am.LastRunAt.In(loc)
			if last.Format("2006-01-02 15:04") == current.Format("2006-01-02 15:04") {
				continue
			}
		}

		sendAutoMessage(dg, am)
	}
}

func containsDay(days []string, today string) bool {
	for _, d := range days {
		if strings.EqualFold(d, today) {
			return true
		}
	}
	return false
}

func sendAutoMessage(dg *discordgo.Session, am models.AutoMessager) {

	if am.ChannelID == nil {
		return
	}

	_, err := dg.ChannelMessageSend(*am.ChannelID, am.Message)
	if err != nil {
		log.Println("Failed auto message:", err)
		return
	}

	now := time.Now()

	db.DB.Model(&am).Updates(map[string]interface{}{
		"last_run_at": now,
	})

	log.Println("Auto message sent:", am.Name)
}