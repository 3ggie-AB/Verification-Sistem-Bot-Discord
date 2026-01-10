package bot

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"crypto-member/config"
)

func Start() {
	dg, _ := discordgo.New("Bot " + config.Get("DISCORD_TOKEN"))

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.Bot {
			return
		}

		if len(m.Content) > 7 && m.Content[:7] == "!verify" {
			code := m.Content[8:]

			body, _ := json.Marshal(map[string]string{
				"code": code,
			})

			res, err := http.Post(
				"http://localhost:3000/api/discord/verify",
				"application/json",
				bytes.NewBuffer(body),
			)

			if err != nil || res.StatusCode != 200 {
				s.ChannelMessageSend(m.ChannelID, "❌ kode tidak valid")
				return
			}

			var data struct {
				RoleID string `json:"role_id"`
			}
			json.NewDecoder(res.Body).Decode(&data)

			s.GuildMemberRoleAdd(
				config.Get("DISCORD_GUILD_ID"),
				m.Author.ID,
				data.RoleID,
			)

			s.ChannelMessageSend(m.ChannelID, "✅ role aktif, welcome member 🚀")
		}
	})

	dg.Open()
	select {}
}
