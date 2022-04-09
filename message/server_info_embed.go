package message

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
)

func CreateServerInfoEmbed(info *aternos.ServerInfo) *discordgo.MessageEmbed {
	if info.DynIP == "" {
		info.DynIP = "unavailable"
	}

	return &discordgo.MessageEmbed{
		Title:       "Server info",
		Description: fmt.Sprintf("Server '%s' is currently **%s**.", info.Name, info.StatusLabel),
		Color:       colorMap[info.Status],
		URL:         "https://aternos.org/server/",
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{
				Name:   "Players online",
				Value:  fmt.Sprintf("%d/%d", info.Players, info.MaxPlayers),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Detected problems",
				Value:  fmt.Sprintf("%d", info.Problems),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Installed software",
				Value:  fmt.Sprintf("%s v%s", info.Software, info.Version),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Server address",
				Value:  fmt.Sprintf("`%s:%d`", info.Address, info.Port),
				Inline: true,
			},
			&discordgo.MessageEmbedField{
				Name:   "Dyn IP",
				Value:  fmt.Sprintf("`%s`", info.DynIP),
				Inline: true,
			},
		},
	}
}
