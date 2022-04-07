package message

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
)

func CreateServerInfoEmbed(serverInfo *aternos.ServerInfo) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Server info",
		Description: fmt.Sprintf("Server '%s' is currently **%s**.", serverInfo.Name, serverInfo.StatusLabel),
		Color:       colorMap[serverInfo.Status],
		URL:         "https://aternos.org/server/",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Players online",
				Value:  fmt.Sprintf("%d/%d", serverInfo.Players, serverInfo.MaxPlayers),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Detected problems",
				Value:  fmt.Sprintf("%d", serverInfo.Problems),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Installed software",
				Value:  fmt.Sprintf("%s v%s", serverInfo.Software, serverInfo.Version),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server address",
				Value:  fmt.Sprintf("`%s:%d`", serverInfo.Address, serverInfo.Port),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Dyn IP",
				Value:  fmt.Sprintf("`%s`", serverInfo.DynIP),
				Inline: true,
			},
		},
	}
}
