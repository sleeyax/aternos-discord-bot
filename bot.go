package aternos_discord_bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"net/http"
	"strings"
)

type AternosBot struct {
	// Command prefix.
	// Defaults to "!".
	Prefix string

	// Discord bot token.
	DiscordToken string

	// Aternos discord cookie (ATERNOS_SESSION).
	SessionCookie string

	// Aternos server cookie (ATERNOS_SERVER).
	ServerCookie string

	// Current discord session.
	discord *discordgo.Session

	// Aternos api instance.
	api *aternos.Api

	// Current server status.
	// This is basically only used as a cache for the info command.
	serverInfo *aternos.ServerInfo
}

func (ab *AternosBot) setup() {
	if ab.Prefix == "" {
		ab.Prefix = "!"
	}

	ab.api = aternos.New(&aternos.Options{
		Cookies: []*http.Cookie{
			{
				Name:  "ATERNOS_LANGUAGE",
				Value: "en",
			},
			{
				Name:  "ATERNOS_SESSION",
				Value: ab.SessionCookie,
			},
			{
				Name:  "ATERNOS_SERVER",
				Value: ab.ServerCookie,
			},
		},
	})
	ab.discord.Identify.Intents = discordgo.IntentsGuildMessages
	ab.discord.AddHandler(ab.readMessages)
}

func (ab *AternosBot) Start() error {
	var err error
	ab.discord, err = discordgo.New("Bot " + ab.DiscordToken)
	if err != nil {
		return err
	}

	ab.setup()

	return ab.discord.Open()
}

func (ab *AternosBot) Stop() error {
	return ab.discord.Close()
}

func (ab *AternosBot) GetServerInfo() (*aternos.ServerInfo, error) {
	if ab.serverInfo == nil {
		info, err := ab.api.GetServerInfo()
		if err != nil {
			return nil, err
		}
		ab.serverInfo = &info
	}

	return ab.serverInfo, nil
}

// Called whenever a message is created on any channel that the authenticated bot has access to.
func (ab *AternosBot) readMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore if message doesn't start with prefix.
	if !strings.HasPrefix(m.Content, ab.Prefix) {
		return
	}

	msg := strings.TrimLeft(m.Content, ab.Prefix)

	switch msg {
	case "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "pong":
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "status":
		fallthrough
	case "info":
		fallthrough
	case "players":
		info, err := ab.GetServerInfo()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**Unexpected error while fetching info.**")
			log.Println(err)
			break
		}

		if msg == "status" {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server %s is currently **%s**.", info.Name, info.StatusLabel))
			break
		} else if msg == "players" {
			if len(info.PlayerList) == 0 {
				s.ChannelMessageSend(m.ChannelID, "No active players found.")
				break
			}

			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Active players: %s.", strings.Join(info.PlayerList, ", ")))

			break
		}

		if info.DynIP == "" {
			info.DynIP = "unavailable"
		}

		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			Title:       "Server info",
			Description: fmt.Sprintf("Server '%s' is currently **%s**.", info.Name, info.StatusLabel),
			Color:       colorMap[info.Status],
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Made with <3 by Sleeyax.",
			},
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
		})
	case "help":
		s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			URL:   "https://github.com/sleeyax/aternos-discord-bot/",
			Title: "Available commands",
			Color: 0x00ff00,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Made with <3 by Sleeyax.",
			},
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   "start",
					Value:  "Start the server.",
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "stop",
					Value:  "Stop the server.",
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "info",
					Value:  "Show the current server status and more details.",
					Inline: false,
				},
				&discordgo.MessageEmbedField{
					Name:   "status",
					Value:  "Show the current server status.",
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   "players",
					Value:  "Show active players.",
					Inline: true,
				},
			},
		})
	case "start":
	case "stop":
	default:
		s.ChannelMessageSend(m.ChannelID, "*Command not found. Type `!help` for instructions.*")
	}
}
