package aternos_discord_bot

import (
	"github.com/bwmarrin/discordgo"
	aternos "github.com/sleeyax/aternos-api"
	"log"
	"net/http"
)

func (ab *Bot) setupHandlers() {
	// TODO: move api
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
	ab.discord.AddHandler(ab.handleCommands)
	ab.discord.AddHandler(ab.handleJoinServer)
	ab.discord.AddHandler(ab.handleLeaveServer)
}

func (ab *Bot) handleJoinServer(s *discordgo.Session, e *discordgo.GuildCreate) {
	log.Printf("Joined server %s (ID: %s)", e.Name, e.ID)
	if err := ab.registerCommands(); err != nil {
		log.Panicf("Failed to register commands: %e\n", err)
	}
}

func (ab *Bot) handleLeaveServer(s *discordgo.Session, e *discordgo.GuildDelete) {
	log.Printf("Left server %s (ID: %s)", e.Name, e.ID)
	if err := ab.removeCommands(); err != nil {
		log.Panicf("Failed to remove commands: %e\n", err)
	}
}

func (ab *Bot) Start() error {
	session, err := discordgo.New("Bot " + ab.DiscordToken)
	if err != nil {
		return err
	}

	ab.discord = session
	ab.setupHandlers()

	return ab.discord.Open()
}

func (ab *Bot) Stop() error {
	return ab.discord.Close()
}

// Called whenever a message is created on any channel that the authenticated bot has access to.
/*func (ab *Bot) readMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Ignore if message doesn't start with prefix.
	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	msg := strings.TrimLeft(m.Content, "!")

	footer := &discordgo.MessageEmbedFooter{
		Text: "Made with <3 by Sleeyax.",
	}

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
		info, err := ab.getServerInfo()
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**Unexpected error while fetching info.**")
			log.Println("failed to fetch server info:", err)
			break
		}

		if msg == "status" {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server '%s' is currently **%s**.", info.Name, info.StatusLabel))
			break
		} else if msg == "players" {
			if len(info.PlayerList) == 0 {
				s.ChannelMessageSend(m.ChannelID, "No one is playing right now :(")
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
			Footer:      footer,
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
			URL:    "https://github.com/sleeyax/aternos-discord-bot/",
			Title:  "Available commands",
			Color:  0x00ff00,
			Footer: footer,
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
		if ab.serverInfo != nil && ab.serverInfo.Status != aternos.Offline && ab.serverInfo.Status != aternos.Stopping {
			s.ChannelMessageSend(m.ChannelID, "Server already started! Type `!status` or `!info` to fetch the status.")
			break
		}

		// Connect to the websocket server.
		if _, err := ab.getWSS(); err != nil {
			s.ChannelMessageSend(m.ChannelID, "**Failed to connect to WSS.**")
			log.Println("failed to connect to WSS:", err)
			break
		}

		s.ChannelMessageSend(m.ChannelID, "Starting the server, please wait...")

		ctxHeartBeat, cancelHeartBeat := context.WithCancel(context.Background())
		ctxConfirm, cancelConfirm := context.WithCancel(context.Background())

		go func() {
			defer func() {
				cancelHeartBeat()
				cancelConfirm()
				ab.wss.Close()
				ab.wss = nil
				log.Println("Background routines stopped & connections closed.")
			}()

			for {
				msg, ok := <-ab.wss.Message

				// if the msg channel is closed it means the server unsuspectedly closed the connection, so we should try to reconnect at least once.
				if !ok {
					log.Println("Failed to read message. Tying to reconnect...")
					if _, err := ab.getWSS(); err != nil {
						s.ChannelMessageSend(m.ChannelID, "**Failed to reconnect to WSS.**")
						log.Println("failed to reconnect to WSS:", err)
						break
					}
				}

				switch msg.Type {
				case "ready":
					// Start the server.
					if err := ab.api.StartServer(); err != nil {
						s.ChannelMessageSend(m.ChannelID, "**Failed to start server.**")
						log.Println("failed to start server:", err)
						return
					}

					// Start sending keep-alive requests in the background (until the server is offline, see below).
					go ab.sendHeartBeats(ctxHeartBeat)
				case "status":
					var info aternos.ServerInfo
					json.Unmarshal(msg.MessageBytes, &info)

					switch info.Status {
					case aternos.Online:
						if info.StatusLabelClass == "online" && ab.serverInfo.StatusLabelClass != "online" {
							s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
								Title:       "Server is online",
								Description: fmt.Sprintf("Join now! Only %d seconds left.", info.Countdown),
								Color:       colorMap[aternos.Online],
								Footer:      footer,
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
							})
						}
					case aternos.Offline:
						s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
							Title:       "Server is offline",
							Description: "The server is currently offline.",
							Color:       colorMap[aternos.Offline],
							Footer:      footer,
						})
						ab.serverInfo = &info
						return
					case aternos.Preparing: // stuck in queue (only happens when traffic is high)
						if (info.StatusLabelClass == "queueing" && info.Queue.Status == "pending") && (ab.serverInfo.Queue.Status != "pending") {
							s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Waiting in queue (%d/%d, %s)...", info.Queue.Position, info.Queue.Count, info.Queue.Time))
							go ab.api.ConfirmServer(ctxConfirm, 10*time.Second)
						}
					case aternos.Loading:
						cancelConfirm()
					}

					ab.serverInfo = &info
				}
			}
		}()
	case "stop":
		if ab.serverInfo != nil && (ab.serverInfo.Status == aternos.Stopping || ab.serverInfo.Status == aternos.Offline) {
			s.ChannelMessageSend(m.ChannelID, "Server already stopped! Type `!status` or `!info` to fetch the status.")
			break
		}

		if _, err := ab.getWSS(); err != nil {
			s.ChannelMessageSend(m.ChannelID, "**Failed to connect to WSS.**")
			log.Println("failed to connect to WSS:", err)
			break
		}

		if err := ab.api.StopServer(); err != nil {
			s.ChannelMessageSend(m.ChannelID, "**Failed to stop the server.**")
			log.Println("failed to stop the server manually:", err)
			break
		}

		s.ChannelMessageSend(m.ChannelID, "Stopping the server, please wait...")
	default:
		s.ChannelMessageSend(m.ChannelID, "*Command not found. Type `!help` for instructions.*")
	}
}*/
