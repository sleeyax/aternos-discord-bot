package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
)

func respondWithText(s *discordgo.Session, i *discordgo.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

// handleCommands responds to incoming interactive commands on discord.
func (ab *Bot) handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData()

	switch command.Name {
	case PingCommand:
		respondWithText(s, i, "Pong!")
	default:
		respondWithText(s, i, "**Unknown command!**")
	}
}
