package service

import (
	"log"
	"strings"
	"github.com/bwmarrin/discordgo"
	"crypto-member/config"
)

func StartBot() {
	var BotToken = config.Get("BOT_TOKEN")

	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	dg.AddHandler(InteractionCreate)
	dg.AddHandler(ModalSubmit)
	dg.AddHandler(MessageCreate)

	// intents supaya bisa DM & lihat member
	dg.Identify.Intents = discordgo.IntentsDirectMessages |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsMessageContent

	if err := dg.Open(); err != nil {
		log.Fatal("Error opening connection:", err)
	}
	// StartMemberExpiryChecker(dg)

	log.Println("Bot sudah online!")
	select {} // biar bot jalan terus
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// pastikan DM
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		log.Println("Gagal buat DM:", err)
		return
	}

	content := strings.TrimSpace(m.Content)
	if strings.ToLower(content) != "" {
		s.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
			Content: "Perkenalkan Nama Saya Kopi Americano, Official Bot dari CryptoLabs Akademi. di buat Untuk Redeem Kode Aktifasi Member secara Otomatis.",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Klik dan Masukkan Kode Aktifasi Member",
							Style:    discordgo.PrimaryButton,
							CustomID: "input_token",
						},
					},
				},
			},
		})
		return
	}
}
