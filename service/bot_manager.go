package service

import (
	"log"
	"sync"
	"time"

	"crypto-member/db"
	"crypto-member/models"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

var botSessions = map[string]*discordgo.Session{}
var botMutex sync.Mutex

func GetBotSession(botID string, token string) (*discordgo.Session, error) {

	botMutex.Lock()
	defer botMutex.Unlock()

	if session, ok := botSessions[botID]; ok {
		return session, nil
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	if err := dg.Open(); err != nil {
		return nil, err
	}

	botSessions[botID] = dg

	log.Println("Bot connected:", botID)

	return dg, nil
}

func RunAutoMessenger(db *gorm.DB) {

	var autos []models.AutoMessager
	db.Where("is_active = ?", true).Find(&autos)

	now := time.Now()

	for _, auto := range autos {

		if auto.NextRunAt != nil && now.Before(*auto.NextRunAt) {
			continue
		}

		sendAutoMessageDynamic(db, auto)
		updateNextRun(db, auto)
	}
}

func sendAutoMessageDynamic(db *gorm.DB, auto models.AutoMessager) {

	var bot models.Bot

	if err := db.First(&bot, "id = ?", auto.BotID).Error; err != nil {
		log.Println("Bot tidak ditemukan")
		return
	}

	session, err := GetBotSession(bot.ID, bot.Token)
	if err != nil {
		log.Println("Gagal connect bot:", err)
		return
	}

	target := auto.ServerID

	if auto.ChannelID != nil {
		target = *auto.ChannelID
	}

	_, err = session.ChannelMessageSend(target, auto.Message)

	if err != nil {
		log.Println("Gagal kirim message:", err)
	}
}

func StartAutoMessengerWorker() {

	ticker := time.NewTicker(30 * time.Second)

	for range ticker.C {
		RunAutoMessenger(db.DB)
	}
}

func updateNextRun(db *gorm.DB, auto models.AutoMessager) {

	if auto.RunTime == nil {
		return
	}

	loc, err := time.LoadLocation(auto.Timezone)
	if err != nil {
		loc = time.UTC
	}

	now := time.Now().In(loc)

	t, err := time.Parse("15:04", *auto.RunTime)
	if err != nil {
		return
	}

	next := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		loc,
	)

	if now.After(next) {
		next = next.Add(24 * time.Hour)
	}

	db.Model(&auto).Updates(map[string]interface{}{
		"next_run_at": next,
		"last_run_at": now,
	})
}