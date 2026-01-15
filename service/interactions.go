package service

import (
	"github.com/bwmarrin/discordgo"
	"crypto-member/db"
	"crypto-member/config"
	"errors"
	"log"
	"time"
	"fmt"
)

func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	if i.MessageComponentData().CustomID == "input_token" {
		modal := &discordgo.InteractionResponseData{
			Title:    "Input Kode Aktivasi",
			CustomID: "token_modal",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "token_input",
							Label:       "Masukkan kode membermu",
							Style:       discordgo.TextInputShort,
							Placeholder: "Kode dari payment",
							Required:    true,
						},
					},
				},
			},
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseModal,
			Data: modal,
		})
	}
}

// 3. Modal submit
func ModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionModalSubmit {
		return
	}

	if i.ModalSubmitData().CustomID != "token_modal" {
		return
	}
	var username string
    var fullTag string
	
	memberID := ""
	if i.Member != nil && i.Member.User != nil {
		memberID = i.Member.User.ID
	} else if i.User != nil { // DM
		memberID = i.User.ID
	}

    if i.Member != nil && i.Member.User != nil {
        username = i.Member.User.Username
        fullTag = i.Member.User.Username + "#" + i.Member.User.Discriminator
    } else if i.User != nil {
        // DM
        username = i.User.Username
        fullTag = i.User.Username + "#" + i.User.Discriminator
    }

    log.Printf("User yang submit modal: %s (%s). DC ID : %s", username, fullTag, memberID)

	reply := &discordgo.InteractionResponseData{Flags: 1 << 6} // ephemeral

	// Ambil token dengan aman
	token, err := getModalValue(i, "token_input")
	if err != nil {
		reply.Content = "Gagal mengambil token 😢"
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: reply,
		})
		return
	}

	// Redeem token dari DB
	user, dccode, payment, err := db.RedeemDiscordCode(token, username, memberID)
	if err != nil {
		reply.Content = err.Error()
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: reply,
		})
		return
	}

	fmt.Println("Berhasil Step Redeem")
	// Assign role di server
	ServerID := config.Get("ID_SERVER")
	MemberRoleID := config.Get("ID_ROLE")
	LifetimeRoleID := config.Get("ID_ROLE_LIFETIME")

	role := ""
	if payment.MonthCount >= 1000 {
		role = LifetimeRoleID
	}else{
		role = MemberRoleID
	}
	
	_, err = s.GuildMember(ServerID, *user.IDDiscord)
	if err != nil {
		fmt.Println("Kamu belum join server 😢")
		reply.Content = "Kamu belum join server 😢"
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: reply,
		})
		return
	}

	err = s.GuildMemberRoleAdd(ServerID, *user.IDDiscord, role)
	if err != nil {
		fmt.Println("Gagal assign role 😢")
		reply.Content = "Gagal assign role 😢"
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: reply,
		})
		return
	}

	// Sukses
	reply.Content = "Token valid! Kamu sudah masuk ke server 🎉"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: reply,
	})

	now := time.Now()
	dccode.IsUsed = true
	dccode.UsedAt = &now
	if err := db.DB.Save(&dccode).Error; err != nil {
		reply.Content = "Gagal Update Token ke DB !!!"
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: reply,
		})
	}

	NotifBerhasilAktivasiGroup(
		username,
		memberID,
	)
	
	reply.Content = "Berhasil Menjadikanmu Member CryptoLabs Akademi ✨"
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: reply,
	})
	log.Printf("User %s redeem token %s sukses dan role diassign", user.IDDiscord, token)
}

func getModalValue(i *discordgo.InteractionCreate, customID string) (string, error) {
	data := i.ModalSubmitData()
	if len(data.Components) == 0 {
		return "", errors.New("modal data kosong")
	}

	for _, row := range data.Components {
		actionRow, ok := row.(*discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, comp := range actionRow.Components {
			input, ok := comp.(*discordgo.TextInput)
			if !ok {
				continue
			}
			if input.CustomID == customID {
				return input.Value, nil
			}
		}
	}
	return "", errors.New("modal input tidak ditemukan")
}

func NotifBerhasilAktivasiGroup(username, discordUserID string) {
	botToken := config.Get("BOT_TOKEN")
	notifChannel := config.Get("NOTIF_CHANNEL_ID")
	serverID := config.Get("ID_SERVER")

	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Println("Gagal bikin session Discord:", err)
		return
	}

	if err := dg.Open(); err != nil {
		log.Println("Gagal buka session Discord:", err)
		return
	}
	defer dg.Close()

	// ==========================
	// 1️⃣ NOTIF KE CHANNEL ADMIN
	// ==========================
	embed := &discordgo.MessageEmbed{
		Title: "✅ Aktivasi Member Baru",
		Description: fmt.Sprintf(
			"👤 **User:** %s\n🆔 **Discord ID:** `%s`\n🏠 **Server:** %s\n\nStatus: **AKTIF 🚀**",
			username,
			discordUserID,
			serverID,
		),
		Color:     0x00FF99,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = dg.ChannelMessageSendComplex(
		notifChannel,
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	)
	if err != nil {
		log.Println("Gagal kirim notif ke channel:", err)
	}

	// ==========================
	// 2️⃣ JAPRI USER (DM)
	// ==========================
	dm, err := dg.UserChannelCreate(discordUserID)
	if err != nil {
		log.Println("Gagal buat DM channel:", err)
		return
	}

	_, err = dg.ChannelMessageSend(dm.ID,
		"🎉 **Aktivasi Berhasil!**\n\n"+
			"Selamat, akun kamu sudah **resmi aktif** sebagai member **CryptoLabs Akademi** 🚀\n\n"+
			"Kalau ada kendala atau pertanyaan, feel free buat chat admin ya 🔥",
	)
	if err != nil {
		log.Println("Gagal kirim DM ke user:", err)
	}
}