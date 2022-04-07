package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"github.com/sleeyax/aternos-discord-bot/database/models"
	"github.com/sleeyax/aternos-discord-bot/message"
	"strings"
)

// handleCommands responds to incoming interactive commands on discord.
func (ab *Bot) handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	command := i.ApplicationCommandData()

	// wrap functions around our utilities to make life easier
	sendText := func(content string) {
		respondWithText(s, i, content)
	}
	sendHiddenText := func(content string) {
		respondWithHiddenText(s, i, content)
	}
	sendErrorText := func(content string, err error) {
		respondWithError(s, i, content, err)
	}

	switch command.Name {
	case PingCommand:
		sendHiddenText(message.FormatDefault("Pong!"))
	// TODO: only allow admin account to reconfigure
	case ConfigureCommand:
		if ab.Database == nil {
			sendText(message.FormatWarning("Command unavailable (no database configured)."))
			break
		}

		options := optionsToMap(command.Options)

		err := ab.Database.SaveServerSettings(&models.ServerSettings{
			GuildID:       i.GuildID,
			SessionCookie: options[SessionOption].StringValue(),
			ServerCookie:  options[ServerOption].StringValue(),
		})
		if err != nil {
			sendErrorText("Failed to save configuration.", err)
			break
		}

		sendText(message.FormatSuccess("Configuration changed successfully."))
	case StatusCommand:
		fallthrough
	case InfoCommand:
		fallthrough
	case PlayersCommand:
		w, err := ab.getWorker(i.GuildID)
		if err != nil {
			sendErrorText("Failed to get worker", err)
			break
		}

		serverInfo, err := w.GetServerInfo()
		if err != nil {
			if err == aternos.UnauthenticatedError {
				sendText(message.FormatError("Invalid credentials. Use `/configure` to reconfigure the bot."))
				break
			}

			sendErrorText("Failed to get server info", err)

			break
		}

		switch command.Name {
		case InfoCommand:
			respondWithEmbeds(s, i, []*discordgo.MessageEmbed{
				message.CreateServerInfoEmbed(serverInfo),
			})
		case StatusCommand:
			sendText(message.FormatInfo("Server '%s' is currently **%s**.", serverInfo.Name, serverInfo.StatusLabel))
		case PlayersCommand:
			if len(serverInfo.PlayerList) == 0 {
				sendText(message.FormatInfo("No players online right now."))
				break
			}

			sendText(message.FormatInfo("Active players: %s.", strings.Join(serverInfo.PlayerList, ", ")))
		}
	default:
		sendText(message.FormatWarning("Command unavailable. Please try again later or refresh your discord client `CTRL + R`"))
	}
}
