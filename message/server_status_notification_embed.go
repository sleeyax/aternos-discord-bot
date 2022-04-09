package message

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
)

func CreateServerStatusNotificationEmbed(info *aternos.ServerInfo) (*discordgo.MessageEmbed, error) {
	switch info.Status {
	case aternos.Online:
		return &discordgo.MessageEmbed{
			Title:       "Server is online",
			Description: fmt.Sprintf("Join now! Only %d seconds left.", info.Countdown),
			Color:       colorMap[aternos.Online],
			Fields: []*discordgo.MessageEmbedField{
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
		}, nil
	case aternos.Offline:
		return &discordgo.MessageEmbed{
			Title:       "Server is offline",
			Description: "The server is currently offline.",
			Color:       colorMap[aternos.Offline],
		}, nil
	default:
		return nil, fmt.Errorf("unknown server status code '%d'", info.Status)
	}
}
