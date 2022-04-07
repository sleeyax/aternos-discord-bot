package aternos_discord_bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"github.com/sleeyax/aternos-discord-bot/database"
	"github.com/sleeyax/aternos-discord-bot/database/models"
	"log"
	"net/http"
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
	case StartCommand:
		options, err := ab.createOptions(i.GuildID)

		if err != nil {
			if err == database.ErrDataNotFound {
				respondWithText(s, i, formatMessage("Server settings not found. Use `/configure` to get started.", warning))
				return
			}

			log.Printf("failed to get server settings %e", err)
			respondWithText(s, i, formatMessage("Failed to get server settings", danger))
			return
		}

		connection := &Connection{
			guildId: i.GuildID,
			api:     aternos.New(options),
		}

		respondWithText(s, i, formatMessage(fmt.Sprintf("Created options for %s", connection.guildId), success))
	default:
		respondWithText(s, i, formatMessage("**Unknown command!**", danger))
	}
}

// createOptions creates new aternos configuration options.
func (ab *Bot) createOptions(guildId string) (*aternos.Options, error) {
	options := &aternos.Options{
		Cookies: []*http.Cookie{
			{
				Name:  "ATERNOS_LANGUAGE",
				Value: "en",
			},
		},
	}

	if ab.Database != nil {
		settings, err := ab.Database.GetServerSettings(guildId)
		if err != nil {
			return nil, err
		}

		options.Cookies = append(options.Cookies,
			&http.Cookie{
				Name:  "ATERNOS_SESSION",
				Value: settings.SessionCookie,
			},
			&http.Cookie{
				Name:  "ATERNOS_SERVER",
				Value: settings.ServerCookie,
			},
		)
	} else {
		options.Cookies = append(options.Cookies,
			&http.Cookie{
				Name:  "ATERNOS_SESSION",
				Value: ab.SessionCookie,
			},
			&http.Cookie{
				Name:  "ATERNOS_SERVER",
				Value: ab.ServerCookie,
			},
		)
	}

	return options, nil
}
