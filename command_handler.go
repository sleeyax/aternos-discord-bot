package aternos_discord_bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"github.com/sleeyax/aternos-discord-bot/database/models"
	"log"
	"strings"
)

// handleCommands responds to incoming interactive commands on discord.
func (ab *Bot) handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData()

	// TODO: improve error handling

	switch command.Name {
	case PingCommand:
		respondWithText(s, i, formatMessage("Pong!", normal))
	case ConfigureCommand:
		if ab.Database == nil {
			respondWithText(s, i, formatMessage("Command unavailable (no database configured).", warning))
			break
		}

		options := optionsToMap(command.Options)

		err := ab.Database.SaveServerSettings(&models.ServerSettings{
			GuildID:       i.GuildID,
			SessionCookie: options[SessionOption].StringValue(),
			ServerCookie:  options[ServerOption].StringValue(),
		})
		if err != nil {
			log.Printf("failed to save configuration: %e", err)
			respondWithText(s, i, formatMessage("Failed to save configuration.", danger))
			break
		}

		ab.deleteWorker(i.GuildID)

		respondWithText(s, i, formatMessage("Configuration changed successfully.", success))
	case StatusCommand:
		fallthrough
	case InfoCommand:
		fallthrough
	case PlayersCommand:
		w, err := ab.getWorker(i.GuildID)
		if err != nil {
			log.Printf("failed to get worker: %e", err)
			respondWithText(s, i, formatMessage("Failed to get worker", danger))
			break
		}

		serverInfo, err := w.GetServerInfo()
		if err != nil {
			if err == aternos.UnauthenticatedError {
				respondWithText(s, i, formatMessage("Invalid credentials. Use `/configure` to reconfigure the bot.", danger))
				// delete the worker from the pool, so we can re-create it once the next discord command is received (hopefully with valid credentials, then)
				ab.deleteWorker(i.GuildID)
			} else {
				log.Printf("failed to get server info: %s", err)
				respondWithText(s, i, formatMessage("Failed to get server info", danger))
			}
			break
		}

		switch command.Name {
		case StatusCommand:
			// s.ChannelMessageSendEmbed(i.ChannelID, message.CreateServerInfoEmbed(serverInfo))
			break
		case InfoCommand:
			respondWithText(s, i, formatMessage(fmt.Sprintf("Server '%s' is currently **%s**.", serverInfo.Name, serverInfo.StatusLabel), info))
		case PlayersCommand:
			if len(serverInfo.PlayerList) == 0 {
				respondWithText(s, i, formatMessage("No one is playing right now :(", info))
			} else {
				respondWithText(s, i, formatMessage(fmt.Sprintf("Active players: %s.", strings.Join(serverInfo.PlayerList, ", ")), info))
			}
		}
	default:
		respondWithText(s, i, formatMessage("**Unknown command!**", danger))
	}
}
