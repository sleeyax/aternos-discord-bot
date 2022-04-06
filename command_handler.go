package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sleeyax/aternos-discord-bot/database/models"
	"log"
)

// handleCommands responds to incoming interactive commands on discord.
func (ab *Bot) handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData()

	switch command.Name {
	case PingCommand:
		respondWithText(s, i, formatMessage("Pong!", normal))
	case ConfigureCommand:
		if ab.Database == nil {
			respondWithText(s, i, formatMessage("Command unavailable (no database configured).", danger))
			return
		}

		options := optionsToMap(command.Options)

		err := ab.Database.SaveServerSettings(&models.ServerSettings{
			GuildID:       i.GuildID,
			SessionCookie: options[SessionOption].StringValue(),
			ServerCookie:  options[ServerOption].StringValue(),
		})
		if err != nil {
			log.Println(err)
			respondWithText(s, i, formatMessage("Failed to save configuration.", danger))
			return
		}

		respondWithText(s, i, formatMessage("Configuration changed successfully.", success))
	default:
		respondWithText(s, i, formatMessage("**Unknown command!**", danger))
	}
}
