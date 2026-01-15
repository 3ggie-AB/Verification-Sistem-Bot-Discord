package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"

	"crypto-member/config"
	"crypto-member/service"
)

func main() {
	config.LoadEnv()

	botToken := config.Get("BOT_TOKEN")
	if botToken == "" {
		log.Fatal("BOT_TOKEN kosong")
	}

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal(err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	log.Println("Discord connected")

	// 🔥 RUN SEKALI SAAT START
	runCheck(dg)

	// ⏱️ LOOP CRON
	interval := 4 * time.Hour
	for {
		time.Sleep(interval)
		runCheck(dg)
	}
}

func runCheck(dg *discordgo.Session) {
	log.Println("Start expired check")

	service.CheckAndRemoveExpiredMembers(dg)

	log.Println("Expired check done")
}
