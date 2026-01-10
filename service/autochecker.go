package service

import (
	"log"
	"time"

	"crypto-member/config"
	"crypto-member/db"
	"crypto-member/models"

	"github.com/bwmarrin/discordgo"
)

// // StartMemberExpiryChecker: jalankan background checker
// func StartMemberExpiryChecker(dg *discordgo.Session) {
// 	ticker := time.NewTicker(10 * time.Minute) // cek tiap menit
// 	go func() {
// 		for range ticker.C {
// 			checkAndRemoveExpiredMembers(dg)
// 		}
// 	}()
// }

func CheckAndRemoveExpiredMembers(dg *discordgo.Session) {
	ServerID := config.Get("ID_SERVER")
	MemberRoleID := config.Get("ID_ROLE")

	var users []models.User
	if err := db.DB.Where("member_expired_at <= ?", time.Now()).Find(&users).Error; err != nil {
		log.Println("Gagal ambil user expired:", err)
		return
	}

	for _, u := range users {
		if u.IDDiscord == nil || *u.IDDiscord == "" {
			continue // skip kalau user belum connect Discord
		}

		err := dg.GuildMemberRoleRemove(ServerID, *u.IDDiscord, MemberRoleID)
		if err != nil {
			log.Printf("Gagal copot role user %s: %v", *u.NamaDiscord, err)
			continue
		}

		log.Printf("Role dicopot untuk user %s karena membership expired", *u.NamaDiscord)
	}
}
