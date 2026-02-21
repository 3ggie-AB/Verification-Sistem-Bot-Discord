package service

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

var Discord *discordgo.Session

func InitDiscord() {
	botToken := os.Getenv("BOT_TOKEN")

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal("Gagal buat Discord session:", err)
	}

	if err := dg.Open(); err != nil {
		log.Fatal("Gagal buka Discord connection:", err)
	}

	Discord = dg
	log.Println("Discord bot connected 🚀")
}