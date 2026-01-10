package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"crypto-member/config"
	"crypto-member/service"
)

func main() {
	BotToken := config.Get("BOT_TOKEN")
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal(err)
	}

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	// langsung check expired satu kali
	service.CheckAndRemoveExpiredMembers(dg)

	log.Println("Cek expired selesai")
}
